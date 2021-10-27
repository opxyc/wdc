package alert

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Alert msg format
type Alert struct {
	Time     string `json:"time"`
	ID       string `json:"id"`
	From     string `json:"from"`
	TaskName string `json:"taskName"`
	Short    string `json:"short"`  // short message - msg field in client config.json
	Long     string `json:"long"`   // long message - combined output of `cmd`
	Status   int32  `json:"status"` // 0 if success, 1 if failed
}

// ---
// alerts are logged in the below format: (ex)
// 		    2021/10/26 18:37:45
// line 1)  rJeqKsiMECBUpCGG54YhgL								👈 log id
// line 2)  2021-Oct-27 13:40:04								👈 the time alert was generated
// line 3)  mC													👈 hostname
// line 4)  cpu-usage-gt-10										👈 task name
// line 5)  cpu usage on > 10%. take action immediately			👈 short msg
// line 6)  (cpu-usage-gt-10.err) exit status 1					👈 long msg (combined output of the
// line 7)  (cpu-usage-gt-10.out) current CPU usage is 15	  	   cmd executed at client)
// line 8)  (cpu-usage-gt-10.bringing foo down.err) exec: "lss":
// line 9)  executable file not found in $PATH
// line 10) 0
// 			ENDOFrJeqKsiMECBUpCGG54YhgL							👈 end of one alert

// log file name format
const format = "2006-Jan-02"

// ReadFromLog reads the alert with id from log files in dir
func ReadFromLog(dir, id string) (*Alert, error) {
	// search for given alert id in the logs of part 30 days
	const n = 30
	var alert *Alert
	var wg sync.WaitGroup

	// today's time is taken as reference and will subtract one day from it
	// to get previous n days log file names
	referenceTime := time.Now()
	var t time.Time
	t = time.Now()

	// ctx for cancelling file reads that are in progress if some other goroutine
	// found the alert we are looking for
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()

			// get log file name from time
			t = referenceTime.AddDate(0, 0, -1*i)
			fname := filepath.Join(dir, t.Format(format))
			if _, err := os.Stat(fname); err != nil {
				return
			}

			alrt, _ := searchAndRetrieve(ctx, id, fname)
			// if we got the alert, cancel all other reads
			if alrt != nil {
				alert = alrt
				cancelFunc()
			}
		}(i)
	}

	// wait for all goroutines to complete/return
	wg.Wait()

	if alert == nil {
		return nil, fmt.Errorf("log with id '%v' not found", id)
	}

	return alert, nil
}

// searchAndRetrieve searches the file with name fname for an alert with given id
// until it's found, EOF reached or ctx is cancelled
func searchAndRetrieve(ctx context.Context, id, fname string) (*Alert, error) {
	var alertLines []string // to store the lines read from logfile
	var found bool          // to indicate if the log id is found

	f, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("received done signal")
	default:
		for scanner.Scan() {
			line := scanner.Text()
			if line == id {
				alertLines = append(alertLines, line)
				found = true
			} else if found {
				if line == fmt.Sprintf("ENDOF%s", id) {
					return parseAlert(alertLines), nil
				}
				alertLines = append(alertLines, line)
			}
		}
	}

	return nil, fmt.Errorf("alert for id '%v' not found in '%v'", id, fname)
}

func parseAlert(info []string) *Alert {
	alert := &Alert{
		ID:       info[0],
		Time:     info[1],
		From:     info[2],
		TaskName: info[3],
		Short:    info[4],
	}

	len := len(info)
	// last item in the slice is status
	status, _ := strconv.ParseInt(info[len-1], 10, 32)
	alert.Status = int32(status)
	// 4th line till second last is long msg
	alert.Long = strings.Join(info[5:len-1], "\n")

	return alert
}
