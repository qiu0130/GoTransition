package transitions

import (
	"log"
	"fmt"
	"os"
)
const (
	VERSION = "version 0.1.0"
)

var (
	logger *log.Logger
)


type HandleFunc func(ed *EventData)

type ConditionFunc func(ed *EventData) bool


func Info(format string, v...interface{}) {
	logger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Output(2, fmt.Sprintf(format, v...))

}

func Warning(format string, v...interface{}) {
	logger = log.New(os.Stdout, "WARNING", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Output(2, fmt.Sprintf(format, v...))

}

func Debug(format string, v... interface{}) {
	logger = log.New(os.Stdout, "DEBUG", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Output(2, fmt.Sprintf(format, v...))
}

func Error(format string, v...interface{}) {
	logger = log.New(os.Stderr, "ERROR", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	logger.Output(2, fmt.Sprintf(format, v...))

}
