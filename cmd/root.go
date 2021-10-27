package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	l      *log.Logger // logger for logging alerts to logfile
	logDir string      // log file directory
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wdc",
	Short: "wdc listens for incoming alert messages and help you inspect the same.",
	Long: `wdc listens for incoming alert messages and help you inspect the same.
An alert can be inspected using the inspect command or by visiting localhost:8080/{id}`,
}

func init() {
	// set the log directory
	d, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("[!] Could not get your home directory: %v. Exiting..\n", err)
		os.Exit(1)
	}

	logDir = filepath.Join(d, "WDC", "logs")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
