package gnet

import (
	"fmt"
	"ginx/giface"
	"ginx/utils"
	"net"
)

// IServer 的接口实现，定义一个Server的服务器模块
type Server struct {
	//服务器的名称
	Name string
	//服务器绑定的ip版本
	IPVersion string
	//服务器监听的ip
	IP string
	//服务器监听的端口
	Port int
	//当前server的消息管理模块， 用来绑定MsgID和对应的处理业务API
	MsgHandler giface.IMsgHandle
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Ginx] Server Name : %s, Listener at IP : %s, Port : %d is starting\n",
		utils.GlobalObject.Name, utils.GlobalObject.Host, utils.GlobalObject.TcpPort)

	go func() {
		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error :>>", err)
			return
		}

		//2 监听服务器的地址
		litener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen ", s.IPVersion, " err ", err)
			return
		}

		fmt.Println("start Ginx server success,", s.Name, "success listening...")
		var cid uint32
		cid = 0

		//阻塞的等待客户端的链接，处理客户端链接业务(读写)
		for {
			//如果有客户端过来，阻塞会返回
			conn, err := litener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:>>", err)
				continue
			}
			//将处理新链接的业务方法和conn进行绑定，得到链接模块
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++

			//启动当前的链接业务处理
			go dealConn.Start()
		}
	}()
}

// 停止服务器
func (s *Server) Stop() {
	//TODO 将一些服务器的资源，状态或者一些已经开辟的链接信息 进行停止或者回收
}

// 运行服务器
func (s *Server) Server() {
	//启动server的服务功能
	s.Start()

	//TODO 做一些启动服务器之后的额外业务

	//阻塞状态
	select {}
}

// 路由功能：给当前的服务注册一个路由方法
func (s *Server) AddRouter(msgID uint32, router giface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("Add Router Success!!")
}

/*
初始化Server模块的方法
*/
func NewServer(name string) giface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
