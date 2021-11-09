package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/opxyc/goutils/logger"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "Start listening to incoming alerts",
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := cmd.Flags().GetString("addr")
		if err != nil {
			cmd.Usage()
			return
		}
		ep, err := cmd.Flags().GetString("end-point")
		if err != nil {
			cmd.Usage()
			return
		}
		startListener(addr, ep)
		go runHTTPServer()
	},
}

func init() {
	listenCmd.Flags().StringP("addr", "r", "localhost:40080", "network address of alert server")
	listenCmd.Flags().StringP("end-point", "e", "/ws/connect", "websocket connection endpoint of alert server")
	rootCmd.AddCommand(listenCmd)
}

func startListener(addr, ep string) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var err error
	const logFileNameFormat = "2006-Jan-02"
	l, err = logger.NewDailyLogger(ctx, logDir, logFileNameFormat, 00, 00)
	if err != nil {
		fmt.Printf("Could not initialize logger: %v. Exiting..\n", err)
		return
	}

	// make websocket connection
	ws := websocketCon(addr, ep)
	// listen and log incoming messages
	go func() {
		ws.lnl(ctx)
		// if the listener returned, it means something went wrong
		// and cannot continue. so cancel the context and exit..
		cancelFunc()
		os.Exit(1)
	}()

	// start http server
	go runHTTPServer()

	// wait for interrupt (if any)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	<-sigchan
	fmt.Println("[+] Saving logs..")
	// cancel context which will close ws and log file
	cancelFunc()
	time.Sleep(time.Millisecond * 300)
	fmt.Println("[+] Done")
}
