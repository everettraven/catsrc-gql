package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "opmgql",
	Short: "opmgql is a PoC to showcase serving FBC contents via GraphQL",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("must specify a subcommand")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
