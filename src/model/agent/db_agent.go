package agent

import (
	"github.com/gin-gonic/gin"
	"photo/src/function"
	"photo/src/model/common"
	"photo/src/model/h5"
	"strconv"
	"strings"
	"time"
)

// GetAgentLevel 获取二级代理
func GetAgentLevel(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	account := c.Query("account")
	pid := c.Query("pid")
	if account != "" {
		where += "POSITION('" + account + "' IN a.account) && "
	}
	where += "a.pid = " + pid + " && "
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数--------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from agent as a " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("select a.*,IFNULL(o.zong,0) as sales from agent as a LEFT JOIN (SELECT sum(price) as zong,code FROM orderbuy where status=1 GROUP BY code) as o on a.code=o.code " + wheremap + " order by a.id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// AddAgentLevel 添加二级代理
func AddAgentLevel(c *gin.Context) (string, int) {
	var msg string
	var code int
	var myPoint int
	pid := c.Query("pid")
	account := c.PostForm("account")
	passwd := c.PostForm("passwd")
	point, _ := strconv.Atoi(c.PostForm("point"))
	remarks := c.PostForm("remarks")
	plink := c.PostForm("plink")
	ctime := time.Now().Format("2006-01-02")
	coder := function.CreateRandomString(6)
	common.Db.QueryRow("select point from agent where id=?", pid).Scan(&myPoint)
	if myPoint <= point {
		msg = "添加失败，下级点数不能大于或等于自己"
		code = 400
	} else {
		result, _ := common.Db.Exec("insert into agent(`account`,`passwd`,`point`,`remarks`,`ctime`,`code`,`pid`,`plink`)values(?,?,?,?,?,?,?,?)", account, passwd, point, remarks, ctime, coder, pid, plink)
		rowCount, _ := result.RowsAffected()
		if rowCount > 0 {
			msg = "添加成功"
			code = 200
		} else {
			msg = "添加失败"
			code = 400
		}
	}
	return msg, code
}

// EditAgentLevel 修改二级代理
func EditAgentLevel(c *gin.Context) (string, int) {
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

// DelAgentLevel 删除二级代理
func DelAgentLevel(c *gin.Context) (string, int) {
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

// Agentlogin 代理登录
func Agentlogin(c *gin.Context) (string, int, string, string, int64, string, string, int) {
	var msg string
	var code int
	var agentcodeRes, agenttokenRes string
	var id int64
	var accountRes string
	var coder string
	var pid int
	account := c.PostForm("account")
	passwd := c.PostForm("passwd")
	//timer := time.Now().Format("2006-01-02 15:04:05")
	//查询用户账户
	user_row := common.Db.QueryRow("select id,account,code,pid from agent where account=? && passwd=?", account, passwd)
	user_row.Scan(&id, &accountRes, &coder, &pid)
	if id > 0 && accountRes != "" && coder != "" {

		//--------插入登录日志信息-----------------------------------------------------------------------------------------------
		/*ip := c.ClientIP()
		resp, _ := http.Get("https://67ip.cn/check?ip=" + ip + "&token=24194ac2fb7be46f8bc8ea2f980bc5a6")
		defer resp.Body.Close()
		ipdata := make(map[string]interface{})
		body, _ := io.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &ipdata)
		res := ipdata["data"].(map[string]interface{})
		loginmsg := res["country"].(string) + res["province"].(string) + res["city"].(string) + res["service"].(string)
		_, _ = common.Db.Exec("insert into agentloginlog (`aid`,`ip`,`loginmsg`,`login_time`) values (?,?,?,?)", id, ip, loginmsg, timer)*/
		//--------插入登录日志信息-----------------------------------------------------------------------------------------------

		agentcodeRes = function.CreateRandomString(8)
		agenttokenRes = function.CreateRandomString(10)
		redisdb, ctx, _ := function.CreateRedisClient()
		err := redisdb.Set(ctx, agentcodeRes, agenttokenRes, 28800*time.Second).Err()
		function.CheckErr("设置redis数据失败", err)
		msg = "登录成功"
		code = 0
	} else {
		msg = "登录失败，账号或密码不正确"
		code = 400
	}
	return msg, code, agentcodeRes, agenttokenRes, id, accountRes, coder, pid
}

// Agentlogout 代理退出
func Agentlogout(c *gin.Context) (string, int) {
	var msg string
	var code int
	action := c.PostForm("action")
	agentcode := c.PostForm("agentcode")
	//timer := time.Now().Format("2006-01-02 15:04:05")
	if action == "logout" {
		redisdb, ctx, _ := function.CreateRedisClient()
		err := redisdb.Del(ctx, agentcode).Err()
		function.CheckErr("删除redis数据失败", err)
		msg = "退出成功"
		code = 0
	}
	return msg, code
}

// CheckAgentlogin 检测用户登录状态
func CheckAgentlogin(c *gin.Context) (string, int) {
	var msg string
	var code int
	agentcode := c.PostForm("agentcode")
	agenttoken := c.PostForm("agenttoken")
	redisdb, ctx, _ := function.CreateRedisClient()
	val, _ := redisdb.Get(ctx, agentcode).Result()
	if val != "" && val == agenttoken {
		msg = "你已登录过"
		code = 0
	} else {
		msg = "登录超时，请重新登录"
		code = 400
	}
	return msg, code
}

// GetAgentOrder 当前用户读取订单
func GetAgentOrder(c *gin.Context) ([]map[string]string, int) {
	limit := common.Pagecut(c)

	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	orderNo := c.Query("order_no")
	status := c.Query("status")
	code := c.Query("code")
	if orderNo != "" {
		where += "order_no='" + orderNo + "' && "
	}
	if code != "" {
		where += "code='" + code + "' && "
	}
	if status != "" {
		where += "status='" + status + "' && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数--------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from orderbuy " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------

	rows, _ := common.Db.Query("select * from orderbuy " + wheremap + " order by id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// GetAgentOrderLevel 当前下级代理订单
func GetAgentOrderLevel(c *gin.Context) ([]map[string]string, int) {
	limit := common.Pagecut(c)

	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	orderNo := c.Query("order_no")
	status := c.Query("status")
	code := c.Query("code")

	var id int64
	common.Db.QueryRow("SELECT id FROM agent where code=?", code).Scan(&id)
	query, _ := common.Db.Query("SELECT code FROM `agent` where pid=?", id)
	defer query.Close()
	reArr := common.AssemblyData(query)
	codeArr := make([]string, 0)
	for _, val := range reArr {
		codeArr = append(codeArr, val["code"])
	}
	codeString := strings.Join(codeArr, ",")
	where += "FIND_IN_SET(o.code,'" + codeString + "') && "

	if orderNo != "" {
		where += "o.order_no='" + orderNo + "' && "
	}
	if status != "" {
		where += "o.status='" + status + "' && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单-----------------------------------------------

	//---------统计查询条数--------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from orderbuy as o " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------

	rows, _ := common.Db.Query("select o.*,a.account,a.remarks from orderbuy as o left join agent as a on o.code=a.code " + wheremap + " order by id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// GetAgentDrawing 当前用户读取提现订单
func GetAgentDrawing(c *gin.Context) ([]map[string]string, int) {
	limit := common.Pagecut(c)

	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	aid := c.Query("aid")
	if aid != "" {
		where += "aid='" + aid + "' && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数--------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from drawing " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------

	rows, _ := common.Db.Query("select * from drawing " + wheremap + " order by id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// GetAgentInfo 当前用户读取账户
func GetAgentInfo(c *gin.Context) map[string]string {
	_, _, domain := h5.GetDoorDomain(c)
	aid := c.Query("aid")
	userRow, _ := common.Db.Query("select * from agent where id=?", aid)
	defer userRow.Close()
	result := common.AssemblyData(userRow)
	result[0]["domain"] = domain
	return result[0]
}

// AddAgentDrawing 当前用户提现操作
func AddAgentDrawing(c *gin.Context) (string, int) {
	var msg string
	var code int
	aid := c.PostForm("aid")
	alipayAccount := c.PostForm("alipay_account")
	alipayUsername := c.PostForm("alipay_username")
	price, _ := strconv.ParseFloat(c.PostForm("price"), 64)

	var wallet float64
	var point int
	common.Db.QueryRow("select wallet,point from agent where id=?", aid).Scan(&wallet, &point)

	if wallet < price {
		msg = "提现失败，余额不足"
		code = 400
	} else {
		if price >= 100 {
			coder := function.CreateRandomString(8)
			ctime := time.Now().Format("2006-01-02 15:04:05")
			//money := price * ((100 - float64(point)) * 0.01)
			//money := price * (float64(point) * 0.01)
			result, _ := common.Db.Exec("insert into drawing (`code`,`price`,`aid`,`ctime`,`alipay_account`,`money`,`alipay_username`)values (?,?,?,?,?,?,?)", coder, price, aid, ctime, alipayAccount, price, alipayUsername)
			rowCount, err := result.RowsAffected()
			function.CheckErr("用户提现失败", err)
			if rowCount > 0 {
				result_user, _ := common.Db.Exec("update agent set wallet=? where id=?", wallet-price, aid)
				rowCount_user, _ := result_user.RowsAffected()
				if rowCount_user > 0 {
					msg = "提现操作成功"
					code = 0
				} else {
					msg = "提现异常，扣除余额失败，请联系管理员"
					code = 400
				}
			} else {
				msg = "提现异常，请联系管理员"
				code = 400
			}
		} else {
			msg = "提现失败，提现金额不能小于100"
			code = 400
		}
	}
	return msg, code
}

// QueryAgentSales 当前用户读取销售数据
func QueryAgentSales(c *gin.Context) (string, int, map[string]interface{}) {
	var msg string
	var code int
	coder := c.Query("code")
	data := make(map[string]interface{})
	if coder != "" {

		//当天统计
		dayStartTime := function.GetDayStartTime()
		dayEndTime := function.GetDayEndTime()
		todaySales, _ := common.Db.Query("select coalesce(sum(price), 0) as price, count(*) as count from orderbuy where code=? && status=1 && (pay_time >=? && pay_time <=?)", coder, dayStartTime, dayEndTime)
		defer todaySales.Close()
		resultToDay := common.AssemblyData(todaySales)
		data["toDaySales"] = resultToDay[0]["price"]
		data["toDayCountOrder"] = resultToDay[0]["count"]

		//昨天统计
		yesterdayStartTime := function.GetYesterdayStartTime()
		yesterdayEndTime := function.GetYesterdayEndTime()
		yesterDaySales, _ := common.Db.Query("select coalesce(sum(price), 0) as price, count(*) as count from orderbuy where code=? && status=1 && (pay_time >=? && pay_time <=?)", coder, yesterdayStartTime, yesterdayEndTime)
		defer yesterDaySales.Close()
		resultYesterDay := common.AssemblyData(yesterDaySales)
		data["yesterDaySales"] = resultYesterDay[0]["price"]
		data["yesterDayCountOrder"] = resultYesterDay[0]["count"]

		//7天统计
		weekStartTime := function.GetFirstDateOfWeek()
		weekEndTime := function.GetLastWeekFirstDate()
		weekSales, _ := common.Db.Query("select coalesce(sum(price), 0) as price, count(*) as count from orderbuy where code=? && status=1 && (pay_time >=? && pay_time <?)", coder, weekStartTime, weekEndTime)
		defer weekSales.Close()
		resultWeekSales := common.AssemblyData(weekSales)
		data["weekSales"] = resultWeekSales[0]["price"]
		data["weekCountOrder"] = resultWeekSales[0]["count"]

		//当月统计
		monthStartTime := function.GetFirstDateOfMonth(time.Now())
		monthEndTime := function.GetLastDateOfMonth(time.Now())
		monthSales, _ := common.Db.Query("select coalesce(sum(price), 0) as price, count(*) as count from orderbuy where code=? && status=1 && (pay_time >=? && pay_time <=?)", coder, monthStartTime, monthEndTime)
		defer monthSales.Close()
		resultMonthSales := common.AssemblyData(monthSales)
		data["monthSales"] = resultMonthSales[0]["price"]
		data["monthCountOrder"] = resultMonthSales[0]["count"]

		//总订单数量
		totalSales, _ := common.Db.Query("select count(*) as count from orderbuy where code=?", coder)
		defer totalSales.Close()
		resultTotalCountOrder := common.AssemblyData(totalSales)
		data["totalCountOrder"] = resultTotalCountOrder[0]["count"]

		//总成功订单数
		successSales, _ := common.Db.Query("select count(*) as count from orderbuy where code=? && status=1", coder)
		defer successSales.Close()
		resultSuccessSalesCountOrder := common.AssemblyData(successSales)
		data["successSalesCountOrder"] = resultSuccessSalesCountOrder[0]["count"]

		code = 200
		msg = "获取数据成功"
	} else {
		code = 400
		msg = "参数错误，缺少coder"
	}
	return msg, code, data
}

// GetAgentLoginlog 获取代理登录日志
func GetAgentLoginlog(c *gin.Context) ([]map[string]string, int) {
	var count int
	var result []map[string]string
	limit := common.Pagecut(c)
	aid := c.Query("aid")
	if aid != "" {
		_ = common.Db.QueryRow("select count(*) from agentloginlog where uid=?", aid).Scan(&count)
		row, _ := common.Db.Query("select * from agentloginlog where aid=? order by id desc "+limit, aid)
		defer row.Close()
		result = common.AssemblyData(row)
	}
	return result, count
}

// GetVisitAgent 访问统计
func GetVisitAgent(c *gin.Context) (string, int, map[string]interface{}) {
	var msg string
	var code int
	data := make(map[string]interface{})
	coder := c.Query("code")

	//当天统计
	dayStartTime := time.Now().Format("2006-01-02")
	var toDayPv int
	_ = common.Db.QueryRow("select count(*) as count from `visit` where code=? && types=0 && ctime=?", coder, dayStartTime).Scan(&toDayPv)
	var toDayUv int
	_ = common.Db.QueryRow("SELECT count(DISTINCT(ip)) as count FROM `visit` WHERE code=? && types=0 && ctime=?", coder, dayStartTime).Scan(&toDayUv)
	data["todayVisitPv"] = toDayPv
	data["todayVisitUv"] = toDayUv

	//昨天统计
	yesterdayTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var yesterdayPv int
	_ = common.Db.QueryRow("select count(*) as count from `visit` where code=? && types=0 && ctime=?", coder, yesterdayTime).Scan(&yesterdayPv)
	var yesterdayUv int
	_ = common.Db.QueryRow("SELECT count(DISTINCT(ip)) as count FROM `visit` WHERE code=? && types=0 && ctime=?", coder, yesterdayTime).Scan(&yesterdayUv)
	data["yesterdayVisitPv"] = yesterdayPv
	data["yesterdayVisitUv"] = yesterdayUv

	//7天统计
	weekStartTime := function.GetFirstDateOfWeek02()
	weekEndTime := function.GetLastWeekFirstDate02()
	var weekPv int
	_ = common.Db.QueryRow("select count(*) as count from `visit` where code=? && types=0 && (ctime >=? && ctime <=?)", coder, weekStartTime, weekEndTime).Scan(&weekPv)
	var weekUv int
	_ = common.Db.QueryRow("SELECT count(DISTINCT(ip)) as count FROM `visit` WHERE code=? && types=0 && (ctime >=? && ctime <=?)", coder, weekStartTime, weekEndTime).Scan(&weekUv)
	data["weekVisitPv"] = weekPv
	data["weekVisitUv"] = weekUv

	//当月统计
	monthStartTime := function.GetFirstDateOfMonth02(time.Now())
	monthEndTime := function.GetLastDateOfMonth02(time.Now())
	var monthPv int
	_ = common.Db.QueryRow("select count(*) as count from `visit` where code=? && types=0 && (ctime >=? && ctime <=?)", coder, monthStartTime, monthEndTime).Scan(&monthPv)
	var monthUv int
	_ = common.Db.QueryRow("SELECT count(DISTINCT(ip)) as count FROM `visit` WHERE code=? && types=0 && (ctime >=? && ctime <=?)", coder, monthStartTime, monthEndTime).Scan(&monthUv)
	data["monthVisitPv"] = monthPv
	data["monthVisitUv"] = monthUv

	code = 200
	msg = "获取数据成功"
	return msg, code, data
}
