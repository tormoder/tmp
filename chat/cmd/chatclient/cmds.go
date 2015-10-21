package main

import (
	"bytes"
	"fmt"
	"time"
)

const (
	listAllUsers = iota
	sendPublicMessage
	sendPrivateMessage
	logout
)

var cmdTypes = [...]string{
	"List all online users",
	"Send a public message",
	"Send a private message",
	"Logout",
}

var cmdTypesShort = [...]string{
	"List users",
	"Send public",
	"Send private",
	"Logout",
}

var cmdsShort string

func init() {
	var output bytes.Buffer
	output.WriteString("[menu] ")
	for i, userOption := range cmdTypesShort {
		output.WriteString(
			fmt.Sprintf("%d. %s ", i+1, userOption),
		)
	}
	cmdsShort = output.String()
}

func getAvailableCmds() string {
	var out bytes.Buffer
	out.WriteString("Available commands:\n")
	for i, userOption := range cmdTypes {
		out.WriteString(
			fmt.Sprintf("\t%d. %s\n", i+1, userOption),
		)
	}
	return out.String()
}

func getAvailableCmdsShort() string {
	var out bytes.Buffer
	out.WriteString(time.Now().Format(tformat))
	out.WriteString(" ")
	out.WriteString(cmdsShort)
	return out.String()
}
