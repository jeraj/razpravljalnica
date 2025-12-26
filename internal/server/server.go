package server

import (
	"context"
	"sync"
	"log"
	pb "github.com/jeraj/razpravljalnica/gen"
)

type MessageBoardServer struct {
	pb.UnimplementedMessageBoardServer

	mu         sync.Mutex
	users      map[int64]*pb.User
	nextUserID int64
}

func NewMessageBoardServer() *MessageBoardServer {
	return &MessageBoardServer{
		users: make(map[int64]*pb.User),
	}
}

func (s *MessageBoardServer) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.User, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextUserID++
	user := &pb.User{
		Id:   s.nextUserID,
		Name: req.Name,
	}

	s.users[user.Id] = user
	log.Printf("New user created: id=%d name=%s\n", user.Id, user.Name)
	return user, nil
}
