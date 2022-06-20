package main

import (
	"challenge/logger"
	"flag"
	"io/ioutil"
	"time"
	//Import Github van orginele code.
)

var configPath = flag.String("config", "configs/default.json", "configuration file")
var logPath = flag.String("log", "logs/from-"+time.Now().Format("2006-01-02")+".log", "log file")
var address = flag.String("http", ":8080", "address for http server")
var nolog = flag.Bool("nolog", false, "disable logging to file only")
var logfilter = flag.String("logfilter", "", "text to filter log by (both console and file)")

func main() {
	flag.Parse()
	jsonData, err := ioutil.ReadFile(*configPath) //Laad configuratie file van /configs
	if err != nil {
		panic("Fout bij lezen van configuratiefile")
	}

	if *nolog == true {
		logger.Disable()
	}

	if *logfilter != "" {
		logger.Filter(*logfilter)
	}

	logger.SetFilename(*logPath)

	config := NewConfig(jsonData)
	monitor := NewMonitor(config)
	go gossm.RunHttp(*address, monitor)
	gossm.Run() //Start monitoring tool uit bestand
}
