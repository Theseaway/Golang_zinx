// Client
package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	choice     int
}

func NewClient(ServerIp_ string, ServerPort_ int) *Client {

	client := &Client{
		ServerIp:   ServerIp_,
		ServerPort: ServerPort_,
		choice:     999,
	}

	conn_, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp_, ServerPort_))
	if err != nil {
		fmt.Println("net.Dial failed")
		return nil
	}
	client.conn = conn_
	return client
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) Run() {
	for client.choice != 0 {
		for client.menu() != true {
		}

		switch client.choice {
		case 1:
			fmt.Println("Broadcast")
			break
		case 2:
			fmt.Println("private chat")
			break
		case 3:
			fmt.Println("rename")
			break
		case 0:
			return
		}
	}
}
func (client *Client) menu() bool {
	var choice_ int
	fmt.Println("1.广播模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.exit")
	fmt.Scanln(&choice_)

	if choice_ >= 0 && choice_ <= 3 {
		client.choice = choice_
		return true
	} else {
		fmt.Println("Wrong input,please input number 0-3")
		return false
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址")
	flag.IntVar(&serverPort, "port", 9555, "设置服务器端口")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("---------->Connection to Server down")
		return
	}
	fmt.Println("Connect to Server")
	go client.DealResponse()
	client.Run()
}
