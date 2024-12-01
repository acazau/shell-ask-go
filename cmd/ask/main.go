// cmd/ask/main.go
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/acazau/shell-ask-go/internal/config"
	"github.com/acazau/shell-ask-go/internal/copilot"
	"github.com/acazau/shell-ask-go/internal/models"
	"github.com/acazau/shell-ask-go/internal/providers"
	"github.com/acazau/shell-ask-go/pkg/utils"
	"github.com/acazau/shell-ask-go/pkg/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ask [prompt]",
	Short: "CLI tool for asking questions to AI models",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 && !utils.IsPiped() {
			return fmt.Errorf("please provide a prompt")
		}

		// Get model flag and handle default case
		modelFlag, _ := cmd.Flags().GetString("model")
		if modelFlag == "" {
			modelFlag = "gpt-4" // Default model from models.SelectModel
		}

		// Handle the prompt input
		prompt := strings.Join(args, " ")
		if pipeInput, err := utils.ReadPipe(); err != nil {
			return err
		} else if pipeInput != "" {
			prompt = fmt.Sprintf("%s\nInput:\n%s", prompt, pipeInput)
		}

		// Get command-only flag and append instruction if needed
		commandOnly, _ := cmd.Flags().GetBool("command")
		if commandOnly {
			prompt += "\nReturn the command only without any other text."
		}

		// Get breakdown flag and append instruction if needed
		breakdown, _ := cmd.Flags().GetBool("breakdown")
		if breakdown {
			prompt += "\nProvide a detailed breakdown of what the command does."
		}

		// Handle files context
		files, _ := cmd.Flags().GetString("files")
		if files != "" {
			fileContent, err := utils.ReadFiles(strings.Split(files, ","))
			if err != nil {
				return fmt.Errorf("failed to read files: %w", err)
			}
			prompt = fmt.Sprintf("Files content:\n%s\n\nPrompt: %s", fileContent, prompt)
		}

		// Handle URL context
		urls, _ := cmd.Flags().GetStringSlice("url")
		if len(urls) > 0 {
			urlContent, err := utils.FetchURLs(urls)
			if err != nil {
				return fmt.Errorf("failed to fetch URLs: %w", err)
			}
			prompt = fmt.Sprintf("URL content:\n%s\n\nPrompt: %s", urlContent, prompt)
		}

		// Get streaming flag
		noStream, _ := cmd.Flags().GetBool("no-stream")

		// Initialize provider
		provider, err := providers.InitializeProvider(modelFlag)
		if err != nil {
			return fmt.Errorf("failed to initialize provider: %w", err)
		}

		// Process the request
		return providers.ProcessRequest(cmd.Context(), provider, prompt, !noStream)
	},
}

func init() {
	// Add all the command flags
	rootCmd.PersistentFlags().StringP("model", "m", "", "Choose the LLM to use")
	rootCmd.PersistentFlags().BoolP("command", "c", false, "Ask LLM to return a command only")
	rootCmd.PersistentFlags().BoolP("breakdown", "b", false, "Ask LLM to return a command and the breakdown")
	rootCmd.PersistentFlags().String("files", "", "Adding files to model context")
	rootCmd.PersistentFlags().StringP("type", "t", "", "Define the shape of the response")
	rootCmd.PersistentFlags().StringSliceP("url", "u", []string{}, "Fetch URL content as context")
	rootCmd.PersistentFlags().BoolP("search", "s", false, "Enable web search")
	rootCmd.PersistentFlags().Bool("no-stream", true, "Disable streaming output")
	rootCmd.PersistentFlags().BoolP("reply", "r", false, "Reply to previous conversation")

	// Add built-in commands
	addBuiltinCommands()
}

func addBuiltinCommands() {
	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available models",
		RunE: func(cmd *cobra.Command, args []string) error {
			includeOllama, _ := cmd.Flags().GetBool("include-ollama")
			models := models.GetAllModels(includeOllama)

			if len(models) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No models available.")
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Available models:")
			for _, model := range models {
				desc := model.Description
				if desc == "" {
					desc = model.Family
				}
				fmt.Fprintf(cmd.OutOrStdout(), "- %s (%s)\n", model.ID, desc)
			}

			return nil
		},
	}
	listCmd.Flags().Bool("include-ollama", false, "Include Ollama models in the list")
	rootCmd.AddCommand(listCmd)

	// Version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprint(cmd.OutOrStdout(), version.GetVersionInfo())
		},
	})

	// Git commit message command
	cmCmd := &cobra.Command{
		Use:   "cm",
		Short: "Generate git commit message based on git diff output",
		RunE: func(cmd *cobra.Command, args []string) error {
			if !utils.IsPiped() {
				return fmt.Errorf("this command requires git diff input")
			}

			modelFlag, _ := cmd.Flags().GetString("model")
			if modelFlag == "" {
				modelFlag = models.GetCheapModel("gpt-4") // Use cheaper model for commit messages
			}

			provider, err := providers.InitializeProvider(modelFlag)
			if err != nil {
				return fmt.Errorf("failed to initialize provider: %w", err)
			}

			prompt := "Generate a git commit message based on the following diff:"
			return providers.ProcessRequest(cmd.Context(), provider, prompt, true)
		},
	}
	rootCmd.AddCommand(cmCmd)

	// Copilot login command
	copilotLoginCmd := &cobra.Command{
		Use:   "copilot-login",
		Short: "Login to GitHub Copilot",
		RunE: func(cmd *cobra.Command, args []string) error {
			copilotClient := copilot.New(config.GetConfigDir())

			deviceCode, err := copilotClient.RequestDeviceCode()
			if err != nil {
				return fmt.Errorf("failed to request device code: %w", err)
			}

			fmt.Println("First copy your one-time code:")
			fmt.Printf("\033[1m\033[32m%s\033[0m\n\n", deviceCode.UserCode)
			fmt.Printf("Then visit this GitHub URL to authorize: %s\n\n", deviceCode.VerificationURI)
			fmt.Println("Waiting for authentication...")
			fmt.Println("Press Enter to check the authentication status...")

			// Wait for Enter key
			var input string
			fmt.Scanln(&input)

			// Start polling for auth after Enter is pressed
			ticker := time.NewTicker(time.Duration(deviceCode.Interval) * time.Second)
			defer ticker.Stop()

			fmt.Println("Checking authentication status...")
			for {
				auth, err := copilotClient.VerifyAuth(deviceCode.DeviceCode)
				if err != nil {
					return fmt.Errorf("authentication failed: %w", err)
				}

				if auth != nil {
					fmt.Printf("Received token of length: %d\n", len(auth.AccessToken))
					if err := copilotClient.SaveAuthToken(auth.AccessToken); err != nil {
						return fmt.Errorf("failed to save auth token: %w", err)
					}
					fmt.Println("Authentication successful!")
					return nil
				}

				<-ticker.C // Wait for next tick before retrying
			}
		},
	}
	rootCmd.AddCommand(copilotLoginCmd)
	rootCmd.AddCommand(copilotLoginCmd)
	rootCmd.AddCommand(copilotLoginCmd)

	// Copilot logout command
	copilotLogoutCmd := &cobra.Command{
		Use:   "copilot-logout",
		Short: "Logout from GitHub Copilot",
		RunE: func(cmd *cobra.Command, args []string) error {
			copilotClient := copilot.New(config.GetConfigDir())
			if err := copilotClient.RemoveAuthToken(); err != nil {
				return fmt.Errorf("failed to remove auth token: %w", err)
			}
			fmt.Println("Copilot auth token removed")
			return nil
		},
	}
	rootCmd.AddCommand(copilotLogoutCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
