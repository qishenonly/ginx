package main

import (
	"fmt"
	"ginx/giface"
	"ginx/gnet"
)

/*
	基于Ginx框架来开发的 服务器端应用程序
*/

//ping test 自定义路由
type PingROuter struct {
	gnet.BaseRouter
}

//Test PreHandle
func (this *PingROuter) PreHandle(request giface.IRequest) {
	fmt.Println("Call Router PreHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping..."))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

//Test Handle
func (this *PingROuter) Handle(request giface.IRequest) {
	fmt.Println("Call Router Handle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping..."))
	if err != nil {
		fmt.Println("call back ping...ping...ping error")
	}
}

//Test PostHandle
func (this *PingROuter) PostHandle(request giface.IRequest) {
	fmt.Println("Call Router PostHandle...")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("after ping..."))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	//1. 创建一个server句柄，使用Ginx的api
	s := gnet.NewServer("[ginx V0.3]")

	//给当前ginx框架添加一个自定义的router
	s.AddRouter(&PingROuter{})

	//3. 启动server
	s.Server()
}
