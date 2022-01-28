package utils

import (
	"io"

	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func CreateSSHSession(rw io.ReadWriter, host string, config *ssh.ClientConfig) (*ssh.Session, error) {
	// dial remote ssh client
	sshClientConn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, errors.WithMessage(err, "ssh.Dial failed")
	}

	// create ssh session
	session, err := sshClientConn.NewSession()
	if err != nil {
		return nil, errors.WithMessage(err, "NewSession failed")
	}

	// set remote conn
	session.Stdout = rw
	session.Stderr = rw
	session.Stdin = rw

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
		ssh.IXANY: 1,
		ssh.IMAXBEL: 1,
	}
	// Request pseudo terminal
	if err = session.RequestPty("xterm", 40, 80, modes); err != nil {
		_ = session.Close()
		return nil, errors.WithMessage(err, "RequestPty failed")
	}
	// Start remote shell
	if err = session.Shell(); err != nil {
		_ = session.Close()
		return nil, errors.WithMessage(err, "start shell failed")
	}
	// Use session
	return session, errors.WithMessage(err, "session wait failed")
}
