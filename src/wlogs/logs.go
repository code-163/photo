package wlogs

import (
	"os"
	"time"
)

const (
	LOGPATH  = "logs/"
	FORMAT   = "20060102"
	LineFeed = "\r\n"
)

/**
写入日志
*/
func WriteLog(fileName, msg string) error {
	//以天为基准,存日志
	path := LOGPATH + time.Now().Format(FORMAT) + "/"
	if !IsExist(path) {
		_ = CreateDir(path)
	}
	time_name := "_" + time.Now().Format("2006-01-02(15)") + ".log"
	timer := time.Now().Format("2006-01-02 15:04:05")
	fileObj, err := os.OpenFile(path+fileName+time_name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	defer fileObj.Close()
	_, err = fileObj.WriteString(timer + " => " + msg + LineFeed)
	return err
}

/**
文件夹创建
*/
func CreateDir(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	os.Chmod(path, os.ModePerm)
	return nil
}

/**
判断文件夹/文件是否存在  存在返回 true
*/
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}
