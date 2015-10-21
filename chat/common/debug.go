package common

import "log"

func SetVerbose() {
	verbose = true
}

var verbose bool

func Debug(v ...interface{}) {
	if verbose {
		log.Print(v...)
	}
}

func Debugf(format string, v ...interface{}) {
	if verbose {
		log.Printf(format, v...)
	}
}

func Debugln(v ...interface{}) {
	if verbose {
		log.Println(v...)
	}
}
