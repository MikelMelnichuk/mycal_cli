package cmd

import (
	"fmt"
	"time"

	"github.com/MikelMelnichuk/mycal/internal/formatter"
	"github.com/spf13/cobra"
)

const YYYYMMDD = "2006-01-02"

var todayCmd = &cobra.Command{
	Use:     "today",
	Aliases: []string{"t"},
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		jsonOut, _ := cmd.Flags().GetBool("json")
		after, _ := cmd.Flags().GetString("after")

		targetDate := time.Now().Format(YYYYMMDD)
		fmt.Printf("targetDate: %q\n", targetDate)
		events, err := APIClient.GetDayEvents(targetDate, all, after)
		if err != nil {
			return err
		}

		if jsonOut {
			formatter.PrintJSON(events)
		} else {
			formatter.PrintText(events)
		}
		return nil
	},
}

func init() {
	// Add flags to todayCmd
	todayCmd.Flags().Bool("all", false, "Show all events (including past ones)")
	todayCmd.Flags().Bool("json", false, "Output as JSON")
	todayCmd.Flags().String("after", "", "Filter events after given time (e.g., 12:00)")

	// Register the command with rootCmd
	rootCmd.AddCommand(todayCmd)
}
