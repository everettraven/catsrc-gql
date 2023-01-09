package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"catsrc-gql/server"

	"github.com/operator-framework/operator-registry/alpha/action"
	"github.com/operator-framework/operator-registry/pkg/image/containerdregistry"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts a GraphQL API server to serve FBC content",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		//Do something
		return serve(args[0])
	},
}

func serve(ref string) error {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	refs := []string{ref}
	reg, err := containerdregistry.NewRegistry(containerdregistry.WithLog(logrus.NewEntry(logger)))
	if err != nil {
		return fmt.Errorf("error creating new image registry: %v", err)
	}

	defer func() {
		err = reg.Destroy()
		if err != nil {
			fmt.Println("Unable to cleanup registry")
		}
	}()

	render := action.Render{
		Refs:     refs,
		Registry: reg,
	}

	log.SetOutput(io.Discard)
	dcfg, err := render.Run(context.Background())
	log.SetOutput(os.Stdout)
	if err != nil {
		return fmt.Errorf("failed to render FBC: %s", err)
	}

	gqlServer := server.NewGqlServer(dcfg)
	return gqlServer.Run()
}
