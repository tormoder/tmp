package user

import (
	"sort"
	"time"

	"github.com/tormoder/chat/chat"
	c "github.com/tormoder/chat/common"
	pb "github.com/tormoder/chat/proto"
	"github.com/tormoder/chat/storage"

	"golang.org/x/net/context"
)

type Service struct {
	storage storage.UserStorage
	chat    *chat.Service
}

func NewService(chatService *chat.Service, userStorage storage.UserStorage) *Service {
	return &Service{
		storage: userStorage,
		chat:    chatService,
	}
}

func (s *Service) Login(ctx context.Context, lreq *pb.LoginRequest) (*pb.Credentials, error) {
	c.Debugln("login request from", lreq.Nick)
	user, found := s.storage.GetUser(lreq.Nick)
	if found {
		if user.Online {
			return nil, c.AuthenticationError("user already online")
		}
		user.Online = true
		user.TimeLastSeen = time.Now().Unix()
		err := s.storage.UpdateUser(user)
		if err != nil {
			return nil, err
		}
	} else {
		user = storage.User{
			Online:     true,
			MsgChannel: make(chan *pb.ChatServerMsg, 2048),
			User: pb.User{
				Nick:         lreq.Nick,
				TimeLastSeen: time.Now().Unix(),
			},
		}
		err := s.storage.AddUser(user)
		if err != nil {
			return nil, err
		}
	}

	s.chat.BroadcastAllConnectedClients(
		&pb.ChatServerMsg{
			Msg: &pb.ChatServerMsg_UserEvent{
				UserEvent: &pb.UserEvent{
					Event: pb.UserEvent_LOGIN,
					User:  &user.User,
					Time:  time.Now().Unix(),
				},
			},
		},
	)

	c.Debugln("user", user.User.Nick, "logged-in")

	return &pb.Credentials{
		Nick: user.User.Nick,
	}, nil
}

func (s *Service) Logout(ctx context.Context, creds *pb.Credentials) (*pb.LogoutResponse, error) {
	c.Debugln("logout request from", creds.Nick)
	user, err := s.storage.CheckCredentials(creds)
	if err != nil {
		return nil, err
	}

	user.Online = false
	user.TimeLastSeen = time.Now().Unix()
	err = s.storage.UpdateUser(user)
	if err != nil {
		return nil, c.InternalServerError("storage error")
	}

	s.chat.BroadcastAllConnectedClients(
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

	c.Debugln("user", user.User.Nick, "logged-out")

	return &pb.LogoutResponse{}, nil
}

func (s *Service) ListUsers(ctx context.Context, creds *pb.Credentials) (*pb.ListUsersResponse, error) {
	c.Debugln("list user request from", creds.Nick)
	_, err := s.storage.CheckCredentials(creds)
	if err != nil {
		return nil, err
	}
	users := s.storage.GetAllOnlineUsersDTO()
	sort.Sort(storage.ByNick(users))
	return &pb.ListUsersResponse{
		Users: users,
	}, nil
}
