package admin

import (
	"github.com/gin-gonic/gin"
	"photo/src/model/common"
	"strings"
	"time"
)

// GetModelList 获取模特列表
func GetModelList(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	name := c.Query("name")
	if name != "" {
		where += "POSITION('" + name + "' IN name) && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数--------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from model " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("select * from model " + wheremap + " order by id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// GetModel 获取模特数据
func GetModel(c *gin.Context) []map[string]string {
	rows, _ := common.Db.Query("select id as value,name as label from model order by id desc")
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result
}

// AddModel 添加模特
func AddModel(c *gin.Context) (string, int) {
	name := c.PostForm("name")
	introduce := c.PostForm("introduce")
	image := c.PostForm("image")
	ctime := time.Now().Format("2006-01-02 15:04:05")
	result, _ := common.Db.Exec("insert into model(`name`,`introduce`,`image`,`ctime`)values(?,?,?,?)", name, introduce, image, ctime)
	rowCount, _ := result.RowsAffected()
	var msg string
	var code int
	if rowCount > 0 {
		msg = "添加成功"
		code = 200
	} else {
		msg = "添加失败"
		code = 400
	}
	return msg, code
}

// EditModel 修改模特
func EditModel(c *gin.Context) (string, int) {
	id := c.PostForm("id")
	name := c.PostForm("name")
	introduce := c.PostForm("introduce")
	image := c.PostForm("image")
	result, _ := common.Db.Exec("update model set name=?,introduce=?,image=? where id=?", name, introduce, image, id)
	rowCount, _ := result.RowsAffected()
	var msg string
	var code int
	if rowCount > 0 {
		msg = "修改成功"
		code = 200
	} else {
		msg = "修改失败"
		code = 400
	}
	return msg, code
}

// DelModel 删除模特
func DelModel(c *gin.Context) (string, int) {
	id := c.Query("id")
	ids := c.Query("ids")
	var idval string
	if id != "" {
		idval = "id=" + id
	} else if ids != "" {
		idval = "id in(" + ids + ")"
	}

	result, _ := common.Db.Exec("delete from model where " + idval)
	rowCount, _ := result.RowsAffected()
	var msg string
	var code int
	if rowCount > 0 {
		msg = "删除成功"
		code = 200
	} else {
		msg = "删除失败"
		code = 400
	}
	return msg, code
}
