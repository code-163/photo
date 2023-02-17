package admin

import (
	"github.com/gin-gonic/gin"
	"photo/src/function"
	"photo/src/model/common"
	"strings"
	"time"
)

// GetOrderList 获取订单列表
func GetOrderList(c *gin.Context) ([]map[string]string, int) {
	//---------查询订单--------------------------------------------------------------
	var where, wheremap, mapres string
	order_no := c.Query("order_no")
	status := c.Query("status")
	if order_no != "" {
		where += "POSITION('" + order_no + "' IN o.order_no) && "
	}
	if status != "" {
		where += "o.status=" + status + " && "
	}
	if mapres = strings.TrimRight(where, "&& "); mapres != "" {
		wheremap = "where " + mapres
	}
	//---------查询订单--------------------------------------------------------------

	//---------统计查询条数--------------------------------------------
	var count int
	_ = common.Db.QueryRow("select count(*) from orderbuy as o " + wheremap).Scan(&count)
	//---------统计查询条数--------------------------------------------
	limit := common.Pagecut(c)
	rows, _ := common.Db.Query("SELECT o.*,a.account,u.username FROM `orderbuy` as o LEFT JOIN agent as a on o.code=a.code left join user as u on o.uid=u.id " + wheremap + " order by id desc " + limit)
	defer rows.Close()
	result := common.AssemblyData(rows)
	return result, count
}

// QueryOrderSales 读取系统所有销量
func QueryOrderSales(c *gin.Context) (string, int, map[string]interface{}) {
	var msg string
	var code int
	data := make(map[string]interface{})

	//当天统计
	dayStartTime := function.GetDayStartTime()
	dayEndTime := function.GetDayEndTime()
	todaySales, _ := common.Db.Query("select coalesce(sum(price), 0) as price, count(*) as count from orderbuy where status=1 && (pay_time >=? && pay_time <=?)", dayStartTime, dayEndTime)
	defer todaySales.Close()
	resultToDay := common.AssemblyData(todaySales)
	data["toDaySales"] = resultToDay[0]["price"]
	data["toDayCountOrder"] = resultToDay[0]["count"]

	//昨天统计
	yesterdayStartTime := function.GetYesterdayStartTime()
	yesterdayEndTime := function.GetYesterdayEndTime()
	yesterDaySales, _ := common.Db.Query("select coalesce(sum(price), 0) as price, count(*) as count from orderbuy where status=1 && (pay_time >=? && pay_time <=?)", yesterdayStartTime, yesterdayEndTime)
	defer yesterDaySales.Close()
	resultYesterDay := common.AssemblyData(yesterDaySales)
	data["yesterDaySales"] = resultYesterDay[0]["price"]
	data["yesterDayCountOrder"] = resultYesterDay[0]["count"]

	//7天统计
	sevenDayStartTime := time.Now().AddDate(0, 0, -6).Format("2006-01-02") + " 00:00:00"
	sevenDayEndTime := time.Now().AddDate(0, 0, 0).Format("2006-01-02") + " 23:59:59"
	sevenDaySales, _ := common.Db.Query("select coalesce(sum(price), 0) as price, count(*) as count from orderbuy where status=1 && (pay_time >=? && pay_time <=?)", sevenDayStartTime, sevenDayEndTime)
	defer sevenDaySales.Close()
	resultSevenDay := common.AssemblyData(sevenDaySales)
	data["sevenDaySales"] = resultSevenDay[0]["price"]
	data["sevenDayCountOrder"] = resultSevenDay[0]["count"]

	//本周统计
	weekStartTime := function.GetFirstDateOfWeek()
	weekEndTime := function.GetLastWeekFirstDate()
	weekSales, _ := common.Db.Query("select coalesce(sum(price), 0) as price, count(*) as count from orderbuy where status=1 && (pay_time >=? && pay_time <?)", weekStartTime, weekEndTime)
	defer weekSales.Close()
	resultWeekSales := common.AssemblyData(weekSales)
	data["weekSales"] = resultWeekSales[0]["price"]
	data["weekCountOrder"] = resultWeekSales[0]["count"]

	//当月统计
	monthStartTime := function.GetFirstDateOfMonth(time.Now())
	monthEndTime := function.GetLastDateOfMonth(time.Now())
	monthSales, _ := common.Db.Query("select coalesce(sum(price), 0) as price, count(*) as count from orderbuy where status=1 && (pay_time >=? && pay_time <=?)", monthStartTime, monthEndTime)
	defer monthSales.Close()
	resultMonthSales := common.AssemblyData(monthSales)
	data["monthSales"] = resultMonthSales[0]["price"]
	data["monthCountOrder"] = resultMonthSales[0]["count"]

	//总销量
	totalSales, _ := common.Db.Query("select coalesce(sum(price), 0) as price from orderbuy where status=1")
	defer totalSales.Close()
	resultTotalCountOrder := common.AssemblyData(totalSales)
	data["totalCountOrder"] = resultTotalCountOrder[0]["price"]

	//总成功订单数
	successSales, _ := common.Db.Query("select count(*) as count from orderbuy where status=1")
	defer successSales.Close()
	resultSuccessSalesCountOrder := common.AssemblyData(successSales)
	data["successSalesCountOrder"] = resultSuccessSalesCountOrder[0]["count"]

	code = 200
	msg = "获取数据成功"
	return msg, code, data
}

// GetVisit 访问统计
func GetVisit(c *gin.Context) (string, int, map[string]interface{}) {
	var msg string
	var code int
	data := make(map[string]interface{})

	//当天统计
	dayStartTime := time.Now().Format("2006-01-02")
	var toDayPv int
	_ = common.Db.QueryRow("select count(*) as count from `visit` where types=0 && ctime=?", dayStartTime).Scan(&toDayPv)
	var toDayUv int
	_ = common.Db.QueryRow("SELECT count(DISTINCT(ip)) as count FROM `visit` WHERE types=0 && ctime=?", dayStartTime).Scan(&toDayUv)
	data["todayVisitPv"] = toDayPv
	data["todayVisitUv"] = toDayUv

	//昨天统计
	yesterdayTime := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	var yesterdayPv int
	_ = common.Db.QueryRow("select count(*) as count from `visit` where types=0 && ctime=?", yesterdayTime).Scan(&yesterdayPv)
	var yesterdayUv int
	_ = common.Db.QueryRow("SELECT count(DISTINCT(ip)) as count FROM `visit` WHERE types=0 && ctime=?", yesterdayTime).Scan(&yesterdayUv)
	data["yesterdayVisitPv"] = yesterdayPv
	data["yesterdayVisitUv"] = yesterdayUv

	//7天统计
	weekStartTime := function.GetFirstDateOfWeek02()
	weekEndTime := function.GetLastWeekFirstDate02()
	var weekPv int
	_ = common.Db.QueryRow("select count(*) as count from `visit` where types=0 && (ctime >=? && ctime <=?)", weekStartTime, weekEndTime).Scan(&weekPv)
	var weekUv int
	_ = common.Db.QueryRow("SELECT count(DISTINCT(ip)) as count FROM `visit` WHERE types=0 && (ctime >=? && ctime <=?)", weekStartTime, weekEndTime).Scan(&weekUv)
	data["weekVisitPv"] = weekPv
	data["weekVisitUv"] = weekUv

	//当月统计
	monthStartTime := function.GetFirstDateOfMonth02(time.Now())
	monthEndTime := function.GetLastDateOfMonth02(time.Now())
	var monthPv int
	_ = common.Db.QueryRow("select count(*) as count from `visit` where types=0 && (ctime >=? && ctime <=?)", monthStartTime, monthEndTime).Scan(&monthPv)
	var monthUv int
	_ = common.Db.QueryRow("SELECT count(DISTINCT(ip)) as count FROM `visit` WHERE types=0 && (ctime >=? && ctime <=?)", monthStartTime, monthEndTime).Scan(&monthUv)
	data["monthVisitPv"] = monthPv
	data["monthVisitUv"] = monthUv

	code = 200
	msg = "获取数据成功"
	return msg, code, data
}
