package transitions

import (
	"fmt"
	"log"
	"os"
)

const (
	VERSION = "0.2.0"
)

var (
	logger *log.Logger
	Debug bool
)

func init () {
	Debug = true
	Info("the latest  of version is %s", VERSION)
}

type HandleFunc func(ed *EventData)

type ConditionFunc func(ed *EventData) bool

func Info(format string, v ...interface{}) {
	if Debug {
		logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
		logger.Output(2, fmt.Sprintf(format, v...))
	}
}

func Error(format string, v ...interface{}) {
	logger = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Output(2, fmt.Sprintf(format, v...))

}