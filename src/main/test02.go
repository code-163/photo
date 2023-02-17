package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"photo/src/middleware"
	"time"
)

func writeGo(myChan chan interface{}, title string) {
	fmt.Println("准备写入通道...")
	time.Sleep(time.Second * 3)
	myChan <- title
}

func readGo(myChan chan interface{}) {
	//time.Sleep(time.Second * 3)
	fmt.Println("通道值=>", <-myChan)
}

func main() {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.GET("/index", func(c *gin.Context) {
		var title = c.Query("title")
		myChan := make(chan interface{})
		//go readGo(myChan)
		go writeGo(myChan, title)
		data := <-myChan
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"msg":  "返回数据",
			"data": data,
		})
	})
	_ = r.Run(":9000")
}
