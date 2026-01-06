package server

import (
	"context"
	"sync"
	"log"
	"fmt"
	pb "github.com/jeraj/razpravljalnica/gen"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
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
	
	likes map[int64]map[int64]bool

}

//konstruktor za strežnik
func NewMessageBoardServer() *MessageBoardServer {
	s := &MessageBoardServer{ //inicializira slovar uporabnikov in tem
		users:  make(map[int64]*pb.User),
		topics: make(map[int64]*pb.Topic),
		messages: make(map[int64]*pb.Message),
		likes:    make(map[int64]map[int64]bool),
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

	//pripravimo, da shranimo sporočila
	resp := &pb.GetMessagesResponse{}

	//gremo skozi shranjena sporocila na serverju
	for _, m := range s.messages {
		//pogleda ali sporocilo pripada temi in ce je spocilo novejse od zadnjega, pogleda katerih sporocil client se nima
		if m.TopicId == req.TopicId && m.Id > req.FromMessageId {
			//dodamo v response
			resp.Messages = append(resp.Messages, m)
			//ce smo presegli limit, prekinemo
			if req.Limit > 0 && int32(len(resp.Messages)) >= req.Limit {
				break
			}
		}
	}
	//vrne seznam sporocil
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

func (s *MessageBoardServer) LikeMessage(
	ctx context.Context,
	req *pb.LikeMessageRequest,
) (*pb.Message, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := s.messages[req.MessageId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "message not found")
	}

	if msg.TopicId != req.TopicId {
		return nil, status.Errorf(codes.InvalidArgument, "message not in this topic")
	}

	// inicializiraj mapo, če prvič
	if _, ok := s.likes[msg.Id]; !ok {
		s.likes[msg.Id] = make(map[int64]bool)
	}

	// preveri, ali je uporabnik že lajkal
	if s.likes[msg.Id][req.UserId] {
		return nil, status.Errorf(codes.AlreadyExists, "message already liked by this user")
	}

	//zabeleži like
	s.likes[msg.Id][req.UserId] = true //shranjujejo da vemo kdo je lajkal, da se potem s tem ne
	msg.Likes++

	log.Printf("Message %d liked by user %d (likes=%d)",
		msg.Id, req.UserId, msg.Likes)

	return msg, nil
}


func (s *MessageBoardServer) DeleteMessage(
	ctx context.Context,
	req *pb.DeleteMessageRequest,
) (*emptypb.Empty, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := s.messages[req.MessageId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "message not found")
	}

	if msg.TopicId != req.TopicId {
		return nil, status.Errorf(codes.InvalidArgument, "message not in this topic")
	}

	if msg.UserId != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "cannot delete message")
	}

	// izbriši sporočilo
	delete(s.messages, req.MessageId)

	// izbriši vse like za to sporočilo
	delete(s.likes, req.MessageId)

	log.Printf("Message %d deleted by user %d", req.MessageId, req.UserId)

	return &emptypb.Empty{}, nil
}

func (s *MessageBoardServer) UpdateMessage(
	ctx context.Context,
	req *pb.UpdateMessageRequest,
) (*pb.Message, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	msg, ok := s.messages[req.MessageId]
	if !ok {
		return nil, status.Errorf(codes.NotFound, "message not found")
	}

	if msg.TopicId != req.TopicId {
		return nil, status.Errorf(codes.InvalidArgument, "message not in this topic")
	}

	if msg.UserId != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's message")
	}

	if req.Text == "" {
		return nil, status.Errorf(codes.InvalidArgument, "message text cannot be empty")
	}

	// posodobi besedilo
	msg.Text = req.Text

	log.Printf("Message %d updated by user %d", msg.Id, msg.UserId)

	return msg, nil
}


func (s *MessageBoardServer) Shutdown() { //dodana funkcija, ker sem skos mogla ubijat procese fizicno na racunalniku
	log.Println("Server shutting down")
}


