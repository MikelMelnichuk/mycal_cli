package cmd

import (
	"github.com/MikelMelnichuk/mycal/internal/formatter"
	"github.com/spf13/cobra"
)

// nextWeekCmd represents the "nextweek" command
var nextWeekCmd = &cobra.Command{
	Use:     "nextweek",
	Aliases: []string{"nw", "next-week"},
	Short:   "Display events for the upcoming week",
	Long: `Display a summary of all calendar events scheduled for the upcoming week.
This command provides a quick overview of your schedule for the next 7 days, 
making it easy to plan ahead. Supports outputting the results in JSON format 
for scripting and automation.`,
	Example: `  mycal nextweek
  mycal nextweek --json
  mycal nw`,
	RunE: nextWeekCommand,
}

func nextWeekCommand(cmd *cobra.Command, args []string) error {
	json, _ := cmd.Flags().GetBool("json")

	events, err := APIClient.GetWeekEvents(true, true)
	if err != nil {
		return nil
	}

	if json {
		formatter.PrintJSON(events)
	} else {
		formatter.PrettyPrintWeek(events)
	}

	return nil
}

func init() {
	nextWeekCmd.Flags().Bool("json", false, JsonDescription)

	rootCmd.AddCommand(nextWeekCmd)
}
