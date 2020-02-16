package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var shouldLog = false

func init() {
	flag.BoolVar(&shouldLog, "log", false, "show internal actions")
	log.SetOutput(os.Stderr)
}

func info(msg string) {
	if !shouldLog {
		return
	}
	log.Println("[info]", msg)
}

func infof(fmtString string, args ...interface{}) {
	if !shouldLog {
		return
	}
	log.Println(fmt.Sprintf("[info] "+fmtString, args...))
}

func debug(msg string) {
	if !shouldLog {
		return
	}
	log.Println("[debug]", msg)
}

func debugf(fmtString string, args ...interface{}) {
	if !shouldLog {
		return
	}
	log.Println(fmt.Sprintf("[debug] "+fmtString, args...))
}

func trace(msg string) {
	if !shouldLog {
		return
	}
	log.Println("[trace]", msg)
}

func tracef(fmtString string, args ...interface{}) {
	if !shouldLog {
		return
	}
	log.Println(fmt.Sprintf("[info] "+fmtString, args...))
}
