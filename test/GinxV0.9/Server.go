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

// 创建链接之后执行钩子函数
func DoCollectionBegin(conn giface.IConnection) {
	fmt.Println("====> DoCollectionBegin is Called ...")
	if err := conn.SendMsg(202, []byte("DoCollection Begin!")); err != nil {
		fmt.Println("DoCollectionBegin err:>>", err)
	}

	//给当前的链接设置一些属性
	fmt.Println("Set Conn Property ...")
	conn.SetProperty("Name", "qishenonly")
	conn.SetProperty("Github", "https://github.com/qishenonly")

}

// 链接断开之前执行钩子函数
func DoCollentionLost(conn giface.IConnection) {
	fmt.Println("====> DoCollentionLost is Called ...")
	fmt.Println("connID = ", conn.GetConnID(), " is Lost ...")

	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Name ==> ", name)
	}

	if home, err := conn.GetProperty("Github"); err == nil {
		fmt.Println("Github ==> ", home)
	}
}

func main() {
	//1. 创建一个server句柄，使用Ginx的api
	s := gnet.NewServer("[ginx V0.9]")

	//2. 注册链接钩子函数
	s.SetAfterConnCreate(DoCollectionBegin)
	s.SetBeforeConnDestory(DoCollentionLost)

	//3. 给当前ginx框架添加自定义的router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloGinxRouter{})

	//4. 启动server
	s.Server()
}
