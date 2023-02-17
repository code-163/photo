package admin

import (
	"github.com/gin-gonic/gin"
	"photo/src/function"
	"photo/src/model/common"
	"time"
)

/**
admin表字段
*/
type admin struct {
	id         int64  "自增ID"
	account    string "账户"
	passwd     string "密码"
	login_time string "登录时间"
	count      int    "统计数据"
}

// EditAdminPasswd 修改管理员密码
func EditAdminPasswd(c *gin.Context) (string, int) {
	id := c.PostForm("id")
	account := c.PostForm("account")
	passwd := c.PostForm("passwd")
	point := c.PostForm("point")
	remarks := c.PostForm("remarks")

	result, _ := common.Db.Exec("update admin set account=?,passwd=?,point=?,remarks=? where id=?", account, passwd, point, remarks, id)
	rowCount, err := result.RowsAffected()
	function.CheckErr("修改用户失败", err)
	var msg string
	var code int
	if rowCount > 0 {
		msg = "修改成功"
		code = 200
	}
	return msg, code
}

// Checkadminlogin 检测登录状态
func Checkadminlogin(c *gin.Context) (string, int) {
	var msg string
	var code int
	admincode := c.PostForm("admincode")
	admintoken := c.PostForm("admintoken")
	redisdb, ctx, _ := function.CreateRedisClient()
	val, _ := redisdb.Get(ctx, admincode).Result() //获取admincode
	if val != "" && val == admintoken {
		msg = "你已登录过"
		code = 0
	} else {
		msg = "登录超时，请重新登录"
		code = 400
	}
	return msg, code
}

// Adminlogin 管理员登录
func Adminlogin(c *gin.Context) (string, int, string, string, int64, string) {
	var msg string
	var code int
	var admincode_res, admintoken_res string
	account := c.PostForm("account")
	passwd := c.PostForm("passwd")
	//timer := time.Now().Format("2006-01-02 15:04:05")
	//查询管理员账户
	admin := admin{}
	user_row := common.Db.QueryRow("select id,account,passwd from admin where account=? && passwd=?", account, passwd)
	_ = user_row.Scan(&admin.id, &admin.account, &admin.passwd)
	if admin.id > 0 && admin.account != "" && admin.passwd != "" {

		//--------插入登录日志信息-----------------------------------------------------------------------------------------------
		/*ip := c.ClientIP()
		resp, _ := http.Get("https://67ip.cn/check?ip=" + ip + "&token=24194ac2fb7be46f8bc8ea2f980bc5a6")
		defer resp.Body.Close()
		ipdata := make(map[string]interface{})
		body, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal([]byte(string(body)), &ipdata)
		res := ipdata["data"].(map[string]interface{})
		loginmsg := res["country"].(string) + res["province"].(string) + res["city"].(string) + res["service"].(string)
		_, _ = common.Db.Exec("insert into adminloginlog (`account`,`ip`,`loginmsg`,`login_time`) values (?,?,?,?)", admin.account, ip, loginmsg, timer)*/
		//--------插入登录日志信息-----------------------------------------------------------------------------------------------

		admincode_res = function.CreateRandomString(8)
		admintoken_res = function.CreateRandomString(10)
		redisdb, ctx, _ := function.CreateRedisClient()
		err := redisdb.Set(ctx, admincode_res, admintoken_res, 28800*time.Second).Err()
		function.CheckErr("设置redis数据失败", err)
		msg = "登录成功"
		code = 0
	} else {
		msg = "登录失败，账号或密码不正确"
		code = 400
	}
	return msg, code, admincode_res, admintoken_res, admin.id, admin.account
}

// Adminlogout 管理员退出登录
func Adminlogout(c *gin.Context) (string, int) {
	var msg string
	var code int
	action := c.PostForm("action")
	admincode := c.PostForm("admincode")
	//timer := time.Now().Format("2006-01-02 15:04:05")
	if action == "logout" {
		redisdb, ctx, _ := function.CreateRedisClient()
		err := redisdb.Del(ctx, admincode).Err()
		function.CheckErr("删除redis数据失败", err)
		msg = "退出成功"
		code = 0
	}
	return msg, code
}
