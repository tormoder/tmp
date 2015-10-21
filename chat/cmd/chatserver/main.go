package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/tormoder/chat/chat"
	c "github.com/tormoder/chat/common"
	pb "github.com/tormoder/chat/proto"
	"github.com/tormoder/chat/storage"
	"github.com/tormoder/chat/user"

	"google.golang.org/grpc"
)

var (
	port    = flag.Int("port", 10000, "The chat server `port`")
	verbose = flag.Bool("v", false, "show verbose debugging output")
)

func main() {
	flag.Parse()
	log.SetPrefix("[chatserver] ")
	if *verbose {
		c.SetVerbose()
	}

	listener, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", *port),
	)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		for {
			select {
			case signal := <-signalChan:
				log.Println("Received", signal, "- exiting...")
				os.Exit(0)
			}
		}
	}()

	grpcServer := grpc.NewServer()

	c.Debugln("setting up storage, chat and user service")
	userStorage := storage.NewInMemoryUserStorage()
	chatService := chat.NewService(userStorage)
	userService := user.NewService(chatService, userStorage)

	c.Debugln("registering services with grpc")
	pb.RegisterUserServiceServer(grpcServer, userService)
	pb.RegisterChatServiceServer(grpcServer, chatService)

	c.Debugln("listening on", listener.Addr())
	grpcServer.Serve(listener)
}
