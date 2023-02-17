package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"photo/src/function"
	"photo/src/model/common"
	"strings"
	"time"
)

const qiNiuUrl = "https://7niu.trumall.cn/"

var DBB *sql.DB

//设置数据库链接
func init() {
	DBB, _ = sql.Open("mysql", "photo:photo@tcp(193.134.209.55:3306)/photo?charset=utf8")
	DBB.SetConnMaxLifetime(time.Duration(8*3600) * time.Second) //设置可重用链接的最大时间
	DBB.SetMaxOpenConns(500)                                    //设置最大的连接数
	DBB.SetMaxIdleConns(10)                                     //设置闲置的连接数
	err := DBB.Ping()
	function.CheckErr("数据库链接失败", err)
}

func editAtlas(val string) string {
	//str := "http://7niu.misihu.cn/image/167146213557559.txt,http://7niu.misihu.cn/image/167146213584467.txt,http://7niu.misihu.cn/image/167146213868132.txt,http://7niu.misihu.cn/image/167146213874662.txt,http://7niu.misihu.cn/image/167146214046476.txt,http://7niu.misihu.cn/image/167146214275988.txt,http://7niu.misihu.cn/image/167146214229974.txt,http://7niu.misihu.cn/image/167146214229487.txt,http://7niu.misihu.cn/image/167146214535678.txt,http://7niu.misihu.cn/image/167146214953572.txt,http://7niu.misihu.cn/image/167146214959818.txt,http://7niu.misihu.cn/image/167146215094473.txt,http://7niu.misihu.cn/image/167146215148529.txt,http://7niu.misihu.cn/image/167146215585232.txt,http://7niu.misihu.cn/image/167146215568435.txt,http://7niu.misihu.cn/image/167146215566925.txt,http://7niu.misihu.cn/image/167146215895859.txt,http://7niu.misihu.cn/image/167146216039589.txt,http://7niu.misihu.cn/image/167146216297197.txt,http://7niu.misihu.cn/image/167146216519334.txt,http://7niu.misihu.cn/image/167146216613336.txt,http://7niu.misihu.cn/image/167146216688374.txt,http://7niu.misihu.cn/image/167146216728157.txt,http://7niu.misihu.cn/image/167146216938795.txt,http://7niu.misihu.cn/image/167146216984186.txt,http://7niu.misihu.cn/image/167146217247984.txt,http://7niu.misihu.cn/image/167146217524941.txt,http://7niu.misihu.cn/image/167146217517552.txt,http://7niu.misihu.cn/image/167146217638487.txt,http://7niu.misihu.cn/image/167146218055918.txt,http://7niu.misihu.cn/image/167146218191512.txt,http://7niu.misihu.cn/image/167146218141121.txt,http://7niu.misihu.cn/image/167146218211143.txt,http://7niu.misihu.cn/image/167146218393375.txt,http://7niu.misihu.cn/image/167146218324151.txt,http://7niu.misihu.cn/image/167146219073282.txt,http://7niu.misihu.cn/image/167146219094624.txt,http://7niu.misihu.cn/image/167146219293123.txt,http://7niu.misihu.cn/image/167146219335851.txt,http://7niu.misihu.cn/image/167146219559683.txt,http://7niu.misihu.cn/image/167146219692481.txt,http://7niu.misihu.cn/image/167146219739448.txt,http://7niu.misihu.cn/image/167146219827337.txt,http://7niu.misihu.cn/image/167146219989189.txt,http://7niu.misihu.cn/image/167146219913358.txt"
	if val != "" {
		resArr := strings.Split(val, ",")
		newArr := make([]string, 0)
		for _, v := range resArr {
			res := editImage(v)
			newArr = append(newArr, res)
		}
		//fmt.Println(strings.Join(newArr, ","))
		return strings.Join(newArr, ",")
	} else {
		return ""
	}
}

func editImage(val string) string {
	//str := "http://7niu.misihu.cn/image/167146211824334.txt"
	if val != "" {
		resArr := strings.Split(val, "/")
		return qiNiuUrl + resArr[3] + "/" + resArr[4]
	} else {
		return ""
	}
}

func updateVideoCover(val string) string {
	if val != "" {
		resArr := strings.Split(val, "/")
		if len(resArr) == 8 {
			return qiNiuUrl + resArr[3] + "/" + resArr[4] + "/" + resArr[5] + "/" + resArr[6] + "/" + resArr[7]
		} else {
			return qiNiuUrl + resArr[3] + "/" + resArr[4]
		}
	} else {
		return ""
	}
}

func updateVideoUrl(val string) string {
	if val != "" {
		resArr := strings.Split(val, "/")
		return qiNiuUrl + resArr[3] + "/" + resArr[4] + "/" + resArr[5] + "/" + resArr[6]
	} else {
		return ""
	}
}

func main() {
	rows, _ := DBB.Query("SELECT * FROM model")
	defer rows.Close()
	result := common.AssemblyData(rows)
	for _, val := range result {
		id := val["id"]
		image := editImage(val["image"])
		exec, _ := DBB.Exec("update model set image=? where id=?", image, id)
		affected, _ := exec.RowsAffected()
		if affected > 0 {
			log.Println("处理成功：", id)
		} else {
			log.Println("处理失败：", id)
		}
	}
}
