# 介绍一下Gin

Gin 是一个 Go 编写的轻量级 Web 开发框架。

# Defalut和New有什么不同

`gin.Default()` 默认使用了 `Logger` 和 `Recovery` 中间件，其中：
- `Logger` 中间件将日志写入 `gin.DefaultWriter`, 即使配置了 `GIN_MODE=release`
- `Recovery` 中间件会 `recover` 任何 `panic`，如果有 `panic` 的话，为写入 500 响应码

如果不想使用上面的默认中间件，可以使用 `gin.New()` 新建一个没有任何中间件的路由。

# 针对表单做 restful 的处理

```
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	
	r.POST("/users", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.DefaultPostForm("password", "123456")
		
		c.JSON(http.StatusOK, gin.H{
            "username": username,
            "password": password,
        })
	})

	r.Run(":8080")
}

```


# 针对 JSON 做 restful 的处理

```
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func main() {
	r := gin.Default()

	r.POST("/users", func(c *gin.Context) {
		var req AddUserRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "bad request",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"username": req.Username,
				"password": req.Password,
			})
		}
	})

	_ = r.Run("0.0.0.0:8080")
}

```

# 举例如何用到中间件

```
r := gin.New()

// 作用于全局
r.Use(gin.Logger())
r.Use(gin.Recovery())

// 作用于单个路由
r.POST("/users", gin.Logger(), func(c *gin.Context) {
    // ...
})

_ = r.Run("0.0.0.0:8080")

```

# 如果多个中间件，调用顺序如何？（代码演示）

```

func MiddlewareOne() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("before MiddlewareOne: " + time.Now().String())
		c.Next()
		fmt.Println("after MiddlewareOne: " + time.Now().String())
	}
}

func MiddlewareTwo() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("before MiddlewareTwo: " + time.Now().String())
		c.Next()
		fmt.Println("after MiddlewareTwo: " + time.Now().String())
	}
}

func MiddlewareThree() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("before MiddlewareThree: " + time.Now().String())
		c.Next()
		fmt.Println("after MiddlewareThree: " + time.Now().String())
	}
}

r := gin.New()

// 作用于全局
r.Use(gin.Logger())

r.Use(MiddlewareOne())
r.Use(MiddlewareTwo())
r.Use(MiddlewareThree())

r.Use(gin.Recovery())

// terminal output
[GIN-debug] Listening and serving HTTP on 0.0.0.0:8080
before MiddlewareOne: 2023-03-15 15:17:51.415397 +0800 CST m=+15.635784001
before MiddlewareTwo: 2023-03-15 15:17:51.415991 +0800 CST m=+15.636377251
before MiddlewareThree: 2023-03-15 15:17:51.416 +0800 CST m=+15.636386793
[GIN] 2023/03/15 - 15:17:51 | 200 |      667.25µs |       127.0.0.1 | POST     /users
after MiddlewareThree: 2023-03-15 15:17:51.416711 +0800 CST m=+15.637097876
after MiddlewareTwo: 2023-03-15 15:17:51.416716 +0800 CST m=+15.637103084
after MiddlewareOne: 2023-03-15 15:17:51.41672 +0800 CST m=+15.637106626
[GIN] 2023/03/15 - 15:17:51 | 200 |     1.32675ms |       127.0.0.1 | POST     /users

// conclusion
1. 对于 next() 之前的代码，是按 use 的顺序，先进先出地执行
2. 对于 next() 之后的代码，是按 use 的顺序，后进先出地执行
```

# Gin终止其中一个中间件，要如何做？

通过 `c.Abort()` 可以控制是否要终止后续 pending 的 handlers

# 如何优雅退出Gin的程序

所谓优雅退出指的是在程序完全退出前，需完成，
1. 关闭所有的监听端口
2. 关闭所有的空闲连接
3. 等待活动的连接处理完毕转为空闲并关闭

```
func gracefulShutdown(srv *http.Server) {
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	<-quit // 阻塞

	// 5s 后强制退出
	ctx, channel := context.WithTimeout(context.Background(), 5*time.Second)

	defer channel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown error")
	}
	log.Println("server exiting...")
}

func main() {
	r := gin.New()

	// 作用于全局
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 作用于单个路由
	r.POST("/users", gin.Logger(), func(c *gin.Context) {
		var req AddUserRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": "bad request",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"username": req.Username,
				"password": req.Password,
			})
		}
	})

	srv := &http.Server{
		Addr:      ":8080",
		Handler:   r,
	}

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	gracefulShutdown(srv)
}
```

# 利用之前的Go语言学到的知识，可以总结出一次请求处理的大体流程？

1. `gin.New()` 创建一个 gin 实例 
2. `ListenAndServe` 启动一个监听端口，等待客户端请求进来
3. 接收请求并构建 *conn，启动 goroutine 执行 `serve()` 方法开始处理请求
4. 匹配路由器，执行 `ServeHTTP()` 方法
5. `ServeHTTP()` 方法将请求分发给与请求 url 匹配最近的 handler
6. 返回处理结果给客户端，完成一次请求


# gin返回html的处理（选做）

```
//<!DOCTYPE html>
//<html lang="en">
//<head>
//    <meta charset="UTF-8">
//    <meta http-equiv="X-UA-Compatible" content="IE=edge">
//    <meta name="viewport" content="width=device-width, initial-scale=1.0">
//    <title>{{ .title }}</title>
//</head>
//<body>
//</body>
//</html>
r.LoadHTMLGlob("./templates/index.tmpl")

r.GET("/html", func(c *gin.Context) {
    c.HTML(http.StatusOK, "index.tmpl", gin.H{
        "title": "Hello Golang",
    })
})

```

# gin如何处理静态文件（选做）

```
r := gin.Default()
r.Static("/assets", "./assets")
```
