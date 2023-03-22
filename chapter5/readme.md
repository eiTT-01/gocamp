# 介绍一下RPC

RPC 全称远程过程调用（Remote Procedure Call）。一种通过网络从远程计算机上请求服务，而不需要了解底层网络技术的通信协议。

分布式系统的关键部分，RPC理想上想把网络通信实现同本地函数调用一样，RPC的目的：

- 更容易编写网络通信程序
- 隐藏客户端服务器通信的细节
- 客户端调用更加像本地的过程调用
- 服务端处理更加像本地的过程调用

# 写一个Post的http的客户端调用

```
type AddUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AddUserResponse struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func main() {
	payload := AddUserRequest{
		Username: "channer",
		Password: "111111",
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post("http://127.0.0.1:8080/users", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}

	var res AddUserResponse

	_ = json.NewDecoder(resp.Body).Decode(&res)

	fmt.Println(res)
}
```

# 描述一下RPC四要素

- RPC Client: RPC 协议的调用方 
- RCP Proxy/Stub: RPC 代理存在于客户端，以便于调用者可以「透明」的调用服务端函数
- RPC Server: RPC 协议的服务端，用于实现远程服务的各类方法
- RPC Selector/Processor: 存在于服务端，是一个负责执行 RPC 接口实现的角色，也有诸如管理接口注册、客户端权限检查等工作。

# 如何实现一个GO语言原生的，跨语言的，并发高的RPC示例

```
// server
type EchoService struct{}

func (s *EchoService) Echo(request string, response *string) error {
	*response = "Hello, " + request
	return nil
}

func main() {
	// 建立一个 TCP 监听
	listen, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalln(err)
		return
	}

	// 注册 RCP 服务
	err = rpc.RegisterName("EchoService", &EchoService{})
	if err != nil {
		log.Fatalln(err)
		return
	}

	// 接收&处理客户端请求
	for {
		// 接收请求，并建立 TCP 连接
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			return
		}
		// 反序列化并处理请求
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}

// client
func main() {
	// 创建 TCP 连接
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatalln(err)
		return
	}

	// 创建 JSON-RPC 客户端
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	// 发起请求
	var reply string
	err = client.Call("EchoService.Echo", "Channer", &reply)
	if err != nil {
		log.Fatalln(err)
		return
	}

	fmt.Println(reply)
}
```

# protobuf和json的优势是什么

- 多语言、跨平台
- 传输快、压缩率高
- 易维护


# 写一个grpc+probufer，客户端调用服务端的案例

```
// server
type EchoService struct {
	pb.UnimplementedEchoServiceServer
}

func (s *EchoService) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Message: "Hello, " + request.GetName()}, nil
}

func main() {
	srv := grpc.NewServer()
	pb.RegisterEchoServiceServer(srv, &EchoService{})
	listen, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalln(err)
	}
	err = srv.Serve(listen)
	if err != nil {
		log.Fatalln(err)
	}
}

// client
func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}

	client := pb.NewEchoServiceClient(conn)

	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "I'm Channer"})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp.GetMessage())
}
```

# 介绍一下GRPC的3种流模式，举例说明

- 服务端数据流模式（Server-side Streaming RPC）: 这种模式是客户端发起一次请求，服务端返回一段连续的数据流。典型的例子是客户端向服务端发送一个股票代码，服务端就把该股票的实时数据源源不断的返回给客户端。典型的例子是客户端向服务端发送一个股票代码，服务端就把该股票的实时数据源源不断的返回给客户端。
- 客户端数据流模式（Client-side Streaming RPC）: 与服务端数据流模式相反，这次是客户端源源不断的向服务端发送数据流，而在发送结束后，由服务端返回一个响应。典型的例子是物联网终端向服务器报送数据。
- 双向数据流模式（Bidirectional Streaming RPC）: 这是客户端和服务端都可以向对方发送数据流，这个时候双方的数据可以同时互相发送，也就是可以实现实时交互。典型的例子是聊天机器人。

```
// echo.proto
service EchoService {
  // 普通 RPC 模式
  rpc SayHello(HelloRequest) returns (HelloResponse) {}
  // 客户端数据流模式
  rpc SayHelloClientStream(stream HelloRequest) returns (HelloResponse) {}
  // 服务端数据流模式
  rpc SayHelloServerStream(HelloRequest) returns (stream HelloResponse) {}
  // 双向数据流模式
  rpc SayHelloBiStream(stream HelloRequest) returns (stream HelloResponse) {}
}

// 服务端实现
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

// 客户端实现
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
rs, err := client.SayHelloServerStream(context.Background(), &pb.HelloRequest{Name: "Channer"})
if err != nil {
    log.Fatalln(err)
}
for {
    resp, err = rs.Recv()
    if err != nil {
        log.Fatalln(err)
    }
    log.Println(resp.GetMessage())
}

// 3. 客户端流模式
cs, err := client.SayHelloClientStream(context.Background())
if err != nil {
    log.Fatalln(err)
}

names := []string{"Channer#1", "Channer#2", "Channer#3"}

for _, name := range names {
    err = cs.Send(&pb.HelloRequest{Name: name})
    if err != nil {
        log.Printf("Fail to send %s[err=%v], skip\n", name, err)
        continue
    }
    log.Printf("%s has been sent", name)
    time.Sleep(1 * time.Second)
}

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
```
