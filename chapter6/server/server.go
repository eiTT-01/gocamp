package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/gocamp/chapter6/pb/user"
	uuid2 "github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserController struct {
	user.UnimplementedUserServiceServer
}

func (u UserController) AddUser(ctx context.Context, req *user.AddUserRequest) (*user.AddUserResponse, error) {
	headers, _ := metadata.FromIncomingContext(ctx)
	for k, v := range headers {
		log.Printf("%s = %s\n", k, v)
	}

	data := &user.AddUserResponseData{
		Username:  req.GetUsername(),
		Password:  req.GetPassword(),
		Email:     req.GetPassword(),
		CreatedAt: time.Now().Unix(),
		UpdatedAt: timestamppb.Now(),
		Tags:      req.Tags,
	}
	return &user.AddUserResponse{Code: 1, Message: "Hello, " + req.GetUsername(), Data: data}, nil
}

// RequestIdServerInception RequestId 服务端拦截器
// 该拦截器将在执行处理函数之前被执行
func RequestIdServerInception(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.Println("in RequestIdInception")

	// 1. 获取上下文的 metadata
	headers, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		headers = metadata.Pairs()
	}

	// 2. 如果来源已经加了请求号，则使用客户端提供的请求号
	ids := headers["x-request-id"]
	if len(ids) >= 1 {
		return handler(ctx, req)
	}

	// 3. 否则，生成一个新的请求号
	uuid, err := uuid2.NewUUID()
	if err != nil {
		log.Println(err)
		return handler(ctx, req)
	}

	headers["x-request-id"] = []string{uuid.String()}

	return handler(metadata.NewIncomingContext(ctx, headers), req)
}

func main() {
	port := 8080
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		log.Fatalln(err)
	}

	srv := grpc.NewServer(grpc.UnaryInterceptor(RequestIdServerInception))

	user.RegisterUserServiceServer(srv, &UserController{})

	log.Printf("Server is listening on %d\n", port)

	err = srv.Serve(lis)
	if err != nil {
		log.Fatalln(err)
	}
}
