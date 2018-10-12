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
	Debug = log.New(debugHandle, "SCOPE - DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Info = log.New(infoHandle, "SCOPE - INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Warning = log.New(warningHandle, "SCOPE - WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(errorHandle, "SCOPE - ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
	Fatal = log.New(fatalHandle, "SCOPE - FATAL: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}
