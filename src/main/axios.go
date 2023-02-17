package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"photo/src/middleware"
)

func main() {
	r := gin.Default()
	r.Use(middleware.Cors())
	r.NoRoute(NotPages)
	h5 := r.Group("h5")
	{
		h5.GET("/axiosGet", axiosGet)
		h5.POST("/axiosPost", axiosPost)
	}
	_ = r.Run(":8000")
}

// NotPages 当访问一个错误网站时返回
func NotPages(c *gin.Context) {
	c.JSON(http.StatusNotFound, gin.H{
		"status": 404,
		"error":  "找不到对应路径或没有此文件以及页面",
	})
}

func axiosGet(c *gin.Context) {
	data := make(map[string]interface{})
	data["name"] = c.Query("name")
	data["age"] = 33
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "返回数据",
		"data": data,
	})
}

func axiosPost(c *gin.Context) {
	data := make(map[string]interface{})
	data["name"] = c.PostForm("name")
	data["age"] = 30
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "返回数据",
		"data": data,
	})
}
