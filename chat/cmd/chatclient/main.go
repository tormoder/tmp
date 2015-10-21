package main

import (
	"flag"
	"os"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	pb "github.com/tormoder/chat/proto"
)

var serverAddr = flag.String("saddr", "127.0.0.1:10000", "The chat server address in the format of host:port")

var (
	cui         = ui{os.Stdout}
	tocuiChan   = make(chan string, 2048)
	userService pb.UserServiceClient
	chatService pb.ChatServiceClient
	credentials *pb.Credentials
)

func main() {
	flag.Parse()
	setupSignalHandlers()

	cui.ln("---------------------------------")
	cui.ln("\tSimple Chat Client")
	cui.ln("---------------------------------")

	nick := cui.promptForString("nick")

	cui.ln("Dialing chat server...")
	err := dialServer()
	if err != nil {
		fatalWithErr("Dialing chat server failed:", err)
	}

	cui.ln("Attempting to login...")
	credentials, err = attemptLogin(nick)
	if err != nil {
		fatalWithErr("Login failed", err)
	}

	cui.f("Hello %s, login success!\n", nick)

	cui.ln("Setting up message listener...")
	err = setupMsgListener()
	if err != nil {
		fatalWithErr("Message listener setup failed", err)
	}
	cui.ln("Ready")
	cui.ln()
	cui.ln(getAvailableCmds())
	cui.ln("")

	clientLoop()
}

func dialServer() error {
	clientConn, err := grpc.Dial(
		*serverAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(500*time.Millisecond),
	)
	if err != nil {
		return err
	}
	userService = pb.NewUserServiceClient(clientConn)
	chatService = pb.NewChatServiceClient(clientConn)
	return nil
}

func attemptLogin(nick string) (*pb.Credentials, error) {
	lreq := &pb.LoginRequest{Nick: nick}
	creds, err := userService.Login(context.Background(), lreq)
	if err != nil {
		return nil, err
	}
	return creds, nil
}

func setupMsgListener() error {
	msgStream, err := chatService.ListenForMessages(context.Background(), credentials)
	if err != nil {
		return err
	}
	go func() {
		for {
			msg, err := msgStream.Recv()
			if err != nil {
				fatal("chat server connection error")
			}
			if heartbeat(msg) {
				continue
			}
			select {
			case tocuiChan <- formatMsg(msg):
				// Send OK
			default:
				// UI queue full, drop message
			}
		}
	}()

	return nil
}

func heartbeat(msg *pb.ChatServerMsg) bool {
	if _, ok := msg.Msg.(*pb.ChatServerMsg_Heartbeat); ok {
		return true
	}
	return false
}

func clientLoop() {
	var (
		userChoice      int
		stopRefreshChan = make(chan bool)
	)

	for {
		cui.ln(getAvailableCmdsShort())
		go func() {
			for {
				select {
				case <-time.After(300 * time.Millisecond):
					msgsWasPrinted := pumpNewMsgToUI()
					if msgsWasPrinted {
						cui.ln(getAvailableCmdsShort())
					}
				case <-stopRefreshChan:
					return
				}
			}

		}()
		cui.sln(&userChoice)
		stopRefreshChan <- true

		switch userChoice - 1 {
		case listAllUsers:
			printAllUsers()
			pumpNewMsgToUI()
		case sendPublicMessage:
			sendPublicMsg()
			pumpNewMsgToUI()
		case sendPrivateMessage:
			sendPrivateMsg()
			pumpNewMsgToUI()
		case logout:
			attemptLogout()
			os.Exit(0)
		default:
			pumpNewMsgToUI()
		}
	}
}

func pumpNewMsgToUI() bool {
	hasPrintedMsg := false
	for {
		select {
		case msg := <-tocuiChan:
			cui.ln(msg)
			hasPrintedMsg = true
		default:
			return hasPrintedMsg
		}
	}
}

func printAllUsers() {
	luresp, err := userService.ListUsers(context.Background(), credentials)
	if err != nil {
		cui.ln("Unable to list users:", err)
		return
	}
	cui.ln(formatUserList(luresp.Users))
}

func sendPublicMsg() {
	pmsg := cui.promptForString("public message")
	msg := new(pb.PublicMsgRequest)
	msg.Msg = pmsg
	msg.Creds = credentials
	_, err := chatService.SendPublic(context.Background(), msg)
	if err != nil {
		cui.ln("Error sending message:", err)
	}
}

func sendPrivateMsg() {
	rnick := cui.promptForString("nick of message receiver")
	pmsg := cui.promptForString("private message")
	msg := new(pb.PrivateMsgRequest)
	msg.To = rnick
	msg.Msg = pmsg
	msg.Creds = credentials
	_, err := chatService.SendPrivate(context.Background(), msg)
	if err != nil {
		cui.ln("Error sending message:", err)
	}
}

func attemptLogout() {
	_, err := userService.Logout(context.Background(), credentials)
	if err != nil {
		fatalWithErr("Error on logout", err)
	}
	cui.ln("Logout successful")
}
