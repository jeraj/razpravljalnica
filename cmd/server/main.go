package main

import (
	"log"
	"net"

	pb "github.com/jeraj/razpravljalnica/gen"
	"github.com/jeraj/razpravljalnica/internal/server"
	"google.golang.org/grpc"
)

//TA DATOTOTEKA SAMO ZAZENE SERVER, NIC NI TREBA IMPLEMENTIRATI TUKAJ.
//Nekako kot ON button za server, ostalo se implementira v datoteki internal/server

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterMessageBoardServer(grpcServer, server.NewMessageBoardServer())

	log.Println("Server listening on :50051")
	log.Fatal(grpcServer.Serve(lis))
}
