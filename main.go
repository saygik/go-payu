package main

import (
	"github.com/crgimenes/goconfig"
	_ "github.com/crgimenes/goconfig/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/nuveo/log"
	//"os/user"
)

type mailConfig struct {
	User     string `cfgDefault:"crystalrentalcarsender@gmail.com" cfgRequired:"true"`
	Password string `cfgDefault:"qwe123QWE@" cfgRequired:"true"`
	Server   string `cfgDefault:"smtp.gmail.com" cfgRequired:"true"`
	Port     int    `cfgDefault:"587" cfgDefault:"587"`
	MailTo   string `cfgDefault:"def@def.by" cfgRequired:"true"`
}
type payUConfig struct {
	ClientID      string `cfgDefault:"369485" cfgRequired:"true"`
	Secret        string `cfgDefault:"0c21bfc7a5fd8673b637f05385557004" cfgRequired:"true"`
	PayUBase      string `cfgDefault:"https://secure.snd.payu.com" cfgRequired:"true"`
	MerchantPosId string `cfgDefault:"369485" cfgRequired:"true"`
}

type ConfigApp struct {
	AppPort  string `cfgDefault:":80" cfgRequired:"true"`
	Mail     mailConfig
	PayU     payUConfig
	IgnoreMe string `cfg:"-"`
}

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

var AppConfig = ConfigApp{}

func main() {
	//	var config *configApp
	//	goconfig.File = "config/config.json"
	err := goconfig.Parse(&AppConfig)

	if err != nil {
		//	log.Fatal("Config Error")
		log.Fatal(err)
		return
	}
	// ***************** Print App config ******************
	//j, _ := json.MarshalIndent(Config, "", "\t")
	//log.Println(string(j))

	gin.SetMode(gin.DebugMode)
	//gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.RedirectTrailingSlash = false
	//	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(CORSMiddleware())
	r.Use(sessions.Sessions("gin-boilerplate-session", store))
	//	r.LoadHTMLGlob("./public/html/*")
	//gopath := os.Getenv("GOPATH")
	//if gopath == "" {
	//	gopath = "/root/go"
	//}
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"Code": "200", "Message": "Welcome to API go-payu"})
	})
	r.POST("/notify", setNotify)

	r.GET("/auth", getAuth)

	r.GET("/mail", sendMail)
	r.POST("/orders", createOrder)

	r.NoRoute(NoRoute)

	r.Run(AppConfig.AppPort)
}
func NoRoute(c *gin.Context) {
	c.JSON(404, gin.H{"Code": "404", "Message": "Not Found"})
	c.Abort()

}
