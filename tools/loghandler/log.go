package loghandler

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
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

// log 的基础设置都在这里了
func init() {
	config := fmt.Sprintf(`{"filename":"%s"}`, LogPath)
	logs.SetLogger(ConsoleType, config)
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
func NewLog(id string, path string, mode string) Logger {
	logger := Logger{
		ID: id,
	}
	config := fmt.Sprintf(`{"filename":"%s"}`, path)
	logger.bl = logs.NewLogger()
	logger.bl.SetLogger(mode, config)
	// 开启文件行号显示
	logger.bl.EnableFuncCallDepth(true)
	// 因为是自己封装的需要将包层级给标明,否则文件行号只会显示依赖包中的行号
	logger.bl.SetLogFuncCallDepth(4)
	// 异步 chan 的大小为1k
	logger.bl.Async(1e3)
	return logger
}

// Info 日值类型
func (log *Logger) Info(f interface{}, v ...interface{}) {
	tmpl := GetLogTemplate()
	tmpl.Data = f
	strByte, _ := json.Marshal(tmpl)

	log.bl.Info(string(strByte), v...)
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
	tmpl := GetLogTemplate()
	tmpl.Data = f
	strByte, _ := json.Marshal(tmpl)
	logs.Debug(string(strByte), v...)
}

// Info 方法
func Info(f interface{}, v ...interface{}) {
	tmpl := GetLogTemplate()
	tmpl.Data = f
	strByte, _ := json.Marshal(tmpl)

	logs.Info(string(strByte), v...)
}

// Warn 方法
func Warn(f interface{}, v ...interface{}) {
	tmpl := GetLogTemplate()
	tmpl.Data = f
	strByte, _ := json.Marshal(tmpl)
	logs.Warn(string(strByte), v...)
}

// Error 方法
func Error(f interface{}, v ...interface{}) {
	tmpl := GetLogTemplate()
	tmpl.Data = f
	strByte, _ := json.Marshal(tmpl)
	logs.Error(string(strByte), v...)
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
