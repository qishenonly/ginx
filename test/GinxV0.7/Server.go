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
type PingRouter struct {
	gnet.BaseRouter
}

// Test Handle
func (this *PingRouter) Handle(request giface.IRequest) {
	fmt.Println("Call PingRouter Handle...")
	//先读取客户端的数据，再回写
	fmt.Println("recv from client: MsgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(200, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

// Hello Ginx test 自定义路由
type HelloGinxRouter struct {
	gnet.BaseRouter
}

// Test Handle
func (this *HelloGinxRouter) Handle(request giface.IRequest) {
	fmt.Println("Call HelloGinxRouter Handle...")
	//先读取客户端的数据，再回写
	fmt.Println("recv from client: MsgID = ", request.GetMsgID(),
		", data = ", string(request.GetData()))

	err := request.GetConnection().SendMsg(201, []byte("Hello Welcome To Ginx!!"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	//1. 创建一个server句柄，使用Ginx的api
	s := gnet.NewServer("[ginx V0.5.1]")

	//给当前ginx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloGinxRouter{})

	//3. 启动server
	s.Server()
}
