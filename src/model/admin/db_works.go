package admin

import (
	"github.com/gin-gonic/gin"
	"photo/src/model/common"
	"strings"
	"time"
)

// GetWorks 获取作品列表
func GetWorks(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	title := c.Query("title")
	mid := c.Query("mid")
	status := c.Query("status")
	types := c.Query("types")
	if title != "" {
		where += "POSITION('" + title + "' IN w.title) && "
	}
	if mid != "" {
		where += "w.mid = " + mid + " && "
	}
	if status != "" {
		where += "w.status = " + status + " && "
	}
	if types != "" {
		where += "w.types = " + types + " && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数--------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from works as w " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("SELECT w.*,m.name as m_name FROM `works` as w LEFT JOIN model as m on w.mid = m.id " + wheremap + " order by w.id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// AddWorks 添加作品
func AddWorks(c *gin.Context) (int, string) {
	var code int
	var msg string
	title := c.PostForm("title")
	cover := c.PostForm("cover")
	atlas := c.PostForm("atlas")
	mid := c.PostForm("mid")
	introduce := c.PostForm("introduce")
	bdwp_url := c.PostForm("bdwp_url")
	carry_code := c.PostForm("carry_code")
	de_passwd := c.PostForm("de_passwd")
	price := c.PostForm("price")
	types := c.PostForm("types")
	ctime := time.Now().Format("2006-01-02")
	status := 1
	result, _ := common.Db.Exec("insert into `works` (`title`,`cover`,`atlas`,`mid`,`introduce`,`bdwp_url`,`carry_code`,`de_passwd`,`price`,`ctime`,`status`,`types`)values(?,?,?,?,?,?,?,?,?,?,?,?)", title, cover, atlas, mid, introduce, bdwp_url, carry_code, de_passwd, price, ctime, status, types)
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

// EditWorks 修改作品
func EditWorks(c *gin.Context) (int, string) {
	var code int
	var msg string
	id := c.PostForm("id")
	title := c.PostForm("title")
	cover := c.PostForm("cover")
	atlas := c.PostForm("atlas")
	mid := c.PostForm("mid")
	introduce := c.PostForm("introduce")
	bdwp_url := c.PostForm("bdwp_url")
	carry_code := c.PostForm("carry_code")
	de_passwd := c.PostForm("de_passwd")
	price := c.PostForm("price")
	types := c.PostForm("types")
	result, _ := common.Db.Exec("update `works` set title=?,cover=?,atlas=?,mid=?,introduce=?,bdwp_url=?,carry_code=?,de_passwd=?,price=?,types=? where id=?", title, cover, atlas, mid, introduce, bdwp_url, carry_code, de_passwd, price, types, id)
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

// EditWorkTypes 批量修改作品类型
func EditWorkTypes(c *gin.Context) (int, string) {
	var code int
	var msg string
	var idVal string
	ids := c.PostForm("ids")
	types := c.PostForm("types")
	if ids != "" {
		idVal = "id in(" + ids + ")"
	}
	result, _ := common.Db.Exec("update `works` set types=? where "+idVal, types)
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

// GetWorksData 获取单个作品数据
func GetWorksData(c *gin.Context) (int, string, map[string]string) {
	var code int
	var msg string
	id := c.Query("id")
	rows, _ := common.Db.Query("select * from works where id=?", id)
	defer rows.Close()
	result := common.AssemblyData(rows)
	if len(result) > 0 {
		code = 200
		msg = "获取数据成功"
	} else {
		code = 400
		msg = "获取数据失败"
	}
	return code, msg, result[0]
}

// UpdateWorksImage 更新作品图片
func UpdateWorksImage(c *gin.Context) (int, string) {
	var code int
	var msg string
	id := c.PostForm("id")
	cover := c.PostForm("cover") //主图
	atlas := c.PostForm("atlas") //图册
	result, _ := common.Db.Exec("update `works` set cover=?,atlas=? where id=?", cover, atlas, id)
	resCount, _ := result.RowsAffected()
	if resCount > 0 {
		code = 200
		msg = "修改图片成功"
	} else {
		code = 400
		msg = "修改图片失败"
	}
	return code, msg
}

// DelWorks 删除作品
func DelWorks(c *gin.Context) (string, int) {
	id := c.Query("id")
	ids := c.Query("ids")
	var idval string
	if id != "" {
		idval = "id=" + id
	} else if ids != "" {
		idval = "id in(" + ids + ")"
	}
	result, _ := common.Db.Exec("delete from works where " + idval)
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

// StateWorks 作品开关
func StateWorks(c *gin.Context) (string, int) {
	var msg string
	var code int
	id := c.Query("id")
	state := c.Query("state")
	result, _ := common.Db.Exec("update works set status=? where id=?", state, id)
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
