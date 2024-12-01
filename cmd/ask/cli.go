// cmd/ask/cli.go
package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/acazau/shell-ask-go/internal/config"
	"github.com/acazau/shell-ask-go/internal/models"
	"github.com/acazau/shell-ask-go/internal/providers"
	"github.com/acazau/shell-ask-go/pkg/utils"
	"github.com/spf13/cobra"
)

type CLI struct {
	config  *config.Config
	rootCmd *cobra.Command
}

func NewCLI(cfg *config.Config) *CLI {
	cli := &CLI{
		config: cfg,
	}

	rootCmd := &cobra.Command{
		Use:   "ask [flags] [prompt]",
		Short: "CLI tool for asking questions to AI models",
		RunE:  cli.handleAsk,
	}

	// Add flags
	rootCmd.PersistentFlags().StringP("model", "m", "", "Choose the LLM to use")
	rootCmd.PersistentFlags().BoolP("command", "c", false, "Ask LLM to return a command only")
	rootCmd.PersistentFlags().StringP("type", "t", "", "Define the shape of the response")
	rootCmd.PersistentFlags().StringSliceP("url", "u", []string{}, "Fetch URL content as context")
	rootCmd.PersistentFlags().BoolP("search", "s", false, "Enable web search")
	rootCmd.PersistentFlags().Bool("no-stream", false, "Disable streaming output")
	rootCmd.PersistentFlags().BoolP("reply", "r", false, "Reply to previous conversation")

	cli.rootCmd = rootCmd
	cli.addBuiltinCommands()
	cli.addCustomCommands()

	return cli
}

func (cli *CLI) RootCmd() *cobra.Command {
	return cli.rootCmd
}

func (cli *CLI) Execute() error {
	return cli.RootCmd().Execute()
}

func (cli *CLI) handleAsk(cmd *cobra.Command, args []string) error {
	if len(args) == 0 && !utils.IsPiped() {
		return fmt.Errorf("please provide a prompt")
	}

	model, _ := cmd.Flags().GetString("model")
	if model == "" {
		model = cli.config.DefaultModel
		if model == "" {
			model = "gpt-4o-mini"
		}
	}

	prompt := strings.Join(args, " ")
	if pipeInput, err := utils.ReadPipe(); err != nil {
		return err
	} else if pipeInput != "" {
		prompt = fmt.Sprintf("%s\nInput:\n%s", prompt, pipeInput)
	}

	provider, err := cli.getProvider(model)
	if err != nil {
		return err
	}

	noStream, _ := cmd.Flags().GetBool("no-stream")
	commandOnly, _ := cmd.Flags().GetBool("command")
	if commandOnly {
		prompt += "\nReturn the command only without any other text."
	}

	return cli.processRequest(cmd.Context(), provider, prompt, !noStream)
}

func (cli *CLI) getProvider(modelID string) (providers.Provider, error) {
	model := models.SelectModel(modelID)
	switch {
	case strings.HasPrefix(model, "gpt"):
		return providers.NewOpenAIProvider(cli.config.OpenAIKey, model)
	case strings.HasPrefix(model, "claude"):
		return providers.NewAnthropicProvider(cli.config.AnthropicKey, model), nil
	default:
		return nil, fmt.Errorf("unsupported model: %s", model)
	}
}

func (cli *CLI) processRequest(ctx context.Context, provider providers.Provider, prompt string, stream bool) error {
	reader, err := provider.Complete(ctx, prompt, stream)
	if err != nil {
		return err
	}
	defer func() {
		if closer, ok := reader.(io.Closer); ok {
			closer.Close()
		}
	}()

	_, err = io.Copy(os.Stdout, reader)
	return err
}

func (cli *CLI) addBuiltinCommands() {
	cmCmd := &cobra.Command{
		Use:   "cm",
		Short: "Generate git commit message based on git diff output",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !utils.IsPiped() {
				return fmt.Errorf("this command requires git diff input")
			}
			return cli.handleAsk(cmd, []string{"Generate a git commit message based on the following diff:"})
		},
	}
	cli.rootCmd.AddCommand(cmCmd)
}

func (cli *CLI) addCustomCommands() {
	for _, cmd := range cli.config.Commands {
		c := cmd
		customCmd := &cobra.Command{
			Use:     c.Name,
			Short:   c.Description,
			Example: c.Example,
			RunE: func(cmd *cobra.Command, args []string) error {
				if c.RequireStdin && !utils.IsPiped() {
					return fmt.Errorf("this command requires piped input")
				}
				return cli.handleAsk(cmd, []string{c.Prompt})
			},
		}
		cli.rootCmd.AddCommand(customCmd)
	}
}
