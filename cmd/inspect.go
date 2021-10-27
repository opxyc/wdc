package cmd

import (
	"fmt"

	"github.com/opxyc/wdc/alert"
	"github.com/spf13/cobra"
)

// inspectCmd represents the inspect command
var inspectCmd = &cobra.Command{
	Use:           "inspect ID",
	Args:          cobra.ExactArgs(1),
	Short:         "Inspect an alert using it's ID",
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}

		id := args[0]

		alert, err := alert.ReadFromLog(logDir, id)
		if err != nil {
			fmt.Printf("Could not find info on given log id: %v\n", err)
			return
		}

		fmt.Printf("%-14s %s\n", "ID", alert.ID)
		fmt.Printf("%-14s %s\n", "Time", alert.Time)
		fmt.Printf("%-14s %s\n", "Hostname", alert.From)
		fmt.Printf("%-14s %s\n", "Task Name", alert.TaskName)
		fmt.Printf("%-14s %s\n", "Message", alert.Short)
		var status string
		if alert.Status == 0 {
			status = "OK (actions executed successfully)"
		} else {
			status = "Require manual effort"
		}
		fmt.Printf("%-14s %s\n", "Status", status)
		fmt.Printf("\nCommand output:\n%s\n", alert.Long)
	},
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}
