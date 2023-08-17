package logs

import (
	"strings"

	"github.com/fatih/color"
)

type SpiderError struct {
	msg string
}

var BaseLog = SpiderError{}

// Info
func (err *SpiderError) InfoMsg(msgs ...string) {
	color.Green(
		"[INFO] %s",
		choiceValue(err.msg, msgs),
	)
}

// WarnMsg
func (err *SpiderError) WarnMsg(msgs ...string) {
	color.Yellow(
		"[WARN] %s",
		choiceValue(err.msg, msgs),
	)
}

// DebugMsg
func (err *SpiderError) DebugMsg(msgs ...string) {
	color.Blue(
		"[DEBUG] %s",
		choiceValue(err.msg, msgs),
	)
}

// ErrorMsg
func (err *SpiderError) ErrorMsg(msgs ...string) {
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
