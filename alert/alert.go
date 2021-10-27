package alert

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	// get today's date
	// open corresponding log file and read
	// if not present, read yesterdays... and continue until past 60 days..

	var t time.Time
	now := time.Now()
	t = time.Now()
	var alertLines []string // to store the lines read from logfile
	var found bool          // to indicate if the log id is found
	for i := 0; i < 60 && !found; i++ {
		// get time
		t = now.AddDate(0, 0, -1*i)
		// get log file name from time
		fname := filepath.Join(dir, t.Format(format))
		if _, err := os.Stat(fname); err != nil {
			continue
		}

		f, err := os.Open(fname)
		if err != nil {
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			if line == id {
				alertLines = append(alertLines, line)
				found = true
			} else if found {
				if line == fmt.Sprintf("ENDOF%s", id) {
					break
				}
				alertLines = append(alertLines, line)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	if !found {
		return nil, fmt.Errorf("log with id '%v' not found", id)
	}

	return parseAlert(alertLines), nil
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
