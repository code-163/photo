package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"net/http"
	"photo/src/middleware"
)

func main() {
	r := gin.Default()
	r.Use(middleware.Cors())
	//store := cookie.NewStore([]byte("secret"))
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("mySession", store))
	r.GET("/index", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Options(sessions.Options{MaxAge: 60})
		/*if session.Get("username") == nil {
			session.Set("username", "alicefelitsa")
			_ = session.Save()
		}*/
		fmt.Println("session：", session.Get("username"))
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  session.Get("username"),
		})
	})
	_ = r.Run(":9000")
}
