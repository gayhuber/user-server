package tools

import (
	"crypto/md5"
	"encoding/json"
	// "errors"
	"fmt"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
	"user-server/config"
	logs "user-server/tools/loghandler"
)

// GenerageUniqueID 获取uniqueid
func GenerageUniqueID(src string) (uid, salt string) {
	uuid, _ := uuid.NewV4()
	salt = uuid.String()[:32]
	h := md5.New()
	h.Write([]byte(src + salt))
	uid = fmt.Sprintf("%x", h.Sum(nil))
	return
}

// MD5String 生成 md5 字符串
func MD5String(v []byte) string {
	h := md5.New()
	h.Write(v)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// RandInt 生成随机数
func RandInt(min, max int) int {
	rand.Seed(time.Now().Unix())
	base := max - min
	return rand.Intn(base) + min
}

// RegisteSoyoungHalf 注册新氧半账号
func RegisteSoyoungHalf(openID, platform, avatar, nickname string) (uid int, isTrue bool) {
	api := config.Conf.Params["soyoung_register"]

	reuqestURL := fmt.Sprintf("%sopen_id=%s&nickname=%s&avatar=%s&platform=%s", api, openID, nickname, avatar, platform)

	logs.Info(reuqestURL)

	resp, err := http.Get(reuqestURL)
	if CheckErr(err) {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if CheckErr(err) {
		return
	}
	logs.Info(string(body), "SOYOUNG_HALF_USER_RESPONSE")

	respBody := respSoyoungHalf{}
	json.Unmarshal(body, &respBody)

	uid = respBody.ResponseData.User.ResponseData.UID
	isTrue = true
	return
}

// 以下三条是注册新氧半账号时返回的 json 结构
type respSoyoungHalf struct {
	ErrorCode    int
	ErrorMsg     string
	ResponseData soyoungUserPlugin
}
type soyoungUserPlugin struct {
	User         soyoungUser `json:"user"`
	DebugOpenUID int         `json:"debug_open_uid"`
}
type soyoungUser struct {
	ErrorCode    int
	ErrorMsg     string
	ResponseData soyoungUserBody
}
type soyoungUserBody struct {
	UID int `json:"uid"`
	Ext []interface{}
}

// GetClient 发送 get 请求的工具
func GetClient(url string, headers ...map[string]string) (resp *http.Response, err error) {
	defer func() {
		if err := recover(); err != nil {
			logs.Error(err)
		}
	}()

	req, err := http.NewRequest("GET", url, strings.NewReader("from=user-server"))
	if CheckErr(err) {
		return
	}

	client := &http.Client{}
	resp, err = client.Do(req)
	return
}

type H map[string]string

func SelectType() {
	// mp := map[string]string{
	// 	"key1": "value1",
	// }

	// mp := errors.New("this is error")
	mp := "asdfasdf"
	var ttt interface{}
	ttt = mp

	x := reflect.TypeOf(mp).String()
	fmt.Println("type:", x)

	// 通过反射进行判断
	switch reflect.TypeOf(mp).String() {
	case "map[string]string":
		fmt.Println("map 类型")
	case "*errors.errorString":
		fmt.Println("error 类型")
	case "string":
		fmt.Println("字符串类型")
	}

	// 另一种隐式的判断
	switch t := ttt.(type) {
	case string:
		fmt.Println("字符串")
	case map[string]string:
		fmt.Println("map 类型")
	case error:
		fmt.Println("error 类型")
	default:
		_ = t
		fmt.Println("unknown")
	}
}

// CheckErr 打印错误
func CheckErr(err error) bool {
	if err != nil {
		logs.Error(err)
		return true
	}
	return false
}

// JSONEncode 数据转字符串 (来自一个 phper 的倔强)
func JSONEncode(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// JSONDecode 字符串转数据 (来自一个 phper 的倔强)
func JSONDecode(data string, v interface{}) error {
	by := []byte(data)
	return json.Unmarshal(by, v)
}

// Map2Struct 利用 json 包进行格式转换(效率低, 之后会换成 github.com/mitchellh/mapstructure 这个包)
func Map2Struct(mp map[string]interface{}, v interface{}) error {
	by, err := json.Marshal(mp)
	if err != nil {
		return nil
	}
	return json.Unmarshal(by, v)
}

// CalculateTTL 计算到达十天后凌晨4点的秒数
func CalculateTTL() int {
	now := time.Now()
	h := time.Date(now.Year(), now.Month(), now.Day(), int(4), int(0), int(0), int(0), time.Local)
	target, _ := time.ParseDuration("240h")
	new := h.Add(target)
	ttl := new.Sub(now).Seconds()
	return int(ttl)
}

// String2Int 转义
func String2Int(s string) (int, error) {
	return strconv.Atoi(s)
}

// Strintg2Int64 转义
func Strintg2Int64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

// Int64ToString 转义
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}
