# SSH Jumper

SSH Jumper 是一个提供SSH跳转功能的服务，具备身份认证、操作记录功能。

SSH Jumper is a jump server for SSH, which has function of authentication and command logged. 

> [中文版本 README](https://github.com/grt1st/sshjumper/blob/master/README.md)
>
> [English Version README](https://github.com/grt1st/sshjumper/blob/master/README_EN.md)

## Quickstart

Use `docker` to create a ssh server which we will connect to it use `SSH Jumper` later.
```commandline
docker run -d \
  --name=sshjumer-slave \
  -e PUID=1000 \
  -e PGID=1000 \
  -e TZ=Europe/London \
  -e SUDO_ACCESS=true `#optional` \
  -e PASSWORD_ACCESS=true `#optional` \
  -e USER_PASSWORD=password \
  -e USER_NAME=username \
  -p 2222:2222 \
  -v /tmp/config:/config \
  --restart unless-stopped \
  lscr.io/linuxserver/openssh-server
```

Then we could connect to the docker instance by ssh if it's success.
```commandline
> ssh -p 2222 username@127.0.0.1
The authenticity of host '[127.0.0.1]:2222 ([127.0.0.1]:2222)' can't be established.
ECDSA key fingerprint is SHA256:dRznQpRa4YN11KJpYAFOAEMcSB7FP9PS0KLba8RZ5vk.
Are you sure you want to continue connecting (yes/no)? yes
Warning: Permanently added '[127.0.0.1]:2222' (ECDSA) to the list of known hosts.
username@127.0.0.1's password:
Welcome to OpenSSH Server

bf215bb0398a:~$ whoami
username
bf215bb0398a:~$ exit
logout
Connection to 127.0.0.1 closed.
```

Then we generate a key pair to use as ssh authorization：
```commandline
ssh-keygen -t dsa
Generating public/private dsa key pair.
Enter file in which to save the key (/Users/grt1st/.ssh/id_dsa):
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
Your identification has been saved in /Users/grt1st/.ssh/id_dsa.
Your public key has been saved in /Users/grt1st/.ssh/id_dsa.pub.
The key fingerprint is:
SHA256:p5ZIczS2vV+DGd3i98/MCguj6ryXydQRkJK/SDjp6qs kintenroku@bogon
The key's randomart image is:
+---[DSA 1024]----+
|       ..o       |
|      o . .      |
|     o o+  .     |
|    + .o.+.  . . |
|   . oo.Soo.. o .|
|    ...+o+.. = . |
|   .  .o+o+ + + .|
|  .  . .=. + + *.|
|Eoo. .=+.   o ..B|
+----[SHA256]-----+
```

Then we download the code repository, and edit file.
```commandline
> vim conf/auth.go

var (
    .......
	host        = "127.0.0.1:2222"
)

const (
	......
	PrivateKeyPath = "/Users/grt1st/.ssh/id_dsa"
)
```

Now let's start the service.
```commandline
> go run server.go
```

Create a new terminal, and connect to ssh jumper by ssh. Password is "bar".
```commandline
～ ssh foo@127.0.0.1 -p 2200
foo@127.0.0.1's password:
===============================================================
　  へ　　　　　／|
　　/＼7　　　 ∠＿/
　 /　│　　 ／　／
　│　Z ＿,＜　／　　 /� 　│　　　　　�　　 /　　〉 　 Y　　　　　　 /　　/
　●　　●　　〈　　/
　()  へ　　　　|　＼〈
　　> _　 ィ　 │ ／／
　 / へ　　 /　＜| ＼＼
　 �_　　(_／　 │／／
　　7　　　　　　　|／
　　＞�r￣￣r�＿
===============================================================
2022-01-29 11:04:41  Welcome foo. 我们经常在正确的事情和容易的事情之间做选择.
> ssh
2022-01-29 11:04:43 Remote addr is 127.0.0.1:2222.
===============================================================
Welcome to OpenSSH Server

bf215bb0398a:~$ whoami
username
bf215bb0398a:~$ exit
logout
2022-01-29 11:04:47
 Connection closed. Please press Enter twice to continue.
>
>
> exit
 Goodbye. Good luck.
Connection to 127.0.0.1 closed.
```

## Usage

### Config

The configuration is in `conf/auth.go`.

```go
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

// authorization to ssh jumper
func ConnectSSHPassword(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error)

// authorization to ssh jumper
func ConnectSSHPublicKey(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error)

// authorization to ssh jumper slave
func GetRemoteSSH(command utils.Command, serverConn *ssh.ServerConn) (string, *ssh.ClientConfig, error)

```

### Command

The commands `SSH Jumper` support is in below:

```commandline
> help
Usage: <command> [args]

Commands:
    ssh     Ssh To Remote Host.
    exec    Execute Command.
    exit    Logout
```

## Changelog

- 2022.01.29 Initial Version.
