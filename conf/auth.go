package conf

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/grt1st/sshjumper/utils"
)

var (
	sshUsername = "foo"
	sshPassword = "bar"
	username    = "username"
	password    = "password"
	host        = "remote"
)

const (
	PrivateKeyPath = "key_path"
)

// ConnectSSHPassword 连接当前ssh
func ConnectSSHPassword(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	if c.User() == sshUsername && string(pass) == sshPassword {
		return nil, nil
	}
	return nil, fmt.Errorf("password rejected for %q", c.User())
}

// ConnectSSHPublicKey 连接当前ssh
func ConnectSSHPublicKey(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
	authorizedKeysMap := map[string]bool{}
	if authorizedKeysMap[string(pubKey.Marshal())] {
		return &ssh.Permissions{
			Extensions: map[string]string{
				"pubkey-fp": ssh.FingerprintSHA256(pubKey),
			},
		}, nil
	}
	return nil, fmt.Errorf("unknown public key for %q", c.User())
}

// GetRemoteSSH 获取远端ssh信息
func GetRemoteSSH(command utils.Command, serverConn *ssh.ServerConn) (string, *ssh.ClientConfig, error) {
	// do something...
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password), // todo: 支持证书
		},
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		}, // must have
	}
	return host, sshConfig, nil
}
