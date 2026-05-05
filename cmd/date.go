package cmd

import (
	"fmt"
	"strings"

	"github.com/MikelMelnichuk/mycal/internal/formatter"
	"github.com/spf13/cobra"
)

var dateCmd = &cobra.Command{
	Use:     "date",
	Aliases: []string{"d"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var json, _ = cmd.Flags().GetBool("json")
		var after, _ = cmd.Flags().GetString("after")

		fmt.Printf("Date got json: %t, after: %s\n", json, after)

		all := true
		if after != "" {
			validInput := isValidTime(after)
			if !validInput {
				return fmt.Errorf("Could not parse the given flag after %q, time of the format HH:MM is expected", after)
			}
			// When given a starting hour, all day flag is irrelevant
			all = false
		}

		if len(args) == 0 {
			return fmt.Errorf("No date description supplied.\nPlease provide natural language date description (e.g., 'next monday') or ISO (YYYY-MM-DD).")
		}

		// Combine the words into a sentence to be sent
		targetDate := strings.Join(args, " ")
		dateEvents, err := APIClient.GetDayEvents(targetDate, all, after)
		if err != nil {
			return err
		}

		if json {
			formatter.PrintJSON(dateEvents)
		} else {
			formatter.PrintText(dateEvents)
		}

		return nil
	},
}

func init() {
	// Define the flags that given command expects to receive
	dateCmd.Flags().Bool("json", false, "Output as JSON")
	dateCmd.Flags().String("after", "", "Filter events after given time (e.g., 12:00)")

	rootCmd.AddCommand(dateCmd)
}
