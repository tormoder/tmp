package storage

import pb "github.com/tormoder/chat/proto"

type User struct {
	Online     bool
	MsgChannel chan *pb.ChatServerMsg
	pb.User
}

type ByNick []*pb.User

func (s ByNick) Len() int           { return len(s) }
func (s ByNick) Less(i, j int) bool { return s[i].Nick < s[j].Nick }
func (s ByNick) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
