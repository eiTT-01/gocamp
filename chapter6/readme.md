# protobuf数据类型有哪些，和go语言如何对应的

| ProtoBuf | Go      |
|----------|---------|
| string   | string  |
| int32    | int32   |
| int64    | int64   |
| uint32   | uint32  |
| uint64   | uint64  |
| float    | float32 |
| double   | float64 |
| bool     | bool    |
| bytes    | []byte  |


# 说说go_package

`go_package` 用来定义当前 proto 文件生成的代码路径和包名。比如，

`option go_package = "./pb";`

表示生成的 Go 文件目录是当前目录的 pb 文件夹，包名为 pb。默认情况下，protoc 使用 `go_package` 目录路径最后那部分作为包名，如果想单独指定跟目录不一样的包名，则可以通过分号，如，

`option go_package = "./pb;echo";`

上述选项，将在 `./pb` 目录，生成包名为 `echo` 的 Go 文件。

# 举例说明如何使用 protobuf 的 map, timestamp 和枚举

```
// user.proto
message AddUserRequest {
  string Username = 1;
  string Password = 2;
  string Email = 3;
  Gender gender = 4;
  map<string, string> Tags = 5;
}

message AddUserResponse {
  int32 Code = 1;
  string Message = 2;
  AddUserResponseData Data = 3;
}

message AddUserResponseData {
  string Username = 1;
  string Password = 2;
  string Email = 3;
  Gender gender = 4;
  int64 CreatedAt = 5;
  google.protobuf.Timestamp UpdatedAt = 6;
  map<string, string> Tags = 7;
}

// client
req := &user.AddUserRequest{
    Username: "Channer",
    Password: "123456",
    Email:    "channer.geng@vechain.com",
    Gender:   user.Gender_Male,
    Tags:     map[string]string{"company": "vechain"},
}

// server
data := &user.AddUserResponseData{
    Username:  req.GetUsername(),
    Password:  req.GetPassword(),
    Email:     req.GetPassword(),
    CreatedAt: time.Now().Unix(),
    UpdatedAt: timestamppb.Now(),
    Tags:      req.Tags,
}

```

# grpc的拦截器是如何使用的？

grpc 总共有 4 种拦截器，2 个 UnaryInterceptor 和 2 个 StreamInterceptor，每个都包含一个服务端拦截器和一个客户端拦截器。

## Client Unary Interceptor

客户端拦截器在客户端被执行，可以用于请求的预处理以及响应的后处理。

![client interceptor](https://techdozo.dev/wp-content/uploads/2022/04/image-3.png)

## Server Unary Interceptor

服务端拦截器将在请求到达服务端时被执行，如下图所示，你可以通过服务端拦截器拦截处理发送到服务端的请求，以及发送给客户端的响应。

![server interceptor](https://techdozo.dev/wp-content/uploads/2022/04/image-4.png)

## Client Stream Interceptor

该拦截器的作用跟 `Client Unary Interceptor` 类似。它是在客户端实现的，用于拦截客户端流式发送至服务端的请求。

## Server Stream Interceptor

在 gRPC 服务端接收到流式请求时被执行。

# 举例说明message的嵌套和import的用法

```
import "google/protobuf/timestamp.proto";

message AddUserResponse {
  int32 Code = 1;
  string Message = 2;
  AddUserResponseData Data = 3;
}

message AddUserResponseData {
  string Username = 1;
  string Password = 2;
  string Email = 3;
  Gender gender = 4;
  int64 CreatedAt = 5;
  google.protobuf.Timestamp UpdatedAt = 6;
  map<string, string> Tags = 7;
}
```