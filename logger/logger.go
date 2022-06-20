package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	enabledFileLog = true
	logFilename    = "log.txt"
	mu             sync.Mutex
	filter         string
)

// Functie zet log om in platte tekst en over naar een bestand
func Log(text string) {
	if filter == "" || (filter != "" && strings.Contains(text, filter)) {
		log.Print(text)

		if !enabledFileLog {
			return
		}
		mu.Lock()
		defer mu.Unlock()
		if err := writeToFile(logFilename, text); err != nil {
			log.Println(err)
		}
	}
}

// Schrijft nieuwe log op een nieuwe line in de logfile
func Logln(v ...interface{}) {
	Log(fmt.Sprintln(v...))
}

// Functie schrijft text naar standaard output
func Logf(format string, v ...interface{}) {
	Log(fmt.Sprintf(format, v...))
}

// Update filename
func SetFilename(fileName string) {
	logFilename = fileName
}

// Schakelt logging uit
func Disable() {
	enabledFileLog = false
}

// Schakelt loggin in
func Enable() {
	enabledFileLog = true
}

// Filter op speficieke keywords
func Filter(f string) {
	filter = f
}

// Schrijf log naar logfile en bestaande bestandsnaam, als hij niet bestaat maakt hij deze aan
func writeToFile(fileName, text string) error {
	file, err := os.OpenFile(logFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()
	text = fmt.Sprintf("%s %s", time.Now().String(), text)
	if _, err = file.WriteString(text); err != nil {
		return err
	}
	return nil
}
