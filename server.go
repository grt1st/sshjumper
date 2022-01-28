package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"

	"golang.org/x/crypto/ssh"

	"github.com/grt1st/sshjumper/conf"
	"github.com/grt1st/sshjumper/handlers"
)

func main() {
	Listen("localhost:2500")
}

func Listen(addr string) {
	// listen
	localListener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("net.Listen failed: %v\n", err)
	} else {
		log.Printf("net.Listen at %v\n", addr)
	}

	for {
		// get connect
		localConn, err := localListener.Accept()
		if err != nil {
			fmt.Printf("localListener.Accept failed: %v\n", err)
		}
		go func() {
			CreatePortalSession(localConn)
		}()
	}
}

// CreatePortalSession 创建命令监听
func CreatePortalSession(conn net.Conn) {

	// server config
	config := ssh.ServerConfig{
		PasswordCallback: conf.ConnectSSHPassword,
		PublicKeyCallback: conf.ConnectSSHPublicKey,
	}

	// private key
	keyBytes, err := ioutil.ReadFile(conf.PrivateKeyPath)
	if err != nil {
		panic(err)
	}
	key, err := ssh.ParsePrivateKey(keyBytes)
	if err != nil {
		panic(err)
	}
	config.AddHostKey(key)

	// 创建
	serverConn, chans, reqs, err := ssh.NewServerConn(conn, &config)
	if err != nil {
		log.Fatalf("create sever conn failed: %v\n", err)
	} else {
		log.Printf("create server conn from %s\n", serverConn.RemoteAddr())
	}
	// 消费
	go ssh.DiscardRequests(reqs)
	// 消费
	go func(chans <-chan ssh.NewChannel) {
		for newChannel := range chans {
			channel, requests, err := newChannel.Accept()
			if err != nil {
				log.Printf("newChannel.Accept failed: %v\n", err)
				continue // todo: handler error
			}
			go func(in <-chan *ssh.Request) {
				for req := range in {
					switch req.Type {
					case "shell":
						go handlers.HandlePty(channel, serverConn)
						_ = req.Reply(req.Type == "shell", req.Payload)
					case "exec":
						var msg execMsg
						if err = ssh.Unmarshal(req.Payload, &msg); err != nil {
							log.Printf("error parsing ssh execMsg: %s\n", err)
							_ = req.Reply(false, nil)
						} else {
							go func(msg execMsg, ch ssh.Channel) {
								// ch can be used as a ReadWriteCloser if there should be interactivity
								handlers.HandleExec(msg.Command, ch)
								ex := exitStatusMsg{
									Status: 0,
								}
								// return the status code
								if _, err := ch.SendRequest("exit-status", false, ssh.Marshal(&ex)); err != nil {
									log.Printf("unable to send status: %v", err)
								}
								_ = ch.Close()
							}(msg, channel)
							_ = req.Reply(true, nil) // tell the other end that we can run the request
						}
					//case "env":
					//	 ignore
					//case "pty-req":
					//	 todo: change window size
					default:
						log.Printf("req_type=%v, payload=%v", req.Type, string(req.Payload))
					}
				}
			}(requests)
		}
	}(chans)
}

type exitStatusMsg struct {
	Status uint32
}

// RFC 4254 Section 6.5.
type execMsg struct {
	Command string
}