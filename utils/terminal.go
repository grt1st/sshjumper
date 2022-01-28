package utils

import (
	"io"
	"log"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// TermWrapper 封装 terminal 为一个 io.reader
type TermWrapper struct {
	term *terminal.Terminal
	c    ssh.Channel
}

func NewTermWrapper(channel ssh.Channel) *TermWrapper {
	return &TermWrapper{c: channel}
}

func (t *TermWrapper) Read(p []byte) (int, error) {
	if t.term == nil {
		t.term = terminal.NewTerminal(t.c, "")
	}
	data, err := t.term.ReadLine()
	// 处理 ctrl-c
	if err == io.EOF {
		p[0] = 3
		t.term = nil
		return 1, nil
	}
	// 转化为可写
	i := 0
	for k, v := range []byte(data) {
		p[k] = v
		i += 1
	}
	// 添加换行 \n
	if !(p[0] == 10 && p[1] == 0) {
		p[i] = 10
		i += 1
	}
	// 预备添加 /r
	//p[i] = 13
	//i += 1
	log.Println("terminal read", p[:i], i, err)
	return i, err
}

func (t *TermWrapper) Write(p []byte) (int, error) {
	return t.term.Write(p)
}

func (t *TermWrapper) Close() (err error) {
	if t.term != nil {
		// do something...
		_, _ = t.term.Write(LogSprintf(WordsSSHDone).Bytes())
	}
	return err
}
