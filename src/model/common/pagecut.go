package common

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// Pagecut 处理数据分页
func Pagecut(c *gin.Context) string {
	pageSize := c.Query("limit")
	p, _ := strconv.Atoi(c.Query("page"))
	l, _ := strconv.Atoi(pageSize)
	page := strconv.Itoa((p - 1) * l)
	limit := "limit " + page + "," + pageSize
	return limit
}
