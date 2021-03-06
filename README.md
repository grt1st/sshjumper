# SSH Jumper

SSH Jumper 是一个提供SSH跳转功能的服务，具备身份认证、操作记录功能。

SSH Jumper is a jump server for SSH, which has function of authentication and command logged. 

> [中文版本 README](https://github.com/grt1st/sshjumper/blob/master/README.md)
> 
> [English Version README](https://github.com/grt1st/sshjumper/blob/master/README_EN.md)

## 快速开始

使用 `docker` 启动一个 ssh 服务器作为我们需要跳转的机器：
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

如果 docker 启动成功，我们可以这样连接到实例上：
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

之后我们在本地生成公私钥，用于 ssh 的认证：
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

之后我们下载仓库到本地，修改文件：
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

现在启动服务：
```commandline
> go run server.go
```

新开一个命令行，通过 ssh 连接到服务上（密码为bar）：
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

## 使用文档

### 配置

配置项目集中在 `conf/auth.go` 中：

```go
var (
    sshUsername = "foo"      // 连接到 ssh jumper 的用户名
    sshPassword = "bar"      // 连接到 ssh jumper 的密码
    username    = "username" // 通过 ssh jumper 中转的 ssh 账号用户名
    password    = "password" // 通过 ssh jumper 中转的 ssh 账号密码
    host        = "host"     // 通过 ssh jumper 中转的远端 ssh 的地址
)

const (
    ServerAddr     = "127.0.0.1:2200"   // ssh jumper 的监听地址
    PrivateKeyPath = "private_key_path" // ssh jumper 的 ssh private key 
)

// 连接到 ssh jumper 的账号认证方法
func ConnectSSHPassword(c ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error)

// 连接到 ssh jumper 的公钥认证方法
func ConnectSSHPublicKey(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error)

// 获取连接远端 ssh 的信息，如果没权限返回错误
func GetRemoteSSH(command utils.Command, serverConn *ssh.ServerConn) (string, *ssh.ClientConfig, error)

```

### 命令

连接到 ssh jumper 后，支持如下命令

```commandline
> help
Usage: <command> [args]

Commands:
    ssh     Ssh To Remote Host.
    exec    Execute Command.
    exit    Logout
```

## 变更记录

- 2022.01.29 初始版本。

