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
	Short:   "Display events for a specific date",
	Long: `Display calendar events for a specific date.
Supports natural language dates (e.g., "next monday", "today") or ISO format (YYYY-MM-DD).
Allows filtering events by time or outputting the results in JSON format.`,
	Example: `  mycal date 2026-06-18
  mycal date 2026-03-25  # Any date
  mycal date 2026-06-18 --json
  mycal date 2026-06-18 --after 12:00
  mycal date next monday # We'll do our best to understand your text
  mycal d`,
	RunE: runDateCommand,
}

func runDateCommand(cmd *cobra.Command, args []string) error {
	var json, _ = cmd.Flags().GetBool("json")
	var after, _ = cmd.Flags().GetString("after")

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
		formatter.PrettyPrintSingleDay(dateEvents, true)
	}

	return nil
}

func init() {
	// Define the flags that given command expects to receive
	dateCmd.Flags().Bool("json", false, JsonDescription)
	dateCmd.Flags().String("after", "", AfterDescription)

	rootCmd.AddCommand(dateCmd)
}
