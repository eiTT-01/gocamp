package main

import (
	"context"
	"log"

	"github.com/gocamp/chapter6/pb/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// RequestIdClientInception RequestId 拦截器
// 客户端拦截器在请求发送到服务端之前被调用
func RequestIdClientInception(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	log.Println("in RequestIdClientInception")
	headers, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		headers = metadata.Pairs()
	}

	value := ctx.Value("x-request-id")
	if requestId, ok := value.(string); ok && requestId != "" {
		headers["x-request-id"] = []string{requestId}
	}

	return invoker(metadata.NewOutgoingContext(ctx, headers), method, req, reply, cc, opts...)
}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(RequestIdClientInception))
	if err != nil {
		log.Fatalln(err)
	}

	client := user.NewUserServiceClient(conn)

	req := &user.AddUserRequest{
		Username: "Channer",
		Password: "123456",
		Email:    "channer.geng@vechain.com",
		Gender:   user.Gender_Male,
		Tags:     map[string]string{"company": "vechain"},
	}

	resp, err := client.AddUser(context.Background(), req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(resp)
}
