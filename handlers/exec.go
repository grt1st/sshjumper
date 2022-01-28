package handlers

import (
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/grt1st/sshjumper/utils"
)

// HandleExec 处理 exec 类型的请求，如 scp
func HandleExec(command string, channel ssh.Channel) {
	_, _ = channel.Write(utils.LogSprintf(utils.WordsNotSupport, strings.Split(command, " ")[0]).Bytes())
}
