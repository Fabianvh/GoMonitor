package main

import (
	"challenge/logger"
	"flag"
	"io/ioutil"
	"time"

	"github.com/ssimunic/gossm"
	//Import Github van orginele code.
)

var configPath = flag.String("config", "configs/default.json", "configuration file")              //Geeft pad naar configfile aan
var logPath = flag.String("log", "logs/from-"+time.Now().Format("2006-01-02")+".log", "log file") //Geeft pad aan naar logfile
var address = flag.String("http", ":8080", "address for http server")                             //Geeft de poort aan waar de webserver op gaat draaien
var nolog = flag.Bool("nolog", false, "disable logging to file only")                             //Als de functie true word aangezet word logging uitgezet
var logfilter = flag.String("logfilter", "", "text to filter log by (both console and file)")     //Maakt text naar een logfile

func main() {
	flag.Parse()
	jsonData, err := ioutil.ReadFile(*configPath) //Laad configuratie file van /configs
	if err != nil {
		panic("Fout bij lezen van configuratiefile") //Foutmelding message bij niet kunnen laden van de configfile
	}

	if *nolog == true { //Als de logfile uitstaat, start hij functie logger disable
		logger.Disable()
	}

	if *logfilter != "" { //Logfilter
		logger.Filter(*logfilter)
	}

	logger.SetFilename(*logPath) //Zet bestandnaam van de logfile

	config := NewConfig(jsonData)       //Laad configuratiefile in functie Newconfig
	monitor := NewMonitor(config)       //Laad monitor in NewMonitor
	go gossm.RunHttp(*address, monitor) //Zet HTTP server aan met address en monitoring functie
	gossm.Run()                         //Start monitoring tool
}
