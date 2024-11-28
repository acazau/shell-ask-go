// cmd/ask/main.go
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/acazau/shell-ask-go/internal/cli"
	"github.com/acazau/shell-ask-go/internal/config"
	"github.com/acazau/shell-ask-go/pkg/version"
)

var rootCmd = &cobra.Command{
	Use:   "ask",
	Short: "CLI tool for asking questions to AI models",
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(cmd.OutOrStdout(), version.GetVersionInfo())
		},
	})
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	cli := cli.NewCLI(cfg)
	rootCmd.AddCommand(cli.RootCmd())
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %v\n", err)
		os.Exit(1)
	}
}

