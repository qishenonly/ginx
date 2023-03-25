package main

import (
	"fmt"
	"ginx/giface"
	"ginx/gnet"
)

/*
	基于Ginx框架来开发的 服务器端应用程序
*/

// ping test 自定义路由
type PingROuter struct {
	gnet.BaseRouter
}

// Test Handle
func (this *PingROuter) Handle(request giface.IRequest) {
	fmt.Println("Call Router Handle...")
	//先读取客户端的数据，再回写
	fmt.Println("recv from client: MsgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1. 创建一个server句柄，使用Ginx的api
	s := gnet.NewServer("[ginx V0.5.1]")

	//给当前ginx框架添加一个自定义的router
	s.AddRouter(&PingROuter{})

	//3. 启动server
	s.Server()
}
