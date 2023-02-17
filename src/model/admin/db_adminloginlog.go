package admin

import (
	"github.com/gin-gonic/gin"
	"photo/src/model/common"
)

/**
管理员登录日志表字段
*/
type adminloginlog struct {
	id         int64  "自增ID"
	account    string "账户"
	loginmsg   string "登录信息，地理位置"
	ip         string "IP地址"
	login_time string "登录时间"
	count      int    "统计数据"
}

// QueryAdminLoginLog 获取管理员登录日志
func QueryAdminLoginLog(c *gin.Context) (interface{}, int) {
	adminloginlog := adminloginlog{}
	limit := common.Pagecut(c)

	//---------统计查询条数--------------------------------------------
	count := common.Db.QueryRow("select count(*) from adminloginlog")
	count.Scan(&adminloginlog.count)
	//---------统计查询条数--------------------------------------------

	row, _ := common.Db.Query("select * from adminloginlog order by id desc " + limit)
	defer row.Close()
	tableData := make([]interface{}, 0)
	for row.Next() {
		row.Scan(&adminloginlog.id, &adminloginlog.account, &adminloginlog.login_time, &adminloginlog.ip, &adminloginlog.loginmsg)
		data := make(map[string]interface{})

		//-----组装数据----------------------------
		data["id"] = adminloginlog.id
		data["account"] = adminloginlog.account
		data["login_time"] = adminloginlog.login_time
		data["ip"] = adminloginlog.ip
		data["loginmsg"] = adminloginlog.loginmsg
		tableData = append(tableData, data)
		//-----组装数据----------------------------
	}
	//resjson, _ := json.Marshal(tableData)
	return tableData, adminloginlog.count
}
