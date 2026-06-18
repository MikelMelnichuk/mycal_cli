package cmd

import (
	"github.com/spf13/cobra"
)

// testCMD represents the "test" command (aliased as "health")
var testCmd = &cobra.Command{
	Use:     "test",
	Aliases: []string{"health"},
	Short:   "Check the health and connectivity of the calendar service",
	Long: `Perform a health check to verify that the CLI is properly configured 
and can successfully connect to the calendar service. 
This is useful for debugging connection issues, verifying your credentials, 
or ensuring the API is reachable.`,
	Example: `  mycal test
  mycal health
  mycal test --verbose`,
	RunE: checkHealthCommand,
}

func checkHealthCommand(cmd *cobra.Command, args []string) error {
	var verbose, _ = cmd.Flags().GetBool("verbose")

	// If verbose check of the backend connection
	if verbose {
		var err = APIClient.HealthCheckBackend()
		if err != nil {
			return err
		}
	}

	// In any case check of connection for the backend and DB
	var err = APIClient.HealthCheckDB()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	testCmd.Flags().BoolP("verbose", "v", false, "Enable verbose mode, for more logs")
	rootCmd.AddCommand(testCmd)
}
