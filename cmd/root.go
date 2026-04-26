package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/MikelMelnichuk/mycal/internal/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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
	Long: `mycal is a command-line tool to view and interact with your calendar events.
It communicates with a backend API to fetch events for days, weeks, or specific IDs.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help() // just show help if no subcommand is given
	},
	// Persistent pre-run hook – initializes API client before any command runs
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration (file, env, flags)
		if err := initConfig(); err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Create API client
		timeout := viper.GetDuration("api.timeout")
		if timeout == 0 {
			timeout = 10 * time.Second
		}
		APIClient = api.NewClient(apiBaseURL, timeout)

		// TODO: Optionally test connectivity?
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
	rootCmd.PersistentFlags().StringVar(&apiBaseURL, "api-url", "http://localhost:8080/api/v1", "backend API base URL")

	// Bind flags to viper so they can be overridden by env or config file
	_ = viper.BindPFlag("api.url", rootCmd.PersistentFlags().Lookup("api-url"))
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	// Set environment variable prefix: MYCAL_API_URL, MYCAL_VERBOSE, etc.
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
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mycal")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if there's no config file
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// TODO: create a default
		} else {
			return err
		}
	}

	// Override with environment variables: MYCAL_API_URL, MYCAL_VERBOSE
	if envURL := os.Getenv("MYCAL_API_URL"); envURL != "" {
		apiBaseURL = envURL
	}

	// Also read from viper config (lower priority than env)
	if viper.IsSet("api.url") && apiBaseURL == "http://localhost:8080/api/v1" {
		apiBaseURL = viper.GetString("api.url")
	}

	return nil
}
