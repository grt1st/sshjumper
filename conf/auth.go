package conf

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/grt1st/sshjumper/utils"
)

var (
	sshUsername = "foo"      // ssh jumper username
	sshPassword = "bar"      // ssh jumper password
	username    = "username" // ssh slave username
	password    = "password" // ssh slave password
	host        = "host"     // ssh slave host
)

const (
	ServerAddr     = "127.0.0.1:2200"   // ssh jumper host
	PrivateKeyPath = "private_key_path" // ssh jumper private key
)

// ConnectSSHPassword authorization to ssh jumper
func ConnectSSHPassword(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	if c.User() == sshUsername && string(pass) == sshPassword {
		return nil, nil
	}
	return nil, fmt.Errorf("password rejected for %q", c.User())
}

// ConnectSSHPublicKey authorization to ssh jumper
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

// GetRemoteSSH authorization to ssh jumper slave
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
