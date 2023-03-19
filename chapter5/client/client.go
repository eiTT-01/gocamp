package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/gocamp/chapter5/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewEchoServiceClient(conn)

	// 1. 普通模式
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "Channer"})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp.GetMessage())

	//2. 服务端流模式
	//rs, err := client.SayHelloServerStream(context.Background(), &pb.HelloRequest{Name: "Channer"})
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//for {
	//	resp, err = rs.Recv()
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	log.Println(resp.GetMessage())
	//}

	// 3. 客户端流模式
	//cs, err := client.SayHelloClientStream(context.Background())
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//names := []string{"Channer#1", "Channer#2", "Channer#3"}
	//
	//for _, name := range names {
	//	err = cs.Send(&pb.HelloRequest{Name: name})
	//	if err != nil {
	//		log.Printf("Fail to send %s[err=%v], skip\n", name, err)
	//		continue
	//	}
	//	log.Printf("%s has been sent", name)
	//	time.Sleep(1 * time.Second)
	//}

	// 4. 双向数据流模式
	cs, err := client.SayHelloBiStream(context.Background())
	if err != nil {
		log.Fatalln(err)
	}

	names := []string{"Channer#1", "Channer#2", "Channer#3", "Channer#4", "Channer#5"}

	var wg sync.WaitGroup
	wg.Add(2)

	// 发送消息到服务端
	go func(n []string) {
		defer wg.Done()

		for _, name := range n {
			err = cs.Send(&pb.HelloRequest{Name: name})
			if err != nil {
				log.Printf("Fail to send %s[err=%v], skip\n", name, err)
				continue
			}
			log.Printf("%s has been sent", name)
			time.Sleep(1 * time.Second)
		}

		err = cs.CloseSend()
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Everything is sent successfully")
	}(names)

	// 接收服务端的消息
	go func() {
		defer wg.Done()

		for {
			msg, err := cs.Recv()
			if err != nil {
				log.Printf("Fail to recv msg from server: %v", err)
				return
			}
			log.Printf("Got message: %s\n", msg.GetMessage())
		}
	}()

	wg.Wait()
}
