package cmd

import (
	"fmt"
	"time"

	"github.com/MikelMelnichuk/mycal/internal/formatter"
	"github.com/spf13/cobra"
)

var todayCmd = &cobra.Command{
	Use:     "today",
	Aliases: []string{"t"},
	RunE: func(cmd *cobra.Command, args []string) error {
		all, _ := cmd.Flags().GetBool("all")
		jsonOut, _ := cmd.Flags().GetBool("json")
		after, _ := cmd.Flags().GetString("after")

		// The "after" command takes priority above "all" command
		if after != "" {
			validInputFlag := isValidTime(after)
			if !validInputFlag {
				return fmt.Errorf("Could not parse the given flag after %q, time of the format HH:MM is expected", after)
			}
			// When given a starting hour, all day flag is irrelevant
			all = false
		} else {
			after = time.Now().Format(HHMM)
		}

		// Take today's date in YYYYMMDD format
		var targetDate = time.Now().Format(YYYYMMDD)
		events, err := APIClient.GetDayEvents(targetDate, all, after)
		if err != nil {
			return err
		}

		if jsonOut {
			formatter.PrintJSON(events)
		} else {
			formatter.PrettyPrintSingleDay(events, true)
		}
		return nil
	},
}

func init() {
	// Define expected flags for todayCmd
	todayCmd.Flags().Bool("all", false, AllDescription)
	todayCmd.Flags().Bool("json", false, JsonDescription)
	todayCmd.Flags().String("after", "", AfterDescription)

	// Register the command with rootCmd
	rootCmd.AddCommand(todayCmd)
}
