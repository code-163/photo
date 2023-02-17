package admin

import (
	"github.com/gin-gonic/gin"
	"photo/src/model/common"
	"strings"
	"time"
)

// GetVideoList 获取视频列表
func GetVideoList(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	title := c.Query("title")
	mid := c.Query("mid")
	status := c.Query("status")
	if title != "" {
		where += "POSITION('" + title + "' IN v.title) && "
	}
	if mid != "" {
		where += "v.mid = " + mid + " && "
	}
	if status != "" {
		where += "v.status = " + status + " && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数--------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from video as v " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("SELECT v.*,m.name as m_name FROM `video` as v LEFT JOIN model as m on v.mid = m.id " + wheremap + " order by v.id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// DelMovie 删除视频
func DelMovie(c *gin.Context) (string, int) {
	id := c.Query("id")
	ids := c.Query("ids")
	var idval string
	if id != "" {
		idval = "id=" + id
	} else if ids != "" {
		idval = "id in(" + ids + ")"
	}

	result, _ := common.Db.Exec("delete from movie where " + idval)
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

// AddVideo 添加视频
func AddVideo(c *gin.Context) (int, string) {
	var code int
	var msg string
	title := c.PostForm("title")
	cover := c.PostForm("cover")
	playerUrl := c.PostForm("player_url")
	mid := c.PostForm("mid")
	introduce := c.PostForm("introduce")
	bdwp_url := c.PostForm("bdwp_url")
	carry_code := c.PostForm("carry_code")
	de_passwd := c.PostForm("de_passwd")
	price := c.PostForm("price")
	ctime := time.Now().Format("2006-01-02")
	status := 1
	result, _ := common.Db.Exec("insert into `video` (`title`,`cover`,`player_url`,`mid`,`introduce`,`bdwp_url`,`carry_code`,`de_passwd`,`price`,`ctime`,`status`)values(?,?,?,?,?,?,?,?,?,?,?)", title, cover, playerUrl, mid, introduce, bdwp_url, carry_code, de_passwd, price, ctime, status)
	resCount, _ := result.RowsAffected()
	if resCount > 0 {
		code = 200
		msg = "添加成功"
	} else {
		code = 400
		msg = "添加失败"
	}
	return code, msg
}

// EditVideo 修改视频
func EditVideo(c *gin.Context) (int, string) {
	var code int
	var msg string
	id := c.PostForm("id")
	title := c.PostForm("title")
	cover := c.PostForm("cover")
	playerUrl := c.PostForm("player_url")
	mid := c.PostForm("mid")
	introduce := c.PostForm("introduce")
	bdwp_url := c.PostForm("bdwp_url")
	carry_code := c.PostForm("carry_code")
	de_passwd := c.PostForm("de_passwd")
	price := c.PostForm("price")
	result, _ := common.Db.Exec("update `video` set title=?,cover=?,player_url=?,mid=?,introduce=?,bdwp_url=?,carry_code=?,de_passwd=?,price=? where id=?", title, cover, playerUrl, mid, introduce, bdwp_url, carry_code, de_passwd, price, id)
	resCount, _ := result.RowsAffected()
	if resCount > 0 {
		code = 200
		msg = "修改成功"
	} else {
		code = 400
		msg = "修改失败"
	}
	return code, msg
}

// DelVideo 删除视频
func DelVideo(c *gin.Context) (string, int) {
	id := c.Query("id")
	ids := c.Query("ids")
	var idval string
	if id != "" {
		idval = "id=" + id
	} else if ids != "" {
		idval = "id in(" + ids + ")"
	}

	result, _ := common.Db.Exec("delete from video where " + idval)
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

// StateVideo 视频开关
func StateVideo(c *gin.Context) (string, int) {
	var msg string
	var code int
	id := c.Query("id")
	state := c.Query("state")
	result, _ := common.Db.Exec("update video set status=? where id=?", state, id)
	rowCount, _ := result.RowsAffected()
	if rowCount > 0 {
		msg = "设置成功"
		code = 200
	} else {
		msg = "设置失败"
		code = 400
	}
	return msg, code
}
