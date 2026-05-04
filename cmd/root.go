package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/MikelMelnichuk/mycal/internal/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const YYYYMMDD = "2006-01-02"

var (
	cfgFile    string
	apiBaseURL string

	// Global API client – can be used by all subcommands
	APIClient *api.Client
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mycal",
	Short: "CLI for managing calendar events",
	Long: `mycal is a CLI tool to view and interact with your calendar events.
It communicates with a backend API to fetch events for days, weeks, or specific IDs.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() // just show help if no subcommand is given
	},
	// Persistent pre-run hook initializes API client before any command runs
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration (file, env, flags)
		if err := initConfig(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Create API client
		timeout := viper.GetDuration("api.timeout")
		if timeout == 0 {
			timeout = api.DEFAULT_TIMEOUT
		}
		fmt.Printf("apiBaseURL: %q, timeout: %s\n", apiBaseURL, timeout)
		APIClient = api.NewClient(apiBaseURL, timeout)

		// TODO: test connectivity? we have / and /health
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main().
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Cobra already prints the error, but we can add extra handling
		os.Exit(1)
	}
}

func init() {
	// Global persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mycal.yaml)")
	rootCmd.PersistentFlags().StringVar(&apiBaseURL, "api-url", api.DEFAULT_BASE_IP, "backend API base URL")

	// Bind flags to viper so they can be overridden by env or config file
	_ = viper.BindPFlag("api.url", rootCmd.PersistentFlags().Lookup("api-url"))
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Set environment variable prefix: MYCAL_API_URL, MYCAL_TIMEOUT, etc.
	viper.SetEnvPrefix("MYCAL")
	viper.AutomaticEnv()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() error {
	if cfgFile != "" {
		// Use config file given through the CLI flag
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}
		// TODO: move to ~/.config/
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mycal")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there's no config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := createDefaultConfig(); err != nil {
				return err
			}
			// Attempt to read the config after its creation
			if err := viper.ReadInConfig(); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	// Override with environment variables: MYCAL_API_URL and MYCAL_TIMEOUT
	if envURL := os.Getenv("MYCAL_API_URL"); envURL != "" {
		apiBaseURL = envURL
	}

	// Also read from viper config (lower priority than env)
	if viper.IsSet("api.url") && apiBaseURL == api.DEFAULT_BASE_IP {
		apiBaseURL = viper.GetString("api.url")
	}

	return nil
}

func createDefaultConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".mycal.yaml")

	// Check if config already exists; if yes, do nothing.
	if _, err := os.Stat(configPath); err == nil {
		// File exists – nothing to do.
		return nil
	} else if !os.IsNotExist(err) {
		// Some other error (e.g., permission denied).
		return err
	}

	// Default configuration content.
	defaultConfig := fmt.Sprintf(`api:
  url: %s
  timeout: %s
`, api.DEFAULT_BASE_IP, api.DEFAULT_TIMEOUT)

	// Write the file with read/write permissions for the owner.
	err = os.WriteFile(configPath, []byte(defaultConfig), 0644)
	if err != nil {
		return err
	}

	return nil
}
