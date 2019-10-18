/*
Cross compile for AWS AMI2 image like so:
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"
*/

package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/datamind-dot-no/kerio-checks/config"

	"github.com/datamind-dot-no/kerio-checks/app/notifications"
	"github.com/datamind-dot-no/kerio-checks/app/qcheck"
)

var (
	// Trace log
	Trace *log.Logger

	// Info level logs
	Info *log.Logger

	// Warning level log
	Warning *log.Logger

	// Error level log
	Error *log.Logger
)

// Init the logging stuff
func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {

	// test some logging levels and messages
	Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	//Trace.Println("I have something standard to say")
	//Info.Println("Special Information")
	//Warning.Println("There is something you need to know about")
	//Error.Println("Something has failed")

	conf, err := config.New()
	if err != nil {
		panic(err)
	}
	notifications := notifications.New(conf)
	qchk := qcheck.New(conf, notifications)

	err = qchk.CheckQ()
	if err != nil {
		panic(err)
	}
}
