package loghandler

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"user-server/config"
)

// log 的一些设置
const (
	// LogPath  系统日志文件地址
	LogPath = "./logs/system.log"
	// ConsoleType 类型, 支持 console, file
	ConsoleType = "file"
)

// LogInfoTemplate log 信息模板
type LogInfoTemplate struct {
	UniqueID int64 `json:"uniqueID"`
	Data     interface{}
}

// UniqueID log 的唯一 id
var uniqueID int64

// log 的通用设置
func init() {
	conf := fmt.Sprintf(`{"filename":"%s"}`, config.Conf.Log.Path)
	logs.SetLogger(config.Conf.Log.Mode, conf)
	// 开启文件行号显示
	logs.EnableFuncCallDepth(true)
	// 因为是自己封装的需要将包层级给标明,否则文件行号只会显示依赖包中的行号
	logs.SetLogFuncCallDepth(4)
	// 异步 chan 的大小为1k
	logs.Async(1e3)
}

// Logger 获取 log 结构体
type Logger struct {
	bl *logs.BeeLogger
	ID string
}

type logConfig struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
}

// NewLog 生成一个新的 log 对象
func NewLog(id string) Logger {
	logger := Logger{
		ID: id,
	}
	return logger
}

// Info 日值类型
func (log *Logger) Info(msg interface{}, extra ...string) {
	data := logFormat(msg)

	var category string
	if len(extra) > 0 {
		category = extra[0]
	} else {
		category = "Info"
	}
	logs.Info("[%s] [%s] %s", log.ID, category, data)
}

// Error 日值类型
func (log *Logger) Error(msg interface{}, extra ...string) {
	data := logFormat(msg)

	var category string
	if len(extra) > 0 {
		category = extra[0]
	} else {
		category = "Error"
	}
	logs.Error("[%s] [%s] %s", log.ID, category, data)
}

// SetUniqueID 生成 uniqueID
func SetUniqueID(id int64) {
	uniqueID = id
}

// GetUniqueID 统一获取 uid
func GetUniqueID() int64 {
	return uniqueID
}

// GetLogTemplate 获取 log 模板
func GetLogTemplate() LogInfoTemplate {
	var tmpl LogInfoTemplate
	tmpl.UniqueID = GetUniqueID()
	return tmpl
}

// Debug 方法
func Debug(f interface{}, v ...interface{}) {
	data := logFormat(f)
	var category string
	if len(v) > 0 {
		category = v[0].(string)
	} else {
		category = "Debug"
	}

	logs.Info("[%s] %s", category, data)
}

// Info 方法
func Info(f interface{}, v ...interface{}) {
	data := logFormat(f)
	var category string
	if len(v) > 0 {
		category = v[0].(string)
	} else {
		category = "Info"
	}

	logs.Info("[%s] %s", category, data)
}

// Warn 方法
func Warn(f interface{}, v ...interface{}) {
	data := logFormat(f)
	var category string
	if len(v) > 0 {
		category = v[0].(string)
	} else {
		category = "Warn"
	}

	logs.Info("[%s] %s", category, data)
}

// Error 方法
func Error(f interface{}, v ...interface{}) {
	data := logFormat(f)
	var category string
	if len(v) > 0 {
		category = v[0].(string)
	} else {
		category = "Error"
	}

	logs.Info("[%s] %s", category, data)
}

// Emergency 方法
func Emergency(f interface{}, v ...interface{}) {
	tmpl := GetLogTemplate()
	tmpl.Data = f
	strByte, _ := json.Marshal(tmpl)
	logs.Emergency(string(strByte), v...)
}

// Critical 方法
func Critical(f interface{}, v ...interface{}) {
	tmpl := GetLogTemplate()
	tmpl.Data = f
	strByte, _ := json.Marshal(tmpl)
	logs.Critical(string(strByte), v...)
}

// Alert 方法
func Alert(f interface{}, v ...interface{}) {
	tmpl := GetLogTemplate()
	tmpl.Data = f
	strByte, _ := json.Marshal(tmpl)
	logs.Alert(string(strByte), v...)
}

// Notice 方法
func Notice(f interface{}, v ...interface{}) {
	tmpl := GetLogTemplate()
	tmpl.Data = f
	strByte, _ := json.Marshal(tmpl)
	logs.Notice(string(strByte), v...)
}

func logFormat(f interface{}, v ...interface{}) string {
	str, ok := f.(string)
	if !ok {
		strByte, _ := json.Marshal(f)
		str = string(strByte)
	}
	return str
}
