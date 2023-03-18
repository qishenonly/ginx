package gnet

import (
	"fmt"
	"ginx/giface"
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
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at IP :%s, Port %d, is starting\n", s.IP, s.Port)

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

		//阻塞的等待客户端的链接，处理客户端链接业务(读写)
		for {
			//如果有客户端过来，阻塞会返回
			conn, err := litener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err:>>", err)
				continue
			}

			//已经与客户端建立链接，做一些业务
			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("receive buf err:>>", err)
						continue
					}

					fmt.Printf("receive client buf %s, cnt %d\n", buf, cnt)

					// 回显功能
					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("write back buf err:>>", err)
						continue
					}
				}
			}()
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

/*
初始化Server模块的方法
*/
func NewServer(name string) giface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
