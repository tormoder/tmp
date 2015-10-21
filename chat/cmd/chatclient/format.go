package main

import (
	"bytes"
	"fmt"
	"time"

	pb "github.com/tormoder/chat/proto"
)

const tformat = "15:04:05"

func formatUnixTime(t int64) string {
	return time.Unix(t, 0).Format(tformat)
}

func formatUserList(users []*pb.User) string {
	var output bytes.Buffer
	output.WriteString(time.Now().Format(tformat))
	output.WriteString(" [info] ")
	switch {
	case users == nil, len(users) == 0:
		output.WriteString("No current logged-in users")
	case len(users) == 1:
		output.WriteString("One logged-in user:\n")
		output.WriteString(
			fmt.Sprintf(
				"\t1. %s (last seen %s)",
				users[0].Nick,
				formatUnixTime(users[0].TimeLastSeen),
			),
		)
	default:
		output.WriteString(
			fmt.Sprintf("%d logged-in users:", len(users)),
		)
		for i, user := range users {
			output.WriteString(
				fmt.Sprintf(
					"\n\t%d. %s\t(last seen %s)",
					i+1,
					user.Nick,
					formatUnixTime(user.TimeLastSeen),
				),
			)
		}
	}

	return output.String()
}

func formatMsg(msg *pb.ChatServerMsg) string {
	var output bytes.Buffer
	switch msg.Msg.(type) {
	case *pb.ChatServerMsg_PublicMsg:
		pmsg := msg.GetPublicMsg()
		output.WriteString(
			fmt.Sprintf(
				"%s [%s] %s",
				formatUnixTime(pmsg.TimeSent),
				pmsg.GetFrom().Nick,
				pmsg.Msg,
			),
		)
	case *pb.ChatServerMsg_PrivateMsg:
		pmsg := msg.GetPrivateMsg()
		output.WriteString(
			fmt.Sprintf(
				"%s [%s] [private] %s",
				formatUnixTime(pmsg.TimeSent),
				pmsg.GetFrom().Nick,
				pmsg.Msg,
			),
		)
	case *pb.ChatServerMsg_UserEvent:
		uevent := msg.GetUserEvent()
		output.WriteString(
			fmt.Sprintf(
				"%s [info] User %s just ",
				formatUnixTime(uevent.Time),
				uevent.GetUser().Nick,
			),
		)
		switch uevent.Event {
		case pb.UserEvent_LOGIN:
			output.WriteString("logged-in. ")
		case pb.UserEvent_LOGOUT:
			output.WriteString("logged-out. ")
		default:
			output.WriteString("did somthing unknown. ")
		}
		output.WriteString("Last seen ")
		output.WriteString(formatUnixTime(uevent.GetUser().TimeLastSeen))
	default:
		output.WriteString("Unkown type of message received from chat server")
	}

	return output.String()
}
