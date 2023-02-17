package admin

import (
	"github.com/gin-gonic/gin"
	"photo/src/function"
	"photo/src/model/common"
	"strings"
	"time"
)

// GetDoorDomainList 获取入口域名列表
func GetDoorDomainList(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	domain := c.Query("domain")
	status := c.Query("status")
	if domain != "" {
		where += "POSITION('" + domain + "' IN domain) && "
	}
	if status != "" {
		where += "status=" + status + " && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数-----------------------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from door_domain " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("select * from door_domain " + wheremap + " order by id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// GetJumpDomainList 获取中转域名列表
func GetJumpDomainList(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	domain := c.Query("domain")
	status := c.Query("status")
	if domain != "" {
		where += "POSITION('" + domain + "' IN domain) && "
	}
	if status != "" {
		where += "status=" + status + " && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数-----------------------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from jump_domain " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("select * from jump_domain " + wheremap + " order by id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// GetLuodiDomainList 获取落地域名列表
func GetLuodiDomainList(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	domain := c.Query("domain")
	status := c.Query("status")
	if domain != "" {
		where += "POSITION('" + domain + "' IN domain) && "
	}
	if status != "" {
		where += "status=" + status + " && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数-----------------------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from luodi_domain " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("select * from luodi_domain " + wheremap + " order by id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// AddDoorDomain 添加入口域名
func AddDoorDomain(c *gin.Context) (string, int) {
	domainType := c.PostForm("type")
	domain := c.PostForm("domain")
	domainArr := strings.Split(domain, "\n")
	ctime := time.Now().Format("2006-01-02 15:04:05")
	for i := 0; i < len(domainArr); i++ {
		_, _ = common.Db.Exec("insert into door_domain(`domain`,`ctime`)values(?,?)", domainType+domainArr[i], ctime)
	}
	msg := "添加成功"
	code := 200
	return msg, code
}

// AddJumpDomain 添加中转域名
func AddJumpDomain(c *gin.Context) (string, int) {
	domainType := c.PostForm("type")
	domain := c.PostForm("domain")
	domainArr := strings.Split(domain, "\n")
	ctime := time.Now().Format("2006-01-02 15:04:05")
	for i := 0; i < len(domainArr); i++ {
		_, _ = common.Db.Exec("insert into jump_domain(`domain`,`ctime`)values(?,?)", domainType+domainArr[i], ctime)
	}
	msg := "添加成功"
	code := 200
	return msg, code
}

// AddLuodiDomain 添加落地域名
func AddLuodiDomain(c *gin.Context) (string, int) {
	domainType := c.PostForm("type")
	domain := c.PostForm("domain")
	domainArr := strings.Split(domain, "\n")
	ctime := time.Now().Format("2006-01-02 15:04:05")
	for i := 0; i < len(domainArr); i++ {
		_, _ = common.Db.Exec("insert into luodi_domain(`domain`,`ctime`)values(?,?)", domainType+domainArr[i], ctime)
	}
	msg := "添加成功"
	code := 200
	return msg, code
}

// DelDoorDomain 删除入口域名
func DelDoorDomain(c *gin.Context) (string, int) {
	var msg string
	var code int
	id := c.Query("id")
	ids := c.Query("ids")
	var idval string
	if id != "" {
		idval = "id=" + id
	} else if ids != "" {
		idval = "id in(" + ids + ")"
	}
	result, _ := common.Db.Exec("delete from door_domain where " + idval)
	rowCount, err := result.RowsAffected()
	function.CheckErr("删除域名失败", err)
	if rowCount > 0 {
		msg = "删除成功"
		code = 200
	}
	return msg, code
}

// DelJumpDomain 删除中转域名
func DelJumpDomain(c *gin.Context) (string, int) {
	var msg string
	var code int
	id := c.Query("id")
	ids := c.Query("ids")
	var idval string
	if id != "" {
		idval = "id=" + id
	} else if ids != "" {
		idval = "id in(" + ids + ")"
	}
	result, _ := common.Db.Exec("delete from jump_domain where " + idval)
	rowCount, err := result.RowsAffected()
	function.CheckErr("删除域名失败", err)
	if rowCount > 0 {
		msg = "删除成功"
		code = 200
	}
	return msg, code
}

// DelLuodiDomain 删除落地域名
func DelLuodiDomain(c *gin.Context) (string, int) {
	var msg string
	var code int
	id := c.Query("id")
	ids := c.Query("ids")
	var idval string
	if id != "" {
		idval = "id=" + id
	} else if ids != "" {
		idval = "id in(" + ids + ")"
	}
	result, _ := common.Db.Exec("delete from luodi_domain where " + idval)
	rowCount, err := result.RowsAffected()
	function.CheckErr("删除域名失败", err)
	if rowCount > 0 {
		msg = "删除成功"
		code = 200
	}
	return msg, code
}
