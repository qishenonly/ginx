package giface

// 定义一个服务器接口
type IServer interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//运行服务器
	Server()

	//路由功能：给当前的服务注册一个路由方法，供客户端的链接处理使用
	AddRouter(msgID uint32, router IRouter)
	//获取当前server的链接管理器
	GetConnManager() IConnManager
	//注册AfterConnCreate 钩子函数的方法
	SetAfterConnCreate(func(connection IConnection))
	//注册BeforeConnDestory 钩子函数的方法
	SetBeforeConnDestory(func(connection IConnection))
	//调用AfterConnCreate 钩子函数的方法
	CallAfterConnCreate(connection IConnection)
	//注册BeforeConnDestory 钩子函数的方法
	CallBeforeConnDestory(connection IConnection)
}
