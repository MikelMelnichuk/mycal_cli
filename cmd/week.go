package cmd

import (
	"github.com/MikelMelnichuk/mycal/internal/formatter"
	"github.com/spf13/cobra"
)

// mycal week       # Shows the events that are yet to happen (or in progress), without the past events
// mycal week --all # Shows all the events events of the current weak (including past)
// mycal week --json | jq ".[] | select(.day == \"Fri\")"
// mycal w

var weekCMD = &cobra.Command{
	Use:     "week",
	Aliases: []string{"w"},
	RunE: func(cmd *cobra.Command, args []string) error {
		var json, _ = cmd.Flags().GetBool("json")
		var all, _ = cmd.Flags().GetBool("all")

		var weekEvents, err = APIClient.GetWeekEvents(false, all)
		if err != nil {
			return err
		}

		if json {
			formatter.PrintJSON(weekEvents)
		} else {
			formatter.PrintText(weekEvents)
		}

		return nil
	},
}

func init() {
	weekCMD.Flags().Bool("json", false, JsonDescription)
	weekCMD.Flags().Bool("all", false, AllDescription)

	rootCmd.AddCommand(weekCMD)
}
