package function

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"github.com/mozillazg/go-pinyin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/tuotoo/qrcode"
	"math"
	"math/big"
	randNew "math/rand"
	"mime/multipart"
	"os"
	"path"
	"photo/src/wlogs"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

// CreateRedisClient 创建Redis链接
func CreateRedisClient() (*redis.Client, context.Context, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	// 通过 cient.Ping() 来检查是否成功连接到了 redis 服务器
	_, err := client.Ping(ctx).Result()
	CheckErr("redis链接失败", err)
	return client, ctx, err
}

//Base64EncodeString 编码
func Base64EncodeString(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

//Base64DecodeString 解码
func Base64DecodeString(str string) (string, []byte) {
	resBytes, _ := base64.StdEncoding.DecodeString(str)
	return string(resBytes), resBytes
}

// Md5 md5加密
func Md5(src string) string {
	m := md5.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

// Sha256 Sha256加密
func Sha256(src string) string {
	m := sha256.New()
	m.Write([]byte(src))
	res := hex.EncodeToString(m.Sum(nil))
	return res
}

// CreateRandomNumber 获取随机数字
func CreateRandomNumber(len int) string {
	var numbers = []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
	var container string
	length := bytes.NewReader(numbers).Len()

	for i := 1; i <= len; i++ {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(length)))
		if err != nil {

		}
		container += fmt.Sprintf("%d", numbers[random.Int64()])
	}
	return container
}

// RandNum 取两个数字之间的随机数
func RandNum(min, max int) int {
	var Rander = randNew.New(randNew.NewSource(time.Now().UnixNano()))
	return Rander.Intn(max-min+1) + min
}

// CreateRandomString 获取随机字符串
func CreateRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

// CheckErr 处理错误信息
func CheckErr(msg string, err error) {
	if err != nil {
		wlogs.WriteLog(msg, err.Error())
	}
}

/**
获取昨天时间零时零分
*/
func GetYesterdayStartTime() string {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	restime := startTime.AddDate(0, 0, -1)
	return restime.Format("2006-01-02 15:04:05")
}

/**
获取昨天24点前时间
*/
func GetYesterdayEndTime() string {
	currentTimer := time.Now()
	endTime := time.Date(currentTimer.Year(), currentTimer.Month(), currentTimer.Day(), 23, 59, 59, 0, currentTimer.Location())
	restime := endTime.AddDate(0, 0, -1)
	return restime.Format("2006-01-02 15:04:05")
}

/**
获取当天时间零时零分
*/
func GetDayStartTime() string {
	currentTime := time.Now()
	startTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	return startTime.Format("2006-01-02 15:04:05")
}

/**
获取当天24点前时间
*/
func GetDayEndTime() string {
	currentTimer := time.Now()
	endTime := time.Date(currentTimer.Year(), currentTimer.Month(), currentTimer.Day(), 23, 59, 59, 0, currentTimer.Location())
	return endTime.Format("2006-01-02 15:04:05")
}

/**
获取本周周一的日期
*/
func GetFirstDateOfWeek() string {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStartDate.Format("2006-01-02 15:04:05")
}

// GetFirstDateOfWeek02 获取本周周一的日期
func GetFirstDateOfWeek02() string {
	now := time.Now()
	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return weekStartDate.Format("2006-01-02")
}

/**
获取下周一的日期
*/
func GetLastWeekFirstDate() string {
	thisWeekMonday := GetFirstDateOfWeek()
	TimeMonday, _ := time.Parse("2006-01-02 15:04:05", thisWeekMonday)
	lastWeekMonday := TimeMonday.AddDate(0, 0, +7)
	return lastWeekMonday.Format("2006-01-02 15:04:05")
}

// GetLastWeekFirstDate02 获取下周一的日期
func GetLastWeekFirstDate02() string {
	thisWeekMonday := GetFirstDateOfWeek()
	TimeMonday, _ := time.Parse("2006-01-02 15:04:05", thisWeekMonday)
	lastWeekMonday := TimeMonday.AddDate(0, 0, +6)
	return lastWeekMonday.Format("2006-01-02")
}

/**
获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
*/
func GetFirstDateOfMonth(d time.Time) string {
	d = d.AddDate(0, 0, -d.Day()+1)
	d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	return d.Format("2006-01-02 15:04:05")
}

// GetFirstDateOfMonth02 获取传入的时间所在月份的第一天
func GetFirstDateOfMonth02(d time.Time) string {
	d = d.AddDate(0, 0, -d.Day()+1)
	d = time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
	return d.Format("2006-01-02")
}

/**
获取传入的时间所在月份的最后一天，即某月最后一天的23:59:59。如传入time.Now(), 返回当前月份的最后一天23点时间。
*/
func GetLastDateOfMonth(d time.Time) string {
	d = d.AddDate(0, 0, -d.Day()+1)
	d = time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location()).AddDate(0, 1, -1)
	return d.Format("2006-01-02 15:04:05")
}

// GetLastDateOfMonth02 获取传入的时间所在月份的最后一天
func GetLastDateOfMonth02(d time.Time) string {
	d = d.AddDate(0, 0, -d.Day()+1)
	d = time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 0, d.Location()).AddDate(0, 1, -1)
	return d.Format("2006-01-02")
}

/**
二维码解码
*/
func QrDecode(qrFile string) string {
	fi, err := os.Open(qrFile)
	if err != nil {
		CheckErr("打开二维码失败", err)
		return ""
	}
	defer fi.Close()
	qrmatrix, err := qrcode.Decode(fi)
	if err != nil {
		CheckErr("二维码解码失败", err)
		return ""
	}
	//fmt.Println(qrmatrix.Content)
	return qrmatrix.Content
}

/**
组装查询的数据并返回数组
*/
func AssemblyData(rows *sql.Rows) []map[string]string {
	//获取列名
	columns, _ := rows.Columns()
	//定义一个切片,长度是字段的个数,切片里面的元素类型是sql.RawBytes
	values := make([]sql.RawBytes, len(columns))
	//定义一个切片,元素类型是interface{} 接口
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		//把sql.RawBytes类型的地址存进去了
		scanArgs[i] = &values[i]
	}
	//获取字段值
	result := make([]map[string]string, 0)
	for rows.Next() {
		res := make(map[string]string)
		rows.Scan(scanArgs...)
		for i, col := range values {
			res[columns[i]] = string(col)
		}
		result = append(result, res)
	}
	return result
}

//中文转拼音
func cutPinyiin(hans string) string {
	var result string
	if hans != "" {
		china := pinyin.NewArgs()
		res := pinyin.Pinyin(hans, china)
		for _, v := range res {
			result += v[0]
		}
	} else {
		result = ""
	}
	return result
}

// SlicePage 计算数组分页
func SlicePage(page, pageSize, nums int64) (sliceStart, sliceEnd int64) {
	if page <= 0 {
		page = 1
	}
	if pageSize < 0 {
		pageSize = 20 //设置一页默认显示的记录数
	}
	if pageSize > nums {
		return 0, nums
	}
	// 总页数
	pageCount := int64(math.Ceil(float64(nums) / float64(pageSize)))
	if page > pageCount {
		return 0, 0
	}
	sliceStart = (page - 1) * pageSize
	sliceEnd = sliceStart + pageSize
	if sliceEnd > nums {
		sliceEnd = nums
	}
	return sliceStart, sliceEnd
}

//pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

//pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

//AesEncrypt 加密
func AesEncrypt(data, sKey string) (string, error) {
	resData := []byte(data)
	key := []byte(sKey)
	//创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//判断加密块的大小
	blockSize := block.BlockSize()
	//填充
	encryptBytes := pkcs7Padding(resData, blockSize)
	//初始化加密数据接收切片
	crypted := make([]byte, len(encryptBytes))
	//使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	res := base64.StdEncoding.EncodeToString(crypted)
	return res, nil
}

//AesDecrypt 解密
func AesDecrypt(data, sKey string) (string, error) {
	resData, _ := base64.StdEncoding.DecodeString(data)
	key := []byte(sKey)
	//创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	//获取块的大小
	blockSize := block.BlockSize()
	//使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//初始化解密数据接收切片
	crypted := make([]byte, len(resData))
	//执行解密
	blockMode.CryptBlocks(crypted, resData)
	//去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return "", err
	}
	return string(crypted), nil
}

// GetPpvodSign 获取ppvod签名
func GetPpvodSign(sKey, ip string) string {
	timer := fmt.Sprintf("%v%v", time.Now().Unix(), "000")
	data := fmt.Sprintf("timestamp=%v&ip=%v", timer, ip)
	origData := []byte(data)
	key, iv := byteToKey(sKey, 16)
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	origData = pkcs5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return fmt.Sprintf("%x", encrypted)
}

func byteToKey(password string, keylen int) ([]byte, []byte) {
	pass := []byte(password)
	prev := []byte{}
	key := []byte{}
	iv := []byte{}
	remain := 0
	for len(key) < keylen {
		hash := md5.Sum(append(prev, pass...))
		remain = keylen - len(key)
		if remain < 16 {
			key = append(key, hash[:remain]...)
		} else {
			key = append(key, hash[:]...)
		}
		prev = hash[:]
	}
	hash := md5.Sum(append(prev, pass...))
	if remain < 16 {
		iv = append(prev[remain:], hash[:remain]...)
	} else {
		iv = hash[:]
	}
	return key, iv
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// GetUpToken 获取七牛云上传凭证
func GetUpToken() string {
	accessKey := "MgQ6P-51OifHnEg2zhE5QVkamOB0inbJO9j5O9Bs" //"po0wuUx6fQ8G9xwm_wtwREIvKTMfMIpql78qLSlI"
	secretKey := "jnJwXuTAlrBQlBP_d0NoJOG_UxdZruoKwuuYvX-T" //"9cfW_YG4N2iraQXpvhdSChqxfhHTTsE52pXVV1um"
	bucket := "misihu"
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	putPolicy.Expires = 1800 //示例30分钟有效期
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	return upToken
}

// UploadFileToQny 上传到七牛云对象存储
func UploadFileToQny(fileName, localFile string) string {
	key := "image/" + fileName
	upToken := GetUpToken()
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
	// 是否使用https域名
	cfg.UseHTTPS = true
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	// 构建表单上传的对象
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		fmt.Println(err)
		return ""
	} else {
		return ret.Key
	}
}

// SaveImageFile 保存上传的图片文件到本地
func SaveImageFile(c *gin.Context, pathFile string, file *multipart.FileHeader) (bool, string, string) {
	fileExt := strings.ToLower(path.Ext(file.Filename))
	if fileExt != ".png" && fileExt != ".jpg" && fileExt != ".gif" && fileExt != ".jpeg" {
		return false, "", ""
	} else {
		date := time.Now().Format("20060102")
		if ok := IsFileExist(pathFile + "/" + date); !ok {
			_ = os.Mkdir(pathFile+"/"+date, 0666)
		}
		var fileName = strconv.FormatInt(time.Now().Unix(), 10) + CreateRandomNumber(5) + ".txt"
		var filePath = pathFile + "/" + date + "/" + fileName
		err := c.SaveUploadedFile(file, filePath)
		if err != nil {
			return false, "", ""
		} else {
			return true, fileName, filePath
		}
	}
}

// SaveFile 保存上传的文件到本地
func SaveFile(c *gin.Context, pathFile string, file *multipart.FileHeader) (bool, string, string) {
	fileExt := strings.ToLower(path.Ext(file.Filename))
	date := time.Now().Format("20060102")
	if ok := IsFileExist(pathFile + "/" + date); !ok {
		_ = os.Mkdir(pathFile+"/"+date, 0666)
	}
	var fileName = strconv.FormatInt(time.Now().Unix(), 10) + CreateRandomNumber(5) + fileExt
	var filePath = pathFile + "/" + date + "/" + fileName
	err := c.SaveUploadedFile(file, filePath)
	if err != nil {
		return false, "", ""
	} else {
		return true, fileName, filePath
	}
}

// IsFileExist 判断文件是否存在
func IsFileExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
