package main

import (
	"os"
	"os/signal"
	"syscall"
)

func setupSignalHandlers() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		for {
			select {
			case signal := <-signalChan:
				cui.ln("Received", signal, "- exiting...")
				os.Exit(0)
			}
		}
	}()
}

func fatalWithErr(desc string, err error) {
	cui.ln(desc)
	cui.ln("Reason:", err)
	cui.ln("Bye...")
	cui.ln("")
	os.Exit(1)
}

func fatal(desc string) {
	cui.ln(desc)
	cui.ln("Bye...")
	cui.ln("")
	os.Exit(1)
}
