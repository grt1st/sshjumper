package handlers

import (
	"os/exec"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/grt1st/sshjumper/conf"
	"github.com/grt1st/sshjumper/utils"
)

// HandlePty 处理 shell 类型的请求，处理为交互式命令行
func HandlePty(channel ssh.Channel, serverConn *ssh.ServerConn) {
	_, _ = channel.Write([]byte(utils.WordsWelcome))
	_, _ = channel.Write(utils.LogSprintf(utils.WordsInfo, serverConn.User(), "我们经常在正确的事情和容易的事情之间做选择.").Bytes())
	term := terminal.NewTerminal(channel, "> ")

	defer channel.Close()
	for {
		// get user input
		line, err := term.ReadLine()
		if err != nil {
			break
		}
		// parse input
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		command := utils.NewCommand(line)
		switch command.Name {
		case "ssh":
			host, sshConfig, err := conf.GetRemoteSSH(command, serverConn)
			if err != nil {
				_, _ = term.Write(utils.LogSprintf(utils.WordsSSHNotAvailable, err).Bytes())
				continue
			}
			_, _ = term.Write(utils.LogSprintf(utils.WordsSSHInfo, host).Bytes())
			_, _ = term.Write(utils.LogSprintf(utils.WordsSSHLoading).Bytes())
			wrapper := utils.NewTermWrapper(channel)
			session, _ := utils.CreateSSHSession(wrapper, host, sshConfig)
			_, _ = term.Write([]byte(utils.ClearPreviousLine))
			err = session.Wait()
			_ = session.Close()
			_ = wrapper.Close()
		case "exec":
			cmd, err := exec.Command("bash", "-c", strings.TrimPrefix(line, "exec")).Output()
			if err != nil {
				_, _ = term.Write(utils.LogSprintf(utils.WordsExecWrong, err).Bytes())
			}
			_, _ = term.Write(cmd)
		case "exit":
			_, _ = term.Write([]byte(utils.WordsDone))
			return
		case "help":
			_, _ = term.Write([]byte(utils.GetHelpDoc(map[string]string{
				"ssh": "ssh to remote host.",
				"exec": "execute command.",
				"exit": "logout",
			}, utils.CRLF)))
		default:
			_, _ = term.Write(utils.LogSprintf(utils.WordsNotFound, line).Bytes())
		}
	}
}
