package cmd

import (
	"fmt"
	"time"

	"github.com/MikelMelnichuk/mycal/internal/formatter"
	"github.com/spf13/cobra"
)

// mycal tomorrow
// mycal tomorrow --json
// mycal tomorrow --after 12:00
// mycal tom

var tomorrowCmd = &cobra.Command{
	Use:     "tomorrow",
	Aliases: []string{"tom"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var json, _ = cmd.Flags().GetBool("json")
		var after, _ = cmd.Flags().GetString("after")
		fmt.Printf("Tomorrow got json: %t, after: %s\n", json, after)

		all := true
		if after != "" {
			validInput := isValidTime(after)
			if !validInput {
				return fmt.Errorf("Could not parse the given flag after %q, time of the format HH:MM is expected", after)
			}
			// When given a starting hour, all day flag is irrelevant
			all = false
		}

		// Take tomorrow's date and convert to YYYYMMDD format
		var targetDate = time.Now().AddDate(0, 0, 1).Format(YYYYMMDD)
		tomorrowEvents, err := APIClient.GetDayEvents(targetDate, all, after)
		if err != nil {
			return err
		}

		if json {
			formatter.PrintJSON(tomorrowEvents)
		} else {
			formatter.PrintText(tomorrowEvents)
		}

		return nil

	},
}

func init() {
	// Define the flags that tomorrow command expects to receive
	tomorrowCmd.Flags().Bool("json", false, JsonDescription)
	tomorrowCmd.Flags().String("after", "", AfterDescription)

	rootCmd.AddCommand(tomorrowCmd)
}
