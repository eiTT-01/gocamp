package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/gocamp/chapter5/pb"
	"google.golang.org/grpc"
)

type EchoService struct {
	pb.UnimplementedEchoServiceServer
}

func (e EchoService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("SayHello: %v\n", in.GetName())
	return &pb.HelloResponse{Message: "Hello, " + in.GetName()}, nil
}

func (e EchoService) SayHelloClientStream(server pb.EchoService_SayHelloClientStreamServer) error {
	log.Println("SayHelloClientStream")
	for {
		msg, err := server.Recv()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(msg.GetName())
	}
}

func (e EchoService) SayHelloServerStream(in *pb.HelloRequest, server pb.EchoService_SayHelloServerStreamServer) error {
	log.Printf("SayHelloServerStream: %v\n", in.GetName())
	err := server.Send(&pb.HelloResponse{Message: "Hello, " + in.GetName()})
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	return nil
}

func (e EchoService) SayHelloBiStream(server pb.EchoService_SayHelloBiStreamServer) error {
	var wg sync.WaitGroup
	wg.Add(2)
	c := make(chan string, 2)

	// 接收来自客户端的消息
	go func() {
		defer wg.Done()

		for {
			msg, err := server.Recv()
			if err != nil {
				log.Printf("Error: %v\n", err)
				return
			}
			c <- msg.Name
			log.Println(msg.GetName())
			time.Sleep(1 * time.Second)
		}
	}()

	// 发送消息给客户端
	go func() {
		defer wg.Done()

		for {
			name := <-c
			err := server.Send(&pb.HelloResponse{Message: "Hello, " + name})
			if err != nil {
				log.Printf("Error: %v\n", err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()

	return nil
}

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer()
	pb.RegisterEchoServiceServer(srv, &EchoService{})

	err = srv.Serve(listener)
	if err != nil {
		log.Fatalln(err)
	}
}
