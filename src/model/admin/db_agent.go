package admin

import (
	"github.com/gin-gonic/gin"
	"photo/src/function"
	"photo/src/model/common"
	"strconv"
	"strings"
	"time"
)

// GetAgentList 获取代理列表
func GetAgentList(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	account := c.Query("account")
	if account != "" {
		where += "POSITION('" + account + "' IN account) && "
	}
	where += "pid = 0 && "
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数--------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from agent " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("select * from agent " + wheremap + " order by id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// GetLevelAgent 获取二级代理
func GetLevelAgent(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	account := c.Query("account")
	pid := c.Query("pid")
	if account != "" {
		where += "POSITION('" + account + "' IN account) && "
	}
	where += "pid = " + pid + " && "
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数--------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from agent " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("select * from agent " + wheremap + " order by id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// GetDoorDomain 获取入口域名
func GetDoorDomain(c *gin.Context) []map[string]string {
	rows, _ := common.Db.Query("select domain as value,domain as label from door_domain where status=0 && use_status=0 order by id asc")
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result
}

// AddAgent 添加代理
func AddAgent(c *gin.Context) (string, int) {
	account := c.PostForm("account")
	passwd := c.PostForm("passwd")
	point := c.PostForm("point")
	remarks := c.PostForm("remarks")
	plink := c.PostForm("plink")
	ctime := time.Now().Format("2006-01-02")
	coder := function.CreateRandomString(6)
	result, _ := common.Db.Exec("insert into agent(`account`,`passwd`,`point`,`remarks`,`ctime`,`code`,`plink`)values(?,?,?,?,?,?,?)", account, passwd, point, remarks, ctime, coder, plink)
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

// EditAgent 修改代理
func EditAgent(c *gin.Context) (string, int) {
	id := c.PostForm("id")
	account := c.PostForm("account")
	passwd := c.PostForm("passwd")
	point := c.PostForm("point")
	remarks := c.PostForm("remarks")
	plink := c.PostForm("plink")
	result, _ := common.Db.Exec("update agent set account=?,passwd=?,point=?,remarks=?,plink=? where id=?", account, passwd, point, remarks, plink, id)
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

// DelAgent 删除代理
func DelAgent(c *gin.Context) (string, int) {
	id := c.Query("id")
	ids := c.Query("ids")
	var idval string
	if id != "" {
		idval = "id=" + id
	} else if ids != "" {
		idval = "id in(" + ids + ")"
	}

	result, _ := common.Db.Exec("delete from agent where " + idval)
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

// DelLevelAgent 删除二级代理
func DelLevelAgent(c *gin.Context) (string, int) {
	id := c.Query("id")
	ids := c.Query("ids")
	var idval string
	if id != "" {
		idval = "id=" + id
	} else if ids != "" {
		idval = "id in(" + ids + ")"
	}

	result, _ := common.Db.Exec("delete from agent where " + idval)
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

// AliSettlement 结算代理余额
func AliSettlement(c *gin.Context) (string, int) {
	var msg string
	var code int
	coder := c.PostForm("code")
	amount := c.PostForm("price")
	price, _ := strconv.ParseFloat(amount, 64)
	result, _ := common.Db.Exec("update agent set wallet=(wallet-?) where code=?", price, coder)
	rowCount, _ := result.RowsAffected()
	if rowCount > 0 {
		msg = "结算成功"
		code = 200
	} else {
		msg = "结算失败"
		code = 400
	}
	return msg, code
}
