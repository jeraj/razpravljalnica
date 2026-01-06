package server

import (
	"context"
	"sync"
	"log"
	"fmt"
	pb "github.com/jeraj/razpravljalnica/gen"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type MessageBoardServer struct {
	pb.UnimplementedMessageBoardServer

	mu sync.Mutex

	users      map[int64]*pb.User
	nextUserID int64

	topics      map[int64]*pb.Topic
	nextTopicID int64

	messages      map[int64]*pb.Message
	nextMessageID int64
}

//konstruktor za strežnik
func NewMessageBoardServer() *MessageBoardServer {
	s := &MessageBoardServer{ //inicializira slovar uporabnikov in tem
		users:  make(map[int64]*pb.User),
		topics: make(map[int64]*pb.Topic),
		messages: make(map[int64]*pb.Message),
	}

	//naredim default temo
	s.nextTopicID = 1
	s.topics[1] = &pb.Topic{
		Id:   1,
		Name: "General",
	}

	log.Println("Default topic created: id=1 name=General")

	return s
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

func (s *MessageBoardServer) ListTopics(
	ctx context.Context,
	_ *emptypb.Empty,
) (*pb.ListTopicsResponse, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	resp := &pb.ListTopicsResponse{}
	for _, t := range s.topics {
		resp.Topics = append(resp.Topics, t)
	}
	return resp, nil
}

func (s *MessageBoardServer) CreateTopic(
	ctx context.Context,
	req *pb.CreateTopicRequest,
) (*pb.Topic, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextTopicID++

	topic := &pb.Topic{
		Id:   s.nextTopicID,
		Name: req.Name,
	}

	s.topics[topic.Id] = topic

	log.Printf("New topic created: id=%d name=%s\n", topic.Id, topic.Name)

	return topic, nil
}

func (s *MessageBoardServer) GetMessages(
	ctx context.Context,
	req *pb.GetMessagesRequest,
) (*pb.GetMessagesResponse, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	resp := &pb.GetMessagesResponse{}

	for _, m := range s.messages {
		if m.TopicId == req.TopicId && m.Id > req.FromMessageId {
			resp.Messages = append(resp.Messages, m)

			if req.Limit > 0 && int32(len(resp.Messages)) >= req.Limit {
				break
			}
		}
	}

	return resp, nil
}

func (s *MessageBoardServer) PostMessage(
	ctx context.Context,
	req *pb.PostMessageRequest,
) (*pb.Message, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	// preveri, če obstaja uporabnik
	user, ok := s.users[req.UserId]
	if !ok {
		return nil, fmt.Errorf("user with id %d does not exist", req.UserId)
	}

	// preveri, če obstaja tema
	topic, ok := s.topics[req.TopicId]
	if !ok {
		return nil, fmt.Errorf("topic with id %d does not exist", req.TopicId)
	}

	// ustvari novo sporočilo
	s.nextMessageID++
	msg := &pb.Message{
		Id:        s.nextMessageID,
		TopicId:   topic.Id,
		UserId:    user.Id,
		Text:      req.Text,
		Likes:     0,
		CreatedAt: timestamppb.Now(), // trenutni timestamp
	}

	// shrani sporočilo
	s.messages[msg.Id] = msg

	log.Printf("New message posted: id=%d topic=%d user=%d text=%s\n", msg.Id, topic.Id, user.Id, msg.Text)

	return msg, nil
}
