package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/opxyc/wdc/alert"
)

// WS is handle for websocket connection
type WS struct {
	c *websocket.Conn
}

// webSocketCon creates a websocket connection to addr on endpoint ep
// and returns a handle. You can also specify a logger which will be used
// to log incoming messages.
func websocketCon(addr, ep string) *WS {
	u := url.URL{Scheme: "ws", Host: addr, Path: ep}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		fmt.Printf("[!] Failed to connect to server: %s\n", err)
		os.Exit(1)
	} else {
		fmt.Printf("[+] Connected to %s\n", u.Host)
	}

	ws := WS{c: c}
	return &ws
}

// lnl listens and logs incoming messages
func (ws *WS) lnl(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, message, err := ws.c.ReadMessage()
			if err != nil {
				fmt.Printf("[-] Could not read alert: %v\n", err)
				if websocket.IsUnexpectedCloseError(err) {
					fmt.Printf("[!] Seems like there is some issue at the server.. Exiting\n")
					return
				}
			}

			var alert alert.Alert
			err = json.Unmarshal(message, &alert)
			if err != nil {
				fmt.Printf("could not unmarshal alert: %v\n", err)
			}
			alog(&alert)
		}
	}
}

// Close closes a ws connection
func (ws *WS) Close() {
	err := ws.c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Printf("could not close connection: %v\n", err)
		return
	}
	log.Println("connection closed")
}

var (
	// to keep track of whether header row was printed to console
	printed bool
	mu      sync.Mutex
)

// printHeader prints log header to the console if not already done
func printHeader() {
	if !printed {
		mu.Lock()
		fmt.Printf("%s\n", alogHeader())
		printed = true
		mu.Unlock()
	}
}

// alogHeader returns alert log header
func alogHeader() string {
	// header format:
	// 			 TIME		ID            Host          Message
	titleRow := fmt.Sprintf("%-9v %-23s %-16s %s", "TIME", "ID", "Host", "Message")
	return titleRow
}

// alog logs alert to console and file
func alog(msg *alert.Alert) {
	const maxMsgLen = 70
	// logs received msg in the format:
	// TIME		ID            Host          Message
	// 13:00:10	1635149439253 srv01         cpu usage on > 10%. take action immediately
	// Note: the detailed info is not logged to console. it it stored to log file only.

	printHeader()
	// take time and neglect date for logging to console
	timeOnly := strings.Split(msg.Time, " ")[1]
	info := fmt.Sprintf("%-9v %-23s %-16s %s", timeOnly, msg.ID, msg.From, msg.Short)
	fmt.Printf("%s\n", info)

	// see pkg alert/alert.go to learn about the format in which an alret is logged to file
	l.Printf("\n%s\n%s\n%s\n%s\n%s\n%s%d\nENDOF%s\n\n", msg.ID, msg.Time, msg.From, msg.TaskName, msg.Short, msg.Long, msg.Status, msg.ID)
}
