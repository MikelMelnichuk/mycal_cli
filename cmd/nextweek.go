package cmd

import (
	"github.com/MikelMelnichuk/mycal/internal/formatter"
	"github.com/spf13/cobra"
)

// mycal nextweek
// mycal nextweek --json
// mycal nw

var nextWeekCMD = &cobra.Command{
	Use:     "nextweek",
	Aliases: []string{"nw", "next-week"},
	RunE: func(cmd *cobra.Command, args []string) error {
		json, _ := cmd.Flags().GetBool("json")

		events, err := APIClient.GetWeekEvents(true, true)
		if err != nil {
			return nil
		}

		if json {
			formatter.PrintJSON(events)
		} else {
			formatter.PrintText(events)
		}

		return nil
	},
}

func init() {
	nextWeekCMD.Flags().Bool("json", false, JsonDescription)

	rootCmd.AddCommand(nextWeekCMD)
}
