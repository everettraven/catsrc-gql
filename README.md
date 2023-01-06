## catsrc-gql
`catsrc-gql` is a prototype to showcase what it would be like to serve [File-Based Catalogs (FBC)](https://olm.operatorframework.io/docs/reference/file-based-catalogs/)
via GraphQL as opposed to serving it via gRPC like the current [`opm`](https://github.com/operator-framework/operator-registry) tool does when serving a catalog.

The advantages of using GraphQL:
- Puts the responsibility on the clients to request exactly what they need
- Clients can create highly customized requests
- Response time is faster for smaller requests
    - Example:
        - Using the quay.io/operatorhubio/catalog:latest (image of community-operators built via: [k8s-operatorhub/community-operators](https://github.com/k8s-operatorhub/community-operators))
        - When using GraphQL to request the entire contents of a catalog the response time is 1.63s and the response size is 76.3MB
        - When using GraphQL to request only catalog content that relates to the `cc-operator` package the response time is 5.07ms and response size is 288.2kb (note that this is a pretty small operator. The larger the operator the longer the response time and the larger the response size)

## Tutorial
This tutorial will guide you through serving a FBC via GraphQL with the `catsrc-gql` binary.

### Prerequisites
In order to properly
- Go 1.19
- Docker
- GNU Make (optional but recommended)
- Git
- Any way to send/receive a HTTP request/response (this tutorial will use `curl` to keep it as barebones as possible)

### Initial Setup
Since this is a prototype, there is no binary readily available for install and will need to be built manually.

To do so:
1. Clone the repository with:
```
git clone https://github.com/everettraven/catsrc-gql.git
```
2. Change directory into the newly cloned repo:
```
cd catsrc-gql
```
3. Build the binary:
(if you have GNU Make):
```
make build
```
(if you DO NOT have GNU Make):
```
go build -o catsrc-gql main.go
```

Now that the binary is built we are good to go!

### Serve the FBC
For simplicity we will serve from an image, however the `catsrc-gql` binary does have the ability to serve the FBC from a directory.

For the FBC image, lets use the quay.io/operatorhubio/catalog:latest (image of community-operators built via: [k8s-operatorhub/community-operators](https://github.com/k8s-operatorhub/community-operators)
```
./catsrc-gql serve --from-image quay.io/operatorhubio/catalog:latest
```

After a few seconds you should see an output similar to:
```
2023/01/06 10:49:59 connect to http://localhost:8080/ for GraphQL playground
```

### Querying the Server
Navigating in a web browser to http://localhost:8080/ will load a GraphQL playground that will allow you to write GraphQL queries and run them to see the response. This is a great way to play around and see the results of custom queries.

If you are using a tool other than the GraphQL playground all requests should be sent to http://localhost:8080/query via the `POST` HTTP Method.

To make this tutorial as accessible as possible we will use `curl` for sending requests, but will also include a separate block to show the GraphQL query used so that it can easily be copied and used in other tooling.

**Fetch all Catalog Contents (entire catalog)**
First we will demonstrate fetching the entire catalog:

GraphQL Query:
```graphql
query {
  packages {
    schema
    name
    description
    defaultChannel
    icon {
      mediatype
      base64data
    }
  }
  bundles {
    schema
    package
    name
    image
    properties {
      type
      value
    }
    relatedImages {
      image
      name
    }
  }
  channels {
    schema
    package
    name
    entries {
      name
      replaces
      skips
      skipRange
    }
  }
  metas {
    schema
    package
    blob
  }
}
```

Querying with `curl`:
```
curl -X POST \
 -H "Content-Type: application/json" \
 -d '{"query": "query { packages { schema name description defaultChannel icon { mediatype base64data } } bundles { schema package name image properties { type value } relatedImages { image name } } channels { schema package name entries { name replaces skips skipRange } } metas { schema package blob } }"}' \
 localhost:8080/query
```
Running this query results in a HUGE output, the catalog is ~76.3MB and we just asked for the whole thing!

While this is a query that is possible, clients may want to slim down the information they request. Lets try another query that produces a smaller response.

**Fetching all Catalog Contents for a specific package**
Imagine a scenario where I know what package I want, but I need all the catalog information for that particular package. This query will show how to do that, and for this example we will get all the catalog contents for the `cc-operator`:

GraphQL Query:
```graphql
query {
  packages(packageName: "cc-operator") {
    schema
    name
    description
    defaultChannel
    icon {
      mediatype
      base64data
    }
  }
  bundles(packageName: "cc-operator") {
    schema
    package
    name
    image
    properties {
      type
      value
    }
    relatedImages {
      image
      name
    }
  }
  channels(packageName: "cc-operator") {
    schema
    package
    name
    entries {
      name
      replaces
      skips
      skipRange
    }
  }
  metas(packageName: "cc-operator") {
    schema
    package
    blob
  }
}
```

Querying with `curl`:
```
curl -X POST \
 -H "Content-Type: application/json" \
 -d '{"query": "query { packages(packageName: \"cc-operator\") { schema name description defaultChannel icon { mediatype base64data } } bundles(packageName: \"cc-operator\") { schema package name image properties { type value } relatedImages { image name } } channels(packageName: \"cc-operator\") { schema package name entries { name replaces skips skipRange } } metas(packageName: \"cc-operator\") { schema package blob } }"}' \
 localhost:8080/query
```

The output is still a lot of information, but we can see that the response is significantly faster because we are *only* fetching the catalog resources that are related to the `cc-operator` package.

**Fetching the names of all packages in the catalog**
Now consider a scenario where we don't know what packages are available in a catalog. A client wouldn't want to have to parse the entire catalog when they only want the names of available packages. Luckily with GraphQL we can write a query to do just that:

GraphQL Query:
```graphql
query {
  packages {
    name
  }
}
```

Querying with `curl`:
```
curl -X POST \
 -H "Content-Type: application/json" \
 -d '{"query": "query { packages { name } }"}' \
 localhost:8080/query
```

<details>
    <summary>Output</summary>

```json
{
  "data": {
    "packages": [
      { "name": "ack-apigatewayv2-controller" },
      { "name": "ack-applicationautoscaling-controller" },
      { "name": "ack-cloudtrail-controller" },
      { "name": "ack-dynamodb-controller" },
      { "name": "ack-ec2-controller" },
      { "name": "ack-ecr-controller" },
      { "name": "ack-eks-controller" },
      { "name": "ack-elasticache-controller" },
      { "name": "ack-emrcontainers-controller" },
      { "name": "ack-iam-controller" },
      { "name": "ack-kinesis-controller" },
      { "name": "ack-kms-controller" },
      { "name": "ack-lambda-controller" },
      { "name": "ack-memorydb-controller" },
      { "name": "ack-mq-controller" },
      { "name": "ack-opensearchservice-controller" },
      { "name": "ack-prometheusservice-controller" },
      { "name": "ack-rds-controller" },
      { "name": "ack-s3-controller" },
      { "name": "ack-sagemaker-controller" },
      { "name": "ack-sfn-controller" },
      { "name": "ack-sns-controller" },
      { "name": "ack-sqs-controller" },
      { "name": "aerospike-kubernetes-operator" },
      { "name": "aiven-operator" },
      { "name": "akka-cluster-operator" },
      { "name": "alvearie-imaging-ingestion" },
      { "name": "anchore-engine" },
      { "name": "api-operator" },
      { "name": "apicast-community-operator" },
      { "name": "apicurio-registry" },
      { "name": "apimatic-kubernetes-operator" },
      { "name": "app-director-operator" },
      { "name": "appdynamics-operator" },
      { "name": "application-services-metering-operator" },
      { "name": "appranix" },
      { "name": "appsody-operator" },
      { "name": "aqua" },
      { "name": "argocd-operator" },
      { "name": "argocd-operator-helm" },
      { "name": "atlasmap-operator" },
      { "name": "authorino-operator" },
      { "name": "aws-auth-operator" },
      { "name": "awss3-operator-registry" },
      { "name": "awx-operator" },
      { "name": "azure-service-operator" },
      { "name": "banzaicloud-kafka-operator" },
      { "name": "beegfs-csi-driver-operator" },
      { "name": "bookkeeper-operator" },
      { "name": "camel-k" },
      { "name": "camel-karavan-operator" },
      { "name": "carbonetes-operator" },
      { "name": "cass-operator-community" },
      { "name": "cassandra-operator" },
      { "name": "cc-operator" },
      { "name": "cert-manager" },
      { "name": "chaosblade-operator" },
      { "name": "clever-operator" },
      { "name": "clickhouse" },
      { "name": "cloud-native-postgresql" },
      { "name": "cluster-aas-operator" },
      { "name": "cluster-impairment-operator" },
      { "name": "cluster-manager" },
      { "name": "cockroachdb" },
      { "name": "community-kubevirt-hyperconverged" },
      { "name": "community-trivy-operator" },
      { "name": "composable-operator" },
      { "name": "cos-bucket-operator" },
      { "name": "couchbase-enterprise" },
      { "name": "credstash-operator" },
      { "name": "cryostat-operator" },
      { "name": "datadog-operator" },
      { "name": "datatrucker-operator" },
      { "name": "dell-csi-operator" },
      { "name": "dell-csm-operator" },
      { "name": "deployment-validation-operator" },
      { "name": "ditto-operator" },
      { "name": "dnext-operator" },
      { "name": "druid-operator" },
      { "name": "dynatrace-operator" },
      { "name": "ecr-secret-operator" },
      { "name": "edp-keycloak-operator" },
      { "name": "eginnovations-operator" },
      { "name": "elastic-cloud-eck" },
      { "name": "elastic-phenix-operator" },
      { "name": "ember-csi-operator" },
      { "name": "enc-key-sync" },
      { "name": "enmasse" },
      { "name": "esindex-operator" },
      { "name": "etcd" },
      { "name": "event-streams-topic" },
      { "name": "eventing-kogito" },
      { "name": "ext-postgres-operator" },
      { "name": "external-secrets-operator" },
      { "name": "falco" },
      { "name": "falcon-operator" },
      { "name": "federatorai" },
      { "name": "flagsmith" },
      { "name": "flink-kubernetes-operator" },
      { "name": "flux" },
      { "name": "fossul-operator" },
      { "name": "function-mesh" },
      { "name": "galaxy-operator" },
      { "name": "gitlab-operator-kubernetes" },
      { "name": "gitlab-runner-operator" },
      { "name": "grafana-operator" },
      { "name": "halkyon" },
      { "name": "ham-deploy" },
      { "name": "hawkbit-operator" },
      { "name": "hazelcast-platform-operator" },
      { "name": "hedvig-operator" },
      { "name": "hive-operator" },
      { "name": "horreum-operator" },
      { "name": "hpa-operator" },
      { "name": "hpe-csi-operator" },
      { "name": "hpe-ezmeral-csi-operator" },
      { "name": "hyperfoil-bundle" },
      { "name": "ibm-application-gateway-operator" },
      { "name": "ibm-block-csi-operator-community" },
      { "name": "ibm-quantum-operator" },
      { "name": "ibm-security-verify-access-operator" },
      { "name": "ibm-spectrum-scale-csi-operator" },
      { "name": "ibmcloud-iam-operator" },
      { "name": "ibmcloud-operator" },
      { "name": "infinispan" },
      { "name": "instana-agent-operator" },
      { "name": "integrity-shield-operator" },
      { "name": "intel-device-plugins-operator" },
      { "name": "iot-simulator" },
      { "name": "ipfs-operator" },
      { "name": "istio" },
      { "name": "istio-workspace-operator" },
      { "name": "jaeger" },
      { "name": "jenkins-operator" },
      { "name": "joget-tomcat-operator" },
      { "name": "k8gb" },
      { "name": "keda" },
      { "name": "keycloak-operator" },
      { "name": "kiali" },
      { "name": "klusterlet" },
      { "name": "knative-operator" },
      { "name": "kogito-operator" },
      { "name": "kom-operator" },
      { "name": "kong" },
      { "name": "kong-gateway-operator" },
      { "name": "ks-releaser-operator" },
      { "name": "kuadrant-operator" },
      { "name": "kube-arangodb" },
      { "name": "kube-green" },
      { "name": "kubefed-operator" },
      { "name": "kubeflow" },
      { "name": "kubemod" },
      { "name": "kubemq-operator" },
      { "name": "kubernetes-imagepuller-operator" },
      { "name": "kubernetes-nmstate-operator" },
      { "name": "kubero-operator" },
      { "name": "kubestone" },
      { "name": "kubeturbo" },
      { "name": "lbconfig-operator" },
      { "name": "lib-bucket-provisioner" },
      { "name": "lightbend-console-operator" },
      { "name": "limitador-operator" },
      { "name": "litmuschaos" },
      { "name": "log2rbac" },
      { "name": "logging-operator" },
      { "name": "machine-deletion-operator" },
      { "name": "mariadb-operator-app" },
      { "name": "marin3r" },
      { "name": "mattermost-operator" },
      { "name": "mcad-operator" },
      { "name": "mercury-operator" },
      { "name": "meshery-operator" },
      { "name": "metallb-operator" },
      { "name": "metering-upstream" },
      { "name": "microcks" },
      { "name": "minio-operator" },
      { "name": "mondoo-operator" },
      { "name": "mongodb-atlas-kubernetes" },
      { "name": "mongodb-enterprise" },
      { "name": "mongodb-operator" },
      { "name": "multi-nic-cni-operator" },
      { "name": "multicluster-operators-subscription" },
      { "name": "mysql" },
      { "name": "myvirtualdirectory" },
      { "name": "ndmspc-operator" },
      { "name": "netobserv-operator" },
      { "name": "neuvector-operator" },
      { "name": "nexus-operator-m88i" },
      { "name": "nfd-operator" },
      { "name": "nfs-provisioner-operator" },
      { "name": "node-healthcheck-operator" },
      { "name": "node-maintenance-operator" },
      { "name": "noobaa-operator" },
      { "name": "nsm-operator-registry" },
      { "name": "nuxeo-operator" },
      { "name": "oneagent" },
      { "name": "open-liberty" },
      { "name": "openebs" },
      { "name": "openshift-qiskit-operator" },
      { "name": "opentelemetry-operator" },
      { "name": "opsmx-spinnaker-operator" },
      { "name": "otc-rds-operator" },
      { "name": "ovms-operator" },
      { "name": "patterns-operator" },
      { "name": "pcc-operator" },
      { "name": "percona-postgresql-operator" },
      { "name": "percona-server-mongodb-operator" },
      { "name": "percona-xtradb-cluster-operator" },
      { "name": "pixie-operator" },
      { "name": "planetscale" },
      { "name": "pmem-csi-operator" },
      { "name": "portworx" },
      { "name": "portworx-essentials" },
      { "name": "postgres-operator" },
      { "name": "postgresql" },
      { "name": "postgresql-operator" },
      { "name": "postgresql-operator-dev4devs-com" },
      { "name": "project-quay" },
      { "name": "project-quay-container-security-operator" },
      { "name": "prometheus" },
      { "name": "prometheus-exporter-operator" },
      { "name": "prometurbo" },
      { "name": "pulp-operator" },
      { "name": "pulsar-operator" },
      { "name": "pystol" },
      { "name": "qserv-operator" },
      { "name": "rabbitmq-cluster-operator" },
      { "name": "rabbitmq-messaging-topology-operator" },
      { "name": "rabbitmq-single-active-consumer-operator" },
      { "name": "radanalytics-spark" },
      { "name": "redis-enterprise" },
      { "name": "redis-operator" },
      { "name": "ripsaw" },
      { "name": "robin-operator" },
      { "name": "rocketmq-operator" },
      { "name": "rook-ceph" },
      { "name": "rook-edgefs" },
      { "name": "routernetes-operator" },
      { "name": "runtime-component-operator" },
      { "name": "sap-btp-operator" },
      { "name": "searchpe-operator" },
      { "name": "security-profiles-operator" },
      { "name": "seldon-operator" },
      { "name": "self-node-remediation" },
      { "name": "sematext" },
      { "name": "service-binding-operator" },
      { "name": "shipwright-operator" },
      { "name": "siddhi-operator" },
      { "name": "sigstore-helm-operator" },
      { "name": "silicom-sts-operator" },
      { "name": "skydive-operator" },
      { "name": "snapscheduler" },
      { "name": "snyk-operator" },
      { "name": "sosivio" },
      { "name": "spark-gcp" },
      { "name": "spinnaker-operator" },
      { "name": "splunk" },
      { "name": "starboard-operator" },
      { "name": "steerd-presto-operator" },
      { "name": "storageos" },
      { "name": "strimzi-kafka-operator" },
      { "name": "submariner" },
      { "name": "synapse-helm" },
      { "name": "synapse-operator" },
      { "name": "t8c" },
      { "name": "tackle-operator" },
      { "name": "tagger" },
      { "name": "tektoncd-operator" },
      { "name": "telegraf-operator" },
      { "name": "temporal-operator" },
      { "name": "tf-controller" },
      { "name": "tf-operator" },
      { "name": "tidb-operator" },
      { "name": "topolvm-operator" },
      { "name": "traefikee-operator" },
      { "name": "trident-operator" },
      { "name": "trivy-operator" },
      { "name": "ublhub-operator" },
      { "name": "universal-crossplane" },
      { "name": "varnish-operator" },
      { "name": "vault" },
      { "name": "vault-helm" },
      { "name": "verticadb-operator" },
      { "name": "victoriametrics-operator" },
      { "name": "virt-gateway-operator" },
      { "name": "wavefront" },
      { "name": "wildfly" },
      { "name": "windup-operator" },
      { "name": "wso2am-operator" },
      { "name": "xrootd-operator" },
      { "name": "yaks" },
      { "name": "yugabyte-operator" },
      { "name": "zookeeper-operator" },
      { "name": "zoperator" }
    ]
  }
}
```

> **Note:** 
> This output has been prettified to make it easier to read 
</details>

**Fetching channels for a specific package**
Now imagine a scenario where we know what package we want to install from the catalog, but want to see what channels we can use when installing the catalog. Let's run that GraphQL query! For this example we will use the `api-operator` since it has a couple different channels:

GraphQL query:
```graphql
query {
  channels(packageName: "api-operator") {
    name
    entries {
      name
      replaces
      skips
      skipRange
    }
  }
}
```

Querying with `curl`:
```
curl -X POST \
 -H "Content-Type: application/json" \
 -d '{"query": "query { channels(packageName: \"api-operator\") { name entries { name replaces skips skipRange } } }"}' \
 localhost:8080/query
```

<details>
    <summary>Output</summary>

```json
{
  "data": {
    "channels": [
      {
        "name": "2.x-stable",
        "entries": [
          {
            "name": "api-operator.v2.0.0",
            "replaces": "",
            "skips": null,
            "skipRange": ""
          }
        ]
      },
      {
        "name": "stable",
        "entries": [
          {
            "name": "api-operator.v1.0.1",
            "replaces": "api-operator.v1.2.2",
            "skips": null,
            "skipRange": ""
          },
          {
            "name": "api-operator.v1.1.0",
            "replaces": "api-operator.v1.2.2",
            "skips": null,
            "skipRange": ""
          },
          {
            "name": "api-operator.v1.2.0",
            "replaces": "api-operator.v1.2.2",
            "skips": null,
            "skipRange": ""
          },
          {
            "name": "api-operator.v1.2.1",
            "replaces": "api-operator.v1.2.2",
            "skips": null,
            "skipRange": ""
          },
          {
            "name": "api-operator.v1.2.2",
            "replaces": "api-operator.v1.2.2",
            "skips": null,
            "skipRange": ""
          },
          {
            "name": "api-operator.v1.2.3",
            "replaces": "api-operator.v1.2.2",
            "skips": null,
            "skipRange": ""
          }
        ]
      }
    ]
  }
}
```
> **Note:** 
> This output has been prettified to make it easier to read 
</details>

**Fetching a specific bundle**
Now that I know all the bundles in each channel for the `api-operator`, imagine a scenario where I decided I'd like to install the bundle `api-operator.v2.0.0` from the `2.x-stable` channel and now I need all the information for that specific bundle. Lets formulate that GraphQL query:

GraphQL query:
```graphql
query {
  bundles(bundleName: "api-operator.v2.0.0") {
    package
    name
    image
  }
}
```
in this query we are omitting the `properties` field to keep the response as small as possible, but it is pretty easy to include them if desired. We could also be extra explicit and specify the package name in the query parameters like `bundles(bundleName: "api-operator.v2.0.0", packageName: "api-operator")` in the event that some other bundle existed with the same name for a different package. In this case we know that is not the case so we omit the `packageName` parameter.

Querying with `curl`:
```
curl -X POST \
 -H "Content-Type: application/json" \
 -d '{"query": "query { bundles(bundleName: \"api-operator.v2.0.0\") { package name image } }"}' \
 localhost:8080/query
```

We should see an output similar to:
```json
{
  "data": {
    "bundles": [
      {
        "package": "api-operator",
        "name": "api-operator.v2.0.0",
        "image": "quay.io/operatorhubio/api-operator@sha256:30af98b932f87c58ce659b19651901189059014f23e50eaceb80bf81fd242a90"
      }
    ]
  }
}
```
> **Note:** 
> This output has been prettified to make it easier to read

This is the end of the tutorial as the goal was to just give a very basic feel of what could be done with GraphQL.

Feel free to play around some more and craft some queries of your own! Here are a few resources that may be helpful:
- [sample-queries.graphql](./sample-queries.graphql) - file in this repo with some sample queries used during the process of building and testing this prototype
- [GraphQL Queries Documentation](https://graphql.org/learn/queries/) - great documentation on writing queries from the GraphQL Foundation's website