package admin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"photo/src/function"
	"photo/src/model/common"
	"strings"
)

//获取提现订单
func GetDrawingLisst(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	code := c.Query("code")
	account := c.Query("account")
	if code != "" {
		where += "d.code='" + code + "' && "
	}
	if account != "" {
		var aid int64
		_ = common.Db.QueryRow("select id from agent where account=?", account).Scan(&aid)
		where += "d.aid=" + fmt.Sprintf("%v", aid) + " && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数--------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from drawing d left join agent a on d.aid = a.id " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("select d.*,a.account from drawing d left join agent a on d.aid = a.id " + wheremap + " order by d.id desc " + limit)
	defer rows.Close()
	resultAll := common.AssemblyData(rows)
	return resultAll, count
}

//处理提现状态
func StateDrawing(c *gin.Context) (string, int) {
	var msg string
	var code int
	id := c.Query("id")
	state := c.Query("state")
	if state == "0" {
		msg = "操作失败"
		code = 400
	} else {
		result, _ := common.Db.Exec("update drawing set status=? where id=?", state, id)
		rowCount, err := result.RowsAffected()
		function.CheckErr("更改提现状态失败", err)
		if rowCount > 0 {
			msg = "设置成功"
			code = 200
		}
	}
	return msg, code
}
