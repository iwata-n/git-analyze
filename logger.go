package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
)

func InitLog(isDebug bool, isShowProgress bool) (*log.Logger, *log.Logger) {
	ilog := log.New(io.Discard, "", log.Ltime|log.Lmicroseconds|log.LUTC)
	dlog := log.New(io.Discard, "[DEBUG]", log.Ltime|log.Lmicroseconds|log.LUTC)
	if isShowProgress {
		ilog.SetOutput(os.Stdout)
	}

	if isDebug {
		dlog.SetOutput(os.Stdout)
		ilog.SetOutput(os.Stdout)
		ilog.SetPrefix("[INFO ]")
	}

	return dlog, ilog
}

func JsonString(a ...interface{}) string {
	var buf bytes.Buffer
	m, err := json.Marshal(a)
	checkIfError(err)

	err = json.Indent(&buf, m, "", "  ")

	return buf.String()
}
