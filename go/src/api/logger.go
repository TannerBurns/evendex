package main

import (
	"io"
	"log"
)

var (
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Fatal   *log.Logger
)

func initLogging(
	debugHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	fatalHandle io.Writer,
) {
	Debug = log.New(debugHandle, "EVENDEX - DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "EVENDEX - INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "EVENDEX - WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "EVENDEX - ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(fatalHandle, "EVENDEX - FATAL: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
