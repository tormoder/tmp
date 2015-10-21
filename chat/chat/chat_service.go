package chat

import (
	"errors"
	"sync"
	"time"

	c "github.com/tormoder/chat/common"
	pb "github.com/tormoder/chat/proto"
	"github.com/tormoder/chat/storage"

	"golang.org/x/net/context"
)

type Service struct {
	ustorage storage.UserStorage

	connectedClients map[string]chan *pb.ChatServerMsg
	mu               sync.Mutex // Protects connectedClients
}

func NewService(userStorage storage.UserStorage) *Service {
	return &Service{
		ustorage:         userStorage,
		connectedClients: make(map[string]chan *pb.ChatServerMsg),
	}
}

func (s *Service) BroadcastAllConnectedClients(msg *pb.ChatServerMsg) {
	go func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		for _, clientChan := range s.connectedClients {
			select {
			case clientChan <- msg:
				// Send OK
			default:
				// Client queue full, drop message
			}
		}
	}()
}

func (s *Service) SendPrivate(ctx context.Context, privMsgReq *pb.PrivateMsgRequest) (*pb.SendMsgResponse, error) {
	c.Debugln("send public message request from", privMsgReq.GetCreds().Nick)
	user, err := s.ustorage.CheckCredentials(privMsgReq.GetCreds())
	if err != nil {
		return nil, err
	}

	clientCh, found := s.connectedClients[privMsgReq.To]
	if found {
		select {
		case clientCh <- &pb.ChatServerMsg{
			Msg: &pb.ChatServerMsg_PrivateMsg{
				PrivateMsg: &pb.PrivateMsg{
					To:       privMsgReq.To,
					From:     &user.User,
					Msg:      privMsgReq.Msg,
					TimeSent: time.Now().Unix(),
				},
			},
		}:
			// Send OK
		default:
			// Client queue full, drop message
		}
		return &pb.SendMsgResponse{}, nil
	}

	u, found := s.ustorage.GetUser(privMsgReq.To)
	if !found {
		return nil, errors.New("requested user not found")
	}
	select {
	case u.MsgChannel <- &pb.ChatServerMsg{
		Msg: &pb.ChatServerMsg_PrivateMsg{
			PrivateMsg: &pb.PrivateMsg{
				To:       privMsgReq.To,
				From:     &user.User,
				Msg:      privMsgReq.Msg,
				TimeSent: time.Now().Unix(),
			},
		},
	}:
		// Send OK
	default:
		// Client queue full, drop message
	}

	return &pb.SendMsgResponse{}, nil
}

func (s *Service) SendPublic(ctx context.Context, pubMsgReq *pb.PublicMsgRequest) (*pb.SendMsgResponse, error) {
	c.Debugln("send public message request from", pubMsgReq.GetCreds().Nick)
	user, err := s.ustorage.CheckCredentials(pubMsgReq.GetCreds())
	if err != nil {
		return nil, err
	}

	s.BroadcastAllConnectedClients(
		&pb.ChatServerMsg{
			Msg: &pb.ChatServerMsg_PublicMsg{
				PublicMsg: &pb.PublicMsg{
					From:     &user.User,
					Msg:      pubMsgReq.Msg,
					TimeSent: time.Now().Unix(),
				},
			},
		})

	return &pb.SendMsgResponse{}, nil
}

func (s *Service) ListenForMessages(creds *pb.Credentials, stream pb.ChatService_ListenForMessagesServer) error {
	c.Debugln("listen for messages request from", creds.Nick)
	user, err := s.ustorage.CheckCredentials(creds)
	if err != nil {
		return err
	}

	user.Online = true
	msgChan := user.MsgChannel
	err = s.ustorage.UpdateUser(user)
	if err != nil {
		return c.InternalServerError("storage error")
	}

	s.addClientChanToConnected(creds.Nick, msgChan)

	hbTicker := time.NewTicker(time.Second)
	hb := &pb.ChatServerMsg{
		Msg: &pb.ChatServerMsg_Heartbeat{
			Heartbeat: &pb.Heartbeat{},
		},
	}

	c.Debugln("serving messages for", creds.Nick)

	for {
		select {
		case msg := <-msgChan:
			err = stream.Send(msg)
			if err != nil {
				goto cleanup
			}
		case <-hbTicker.C:
			err = stream.Send(hb)
			if err != nil {
				goto cleanup
			}
		}
	}

cleanup:
	s.rmClientChanFromConnected(creds.Nick)
	user, found := s.ustorage.GetUser(creds.Nick)
	if !found {
		return err
	}
	user.Online = false
	user.TimeLastSeen = time.Now().Unix()
	s.ustorage.UpdateUser(user)

	s.BroadcastAllConnectedClients(
		&pb.ChatServerMsg{
			Msg: &pb.ChatServerMsg_UserEvent{
				UserEvent: &pb.UserEvent{
					Event: pb.UserEvent_LOGOUT,
					User:  &user.User,
					Time:  time.Now().Unix(),
				},
			},
		},
	)

	c.Debugln("user", creds.Nick, "exited from message listing loop")

	return err
}

func (s *Service) addClientChanToConnected(nick string, ch chan *pb.ChatServerMsg) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.connectedClients[nick] = ch
}

func (s *Service) rmClientChanFromConnected(nick string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.connectedClients, nick)
}
