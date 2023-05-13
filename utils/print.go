package utils

import (
	"fmt"
	"strings"
	"time"
)

const (
	CRLF              = "\r\n"
	ClearPreviousLine = "\033[1A\033[K"
	Delimiter         = "==============================================================="
	WordsWelcome      = Delimiter + CRLF +
		"　  へ　　　　　／|" + CRLF +
		"　　/＼7　　　 ∠＿/ " + CRLF +
		"　 /　│　　 ／　／" + CRLF +
		"　│　Z ＿,＜　／　　 /� 　│　　　　　�　　 /　　〉 　 Y　　　　　　 /　　/" + CRLF +
		"　●　　●　　〈　　/" + CRLF +
		"　()  へ　　　　|　＼〈" + CRLF +
		"　　> _　 ィ　 │ ／／" + CRLF +
		"　 / へ　　 /　＜| ＼＼" + CRLF +
		"　 �_　　(_／　 │／／" + CRLF +
		"　　7　　　　　　　|／" + CRLF +
		"　　＞�r￣￣r�＿" + CRLF + Delimiter + CRLF
	WordsInfo            = "\x1b[33;40m Welcome %s. %s \x1b[0m" + CRLF
	WordsDone            = "\x1b[33;40m Goodbye. Good luck. \x1b[0m" + CRLF
	WordsNotFound        = "Command not found: %s" + CRLF
	WordsCommandError    = "Command exec wrong: %s" + CRLF
	WordsNotSupport      = "Sorry, command not support now: %s" + CRLF
	WordsSSHNotAvailable = "SSH not available: %s" + CRLF
	WordsSSHInfo         = "Remote addr is %s. " + CRLF + Delimiter + CRLF
	WordsSSHLoading      = "\x1b[33;40m The connection is being established, please wait. \x1b[0m" + CRLF
	WordsSSHDone         = CRLF + "\x1b[33;40m Connection closed. Please press Enter twice to continue. \x1b[0m" + CRLF
	WordsExecWrong       = "cmd.Run() failed with %s" + CRLF
)

type LogStr string

func (s LogStr) Bytes() []byte {
	return []byte(s)
}

func (s LogStr) S() string {
	return string(s)
}

// LogSprintf format with time
func LogSprintf(format string, a ...interface{}) LogStr {
	return LogStr(fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(format, a...)))
}

// GetHelpDoc get help doc
func GetHelpDoc(commands map[string]string, crlf string) string {
	var maxLength int
	var commandsDoc string
	for k := range commands {
		if len(k) > maxLength {
			maxLength = len(k)
		}
	}
	for k, v := range commands {
		commandsDoc += strings.Repeat(" ", 4) + k + strings.Repeat(" ", 4+maxLength-len(k)) + strings.Title(v) + crlf
	}
	return "Usage: <command> [args]" + crlf + crlf +
		"Commands:" + crlf + commandsDoc
}
