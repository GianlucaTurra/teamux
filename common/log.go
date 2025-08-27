// Package common provides shared utilities and structures for the teamux project.
package common

import (
	"log"
	"os"
)

type Logger struct {
	Infologger    *log.Logger
	Warninglogger *log.Logger
	Errorlogger   *log.Logger
	Fatallogger   *log.Logger
}

func GetLogger(file string) Logger {
	logfile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("Unable to setup syslog:", err.Error())
	}
	defer logfile.Close()
	return Logger{
		Infologger:    log.New(logfile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warninglogger: log.New(logfile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		Errorlogger:   log.New(logfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		Fatallogger:   log.New(logfile, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
