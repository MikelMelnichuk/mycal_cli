package cmd

import (
	"github.com/spf13/cobra"
)

// mycal health
// mycal test
// mycal test --verbose

var testCMD = &cobra.Command{
	Use:     "test",
	Aliases: []string{"health"},
	RunE: func(cmd *cobra.Command, args []string) error {
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
	},
}

func init() {
	testCMD.Flags().BoolP("verbose", "v", false, "Enable verbose mode, for more logs")
	rootCmd.AddCommand(testCMD)
}
