package logs

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

type ErrorLog struct{ msg string }

// Info
func (err *ErrorLog) InfoMsg(msgs ...string) {
	color.Green(
		"[INFO] %s",
		choiceValue(err.msg, msgs),
	)
}

// WarnMsg
func (err *ErrorLog) WarnMsg(msgs ...string) {
	color.Yellow(
		"[WARN] %s",
		choiceValue(err.msg, msgs),
	)
}

// DebugMsg
func (err *ErrorLog) DebugMsg(msgs ...string) {
	color.Blue(
		"[DEBUG] %s",
		choiceValue(err.msg, msgs),
	)
}

// ErrorMsg
func (err *ErrorLog) ErrorMsg(msgs ...string) {
	color.Red(
		"[Error] %s",
		choiceValue(err.msg, msgs),
	)
}

// choiceValue
func choiceValue(a, b any) (tmp string) {
	if a.(string) == "" {
		msgs := b.([]string)
		switch msgLen := len(msgs); msgLen {
		case 0:
			panic("Must have a value to msg.")
		case 1:
			tmp = msgs[0]
		default:
			tmp = strings.Join(msgs, " - ")
		}
		return tmp
	}
	return a.(string)
}

func NewDefaultLog() *ErrorLog {
	return &ErrorLog{}
}

// CustomText 自定义的日志格式
type CustomText struct{}

func (f *CustomText) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := entry.Level.String()
	callerInfo := entry.Message

	// 获取当前文件名和行号
	_, file, line, ok := runtime.Caller(1)
	if ok {
		// 提取文件名并添加到日志消息
		fileName := path.Base(file)
		callerInfo = fmt.Sprintf("%s:%d | %s", fileName, line, callerInfo)
	}

	// 手动设置颜色
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = 36 // Cyan
	case logrus.InfoLevel:
		levelColor = 32 // Green
	case logrus.WarnLevel:
		levelColor = 33 // Yellow
	case logrus.ErrorLevel:
		levelColor = 31 // Red
	default:
		levelColor = 0 // Default color
	}

	return []byte(fmt.Sprintf("\x1b[1;%dm%s | %s | %s\x1b[0m\n", levelColor, timestamp, level, callerInfo)), nil
}

func InitLogger() *logrus.Logger {
	logger := logrus.New()

	// 设置日志输出格式为自定义格式
	logger.SetFormatter(&CustomText{})

	// 设置日志级别
	logger.SetLevel(logrus.DebugLevel)

	return logger
}
