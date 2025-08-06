package internal

import "log"

type Logger struct {
	Infologger    *log.Logger
	Warninglogger *log.Logger
	Errorlogger   *log.Logger
	Fatallogger   *log.Logger
}
