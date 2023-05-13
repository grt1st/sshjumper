package conf

import (
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
)

var (
	// username:password, password encrypted by bcrypt
	jumperUsers = map[string]string{
		"foo": "243261243130246d3058636f30592e41722e70306e2e586c6137566975535066686b324d784e596273506e554b6f485252472e527535455843426d2e",
	}
	// username:public_key_path
	jumpKeys = map[string]string{
		"foo": "/Users/grt1st/.ssh/gogo.pub",
	}
	// ssh-machine info
	jumpMachines = map[string][]string{
		"127.0.0.1:2222": {
			"username", "password",
		},
	}
)

const (
	ServerAddr     = "127.0.0.1:2200"          // ssh jumper host
	PrivateKeyPath = "/Users/grt1st/.ssh/gogo" // ssh jumper private key
)

var authorizedKeysMap map[string]string

// ConnectSSHPassword authorization to ssh jumper
func ConnectSSHPassword(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error) {
	password, ok := jumperUsers[c.User()]
	if ok {
		comparePassword, _ := hex.DecodeString(password)
		result := bcrypt.CompareHashAndPassword(comparePassword, pass)
		if result == nil {
			return nil, nil
		}
	}
	return nil, fmt.Errorf("password rejected for %q", c.User())
}

func InitSSHPublicKey() {
	authorizedKeysMap = make(map[string]string)

	for u, publicKeyPath := range jumpKeys {
		keyBytes, err := ioutil.ReadFile(publicKeyPath)
		if err != nil {
			continue
		}
		pubKey, _, _, _, err := ssh.ParseAuthorizedKey(keyBytes)
		if err != nil {
			continue
		}
		authorizedKeysMap[string(pubKey.Marshal())] = u
	}
}

// ConnectSSHPublicKey authorization to ssh jumper
func ConnectSSHPublicKey(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
	if authorizedKeysMap[string(pubKey.Marshal())] == c.User() {
		return &ssh.Permissions{
			Extensions: map[string]string{
				"pubkey-fp": ssh.FingerprintSHA256(pubKey),
			},
		}, nil
	}
	return nil, fmt.Errorf("unknown public key for %q", c.User())
}

// GetRemoteSSH authorization to ssh jumper slave
func GetRemoteSSH(host string, defaultChoice bool) (string, *ssh.ClientConfig, error) {
	var hostInfo []string
	if defaultChoice {
		if len(jumpMachines) == 0 {
			return "", nil, errors.New("there are no machines to use")
		}
		for k, v := range jumpMachines {
			host = k
			hostInfo = v
			break
		}
	} else {
		var ok bool
		hostInfo, ok = jumpMachines[host]
		if !ok {
			return "", nil, errors.New("machine not found")
		}
	}
	if len(hostInfo) < 2 {
		return "", nil, errors.New("host info not correct")
	}
	// do something...
	sshConfig := &ssh.ClientConfig{
		User: hostInfo[0],
		Auth: []ssh.AuthMethod{
			ssh.Password(hostInfo[1]), // todo: 支持证书
		},
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		}, // must have
	}
	return host, sshConfig, nil
}
