package admin

import (
	"github.com/gin-gonic/gin"
	"photo/src/function"
	"photo/src/model/common"
	"strings"
	"time"
)

// GetWxPayList 获取微信支付列表
func GetWxPayList(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	mchid := c.Query("mchid")
	appid := c.Query("appid")
	status := c.Query("status")
	if mchid != "" {
		where += "POSITION('" + mchid + "' IN mchid) && "
	}
	if appid != "" {
		where += "POSITION('" + appid + "' IN appid) && "
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
	_ = common.Db.QueryRow("select count(*) from wxPayList " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("select * from wxPayList " + wheremap + " order by id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// AddWxPay 添加微信支付
func AddWxPay(c *gin.Context) (string, int) {
	var msg string
	var code int
	mchid := c.PostForm("mchid")
	mchnumber := c.PostForm("mchnumber")
	mchapiv3key := c.PostForm("mchapiv3key")
	mchprivatekey := c.PostForm("mchprivatekey")
	appid := c.PostForm("appid")
	secret := c.PostForm("secret")
	remarks := c.PostForm("remarks")
	right_domain := c.PostForm("right_domain")
	pay_domain := c.PostForm("pay_domain")
	ctime := time.Now().Format("2006-01-02")
	result, _ := common.Db.Exec("insert into wxPayList(`mchid`,`mchnumber`,`mchapiv3key`,`mchprivatekey`,`appid`,`secret`,`remarks`,`ctime`,`right_domain`,`pay_domain`)values(?,?,?,?,?,?,?,?,?,?)", mchid, mchnumber, mchapiv3key, mchprivatekey, appid, secret, remarks, ctime, right_domain, pay_domain)
	affected, _ := result.RowsAffected()
	if affected > 0 {
		msg = "添加成功"
		code = 200
	} else {
		msg = "添加成功"
		code = 200
	}
	return msg, code
}

// EditWxPay 修改微信支付
func EditWxPay(c *gin.Context) (string, int) {
	id := c.PostForm("id")
	mchid := c.PostForm("mchid")
	mchnumber := c.PostForm("mchnumber")
	mchapiv3key := c.PostForm("mchapiv3key")
	mchprivatekey := c.PostForm("mchprivatekey")
	appid := c.PostForm("appid")
	secret := c.PostForm("secret")
	remarks := c.PostForm("remarks")
	right_domain := c.PostForm("right_domain")
	pay_domain := c.PostForm("pay_domain")
	result, _ := common.Db.Exec("update wxPayList set mchid=?,mchnumber=?,mchapiv3key=?,mchprivatekey=?,appid=?,secret=?,remarks=?,right_domain=?,pay_domain=? where id=?", mchid, mchnumber, mchapiv3key, mchprivatekey, appid, secret, remarks, right_domain, pay_domain, id)
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

// DelWxPay 删除微信支付
func DelWxPay(c *gin.Context) (string, int) {
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
	result, _ := common.Db.Exec("delete from wxPayList where " + idval)
	rowCount, err := result.RowsAffected()
	function.CheckErr("删除微信支付失败", err)
	if rowCount > 0 {
		msg = "删除成功"
		code = 200
	}
	return msg, code
}

// StateWxPay 状态开关
func StateWxPay(c *gin.Context) (string, int) {
	var msg string
	var code int
	id := c.Query("id")
	state := c.Query("state")
	result, _ := common.Db.Exec("update wxPayList set status=? where id=?", state, id)
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
