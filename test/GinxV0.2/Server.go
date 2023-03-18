package main

import "ginx/gnet"

/*
	基于Ginx框架来开发的 服务器端应用程序
*/

func main() {
	//1. 创建一个server句柄，使用Ginx的api
	s := gnet.NewServer("[ginx V0.2]")
	//2. 启动server
	s.Server()
}
