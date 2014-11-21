package main

import (
	"io"
	"log"
	"os"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Debug   *log.Logger
)

func initlog(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer,
	debugHandle io.Writer) {

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

	Debug = log.New(debugHandle,
		"DEBUG: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func InitLogging(t, i, w, e, d bool) {
	fileTrace := os.Stdout
	if t == true {
		fileTrace, _ = os.OpenFile("trace.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}

	fileInfo := os.Stdout
	if i == true {
		fileInfo, _ = os.OpenFile("info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}

	fileWarning := os.Stdout
	if w == true {
		fileWarning, _ = os.OpenFile("warning.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}

	fileError := os.Stdout
	if e == true {
		fileError, _ = os.OpenFile("error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}

	fileDebug := os.Stdout
	if d == true {
		fileDebug, _ = os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}

	initlog(fileTrace, fileInfo, fileWarning, fileError, fileDebug)
}
