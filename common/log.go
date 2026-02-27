// Package common provides shared utilities and structures for the teamux project.
package common

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
)

const (
	tempLogFile = "/tmp/teamux.log"
	logDepth    = 2
)

type logger struct {
	infologger    *log.Logger
	warninglogger *log.Logger
	errorlogger   *log.Logger
	fatallogger   *log.Logger
}

var (
	teamuxLogger *logger
	once         sync.Once
)

func GetLogger() *logger {
	once.Do(func() {
		teamuxLogger = createLogger()
	})
	return teamuxLogger
}

func createLogger() *logger {
	logfile, err := os.OpenFile(tempLogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalln("Unable to setup logfile:", err.Error())
	}
	// defer logfile.Close()
	return &logger{
		infologger:    log.New(logfile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warninglogger: log.New(logfile, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorlogger:   log.New(logfile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatallogger:   log.New(logfile, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

func (l *logger) Info(msg string) {
	_ = l.infologger.Output(logDepth, msg)
}

func (l *logger) Warning(msg string) {
	_ = l.warninglogger.Output(logDepth, msg)
}

func (l *logger) Error(msg string) {
	_ = l.errorlogger.Output(logDepth, msg)
}

func (l *logger) Fatal(msg string) {
	_ = l.fatallogger.Output(logDepth, msg)
}

func ShowLogFile(n int) ([]byte, error) {
	var cmd *exec.Cmd
	if n == -1 {
		cmd = exec.Command("cat", tempLogFile)
	} else {
		cmd = exec.Command(
			"sh",
			"-c",
			fmt.Sprintf("tail -n %d %s", n, tempLogFile),
		)
	}
	return cmd.CombinedOutput()
}

func ClearLogFile() error {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("echo -n \"\" > %s", tempLogFile))
	return cmd.Run()
}
