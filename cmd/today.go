package cmd

import (
	"fmt"
	"time"

	"github.com/MikelMelnichuk/mycal/internal/formatter"
	"github.com/spf13/cobra"
)

// todayCmd represents the "today" command
var todayCmd = &cobra.Command{
	Use:     "today",
	Aliases: []string{"t"},
	Short:   "Display events for today",
	Long: `Display calendar events scheduled for today. 
By default, this command shows upcoming events for the current day. 
Use the --all flag to include past events, or filter by time using --after. 
Supports outputting the results in JSON format for scripting and automation.`,
	Example: `  mycal today
  mycal today --all
  mycal today --json
  mycal today --after 13:00
  mycal t`,
	RunE: todayCommand,
}

func todayCommand(cmd *cobra.Command, args []string) error {
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
}

func init() {
	// Define expected flags for todayCmd
	todayCmd.Flags().Bool("all", false, AllDescription)
	todayCmd.Flags().Bool("json", false, JsonDescription)
	todayCmd.Flags().String("after", "", AfterDescription)

	// Register the command with rootCmd
	rootCmd.AddCommand(todayCmd)
}
