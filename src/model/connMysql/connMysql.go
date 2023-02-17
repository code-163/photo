package connMysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"photo/src/function"
	"time"
)

// DB 数据库连接池
var DB *sql.DB

/**
设置数据库链接
*/
func init() {
	DB, _ = sql.Open("mysql", "photo:photo@tcp(127.0.0.1:3306)/photo?charset=utf8")
	DB.SetConnMaxLifetime(time.Duration(8*3600) * time.Second) //设置可重用链接的最大时间
	DB.SetMaxOpenConns(500)                                    //设置最大的连接数
	DB.SetMaxIdleConns(10)                                     //设置闲置的连接数
	err := DB.Ping()
	function.CheckErr("数据库链接失败", err)
}
