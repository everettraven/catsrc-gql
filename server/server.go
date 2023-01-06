package server

import (
	"catsrc-gql/server/graph"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/operator-framework/operator-registry/alpha/declcfg"
)

const defaultPort = "8080"

type GqlServer struct {
	server *handler.Server
	port   string
}

func NewGqlServer(declarativeConfig *declcfg.DeclarativeConfig) *GqlServer {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver := graph.NewResolver(declarativeConfig)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	return &GqlServer{server: srv, port: port}
}

func (gs *GqlServer) Run() error {
	// For demonstration purposes keep the playground around. For a real
	// production grade implementation we would remove this
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", gs.server)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", gs.port)
	return http.ListenAndServe(":"+gs.port, nil)
}
