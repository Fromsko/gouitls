package kratos

import (
	"net/http"

	"github.com/Fromsko/gouitls/logs"
	"github.com/gin-gonic/gin"
	kratosHTTP "github.com/go-kratos/kratos/v2/transport/http"
)

var log = logs.InitLogger()

func init() {
	gin.SetMode(gin.ReleaseMode)
	log.Info("正在启动服务器...")
}

// AuthMiddleware 是一个Gin中间件，用于检查auth token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		// 检查token是否存在
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 2001, "message": "无权访问"})
			c.Abort()
			return
		}
		// 如果token存在，继续处理请求
		c.Next()
	}
}

func NewHTTPServer() *kratosHTTP.Server {
	ginEngine := gin.Default()

	// 添加中间件
	ginEngine.Use(AuthMiddleware())

	// 定义路由
	ginEngine.GET("/demo", func(c *gin.Context) {
		// 处理请求
		c.JSON(http.StatusOK, gin.H{"message": "访问成功"})
	})

	service := kratosHTTP.NewServer(
		kratosHTTP.Address(":8000"),
	)

	service.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("请求来了")
		ginEngine.ServeHTTP(w, r)
	})

	// 适配Gin处理器为Kratos处理器
	// service.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
	// 	ginEngine.ServeHTTP(w, r)
	// })

	return service
}

/*
func main() {
	// 创建HTTP服务器
	httpSrv := NewHTTPServer()

	// 创建gRPC服务器
	grpcSrv := grpc.NewServer(grpc.Address(":9000"))

	// 创建Kratos应用，并添加HTTP和gRPC服务器
	app := kratos.New(
		kratos.Name("myapp"),
		kratos.Server(httpSrv, grpcSrv),
	)

	// 启动应用
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
*/
