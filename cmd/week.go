package cmd

import (
	"github.com/MikelMelnichuk/mycal/internal/formatter"
	"github.com/spf13/cobra"
)

// weekCMD represents the "week" command
var weekCMD = &cobra.Command{
	Use:     "week",
	Aliases: []string{"w"},
	Short:   "Display events for the current week",
	Long: `Display calendar events for the current week. 
By default, this command only shows upcoming and in-progress events, hiding past ones. 
Use the --all flag to include events that have already ended. 
Supports JSON output, making it easy to pipe into tools like jq for advanced filtering.`,
	Example: `  mycal week                                       # Shows upcoming/in-progress events (hides past)
  mycal week --all                                   # Shows all events (including past)
  mycal week --json | jq ".[] | select(.day == \"Fri\")" # Filter JSON output for Friday
  mycal w`,
	RunE: weekCommand,
}

func weekCommand(cmd *cobra.Command, args []string) error {
	var json, _ = cmd.Flags().GetBool("json")
	var all, _ = cmd.Flags().GetBool("all")

	var weekEvents, err = APIClient.GetWeekEvents(false, all)
	if err != nil {
		return err
	}

	if json {
		formatter.PrintJSON(weekEvents)
	} else {
		formatter.PrettyPrintWeek(weekEvents)
	}

	return nil
}

func init() {
	weekCMD.Flags().Bool("json", false, JsonDescription)
	weekCMD.Flags().Bool("all", false, AllDescription)

	rootCmd.AddCommand(weekCMD)
}
