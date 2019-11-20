package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"os"
	//"os/user"
)

//CORSMiddleware ...
func CORSMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Cache-Control", "private, no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Expires", "-1")
		c.Writer.Header().Set("Pragma", "no-cache")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	gin.SetMode(gin.DebugMode)
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.RedirectTrailingSlash = false
	//	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(CORSMiddleware())
	r.Use(sessions.Sessions("gin-boilerplate-session", store))
	//	r.LoadHTMLGlob("./public/html/*")
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = "/root/go"
	}
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"Code": "200", "Message": "Welcome to API go-payu"})
	})
	r.POST("/notify", setNotify)

	r.GET("/auth", getAuth)
	r.POST("/orders", createOrder)

	r.NoRoute(NoRoute)

	r.Run(":80")
}
func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"Code": "404", "Message": "Not Found"})
	c.Abort()

}
