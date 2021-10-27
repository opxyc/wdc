# WatchDogClient - WDC
WDC is a client application for [WatchDog - wc](https://github.com/opxyc/wd), which is a tool to monitor and health check servers through plug-in scripts.

WDC listens to a remote alert server which broadcasts alert messages. The received alert will be logged to `WDC/logs/` in user's home directory. It also comes with a local http server listening on port 8080 with an endpoint `/{id}` that can be used to get details of an alert given it's id. The same can be done via command line using the `inspect` command of wdc.

```
Usage:
  wdc [command]

Available Commands:
  help        Help about any command
  inspect     Inspect an alert using it's ID
  listen      Start listening to incoming alerts
```

####  Start listening for alerts
```
$ wdc listen
```
![](https://user-images.githubusercontent.com/34760059/139127203-80a52bc7-3ad5-407a-b0e7-8debeb56a5bc.png)

But default, wdc connects to `localhost:40080` at endpoint `/ws/connect`.
```
Usage:
  wdc listen [flags]

Flags:
  -r, --addr string        network address of alert server (default "localhost:40080")
  -e, --end-point string   websocket connection endpoint of alert server (default "/ws/connect")
```

#### Inspecting an alert
```
$ wdc inspect ID
```
```
be@mC:~/Desktop/wdc$ wdc inspect PgKt664SLUxceGiuFhTGMn
ID             PgKt664SLUxceGiuFhTGMn
Time           2021-Oct-27 23:31:03
Hostname       dr-pg-l02
Task Name      mount-point-arc-check
Message        Mount point /arc > 90%
Status         OK (actions executed successfully) 

Command output:
(mount-point-arc-check.out) current utilization 91.2%
(mount-point-arc-check.delete-archives.out) deleting archives
(mount-point-arc-check.delete-archives.out) current utilization 63.8%
```

You can also get the same info by visiting `localhost:8080/{id}`.