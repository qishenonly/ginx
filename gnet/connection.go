package gnet

import (
	"errors"
	"fmt"
	"ginx/giface"
	"ginx/utils"
	"io"
	"net"
	"sync"
)

/*
链接模块
*/
type connection struct {
	// 当前Conn隶属于哪个server
	TCPServer giface.IServer
	// 当前链接的 socket TCP套接字
	Conn *net.TCPConn

	// 链接的ID
	ConnID uint32

	//当前的链接状态
	isClosed bool

	//告知当前链接已经退出的/停止 channel -- 由Reader告知Writer退出
	EixtChan chan bool

	//无缓冲的管道，用于读写 Groutine之间的消息通信
	msgChan chan []byte

	//消息管理MsgID和对应的处理业务的api
	MsgHandler giface.IMsgHandle

	//链接属性的集合
	property map[string]interface{}

	//保护链接属性的锁
	propertyLock sync.RWMutex
}

// 初始化链接模块的方法
func NewConnection(server giface.IServer, conn *net.TCPConn, connId uint32, msgHandle giface.IMsgHandle) *connection {
	c := &connection{
		TCPServer:  server,
		Conn:       conn,
		ConnID:     connId,
		MsgHandler: msgHandle,
		isClosed:   false,
		EixtChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		property:   make(map[string]interface{}),
	}

	//将conn加入到ConnManager中
	c.TCPServer.GetConnManager().Add(c)

	return c
}

// 链接的读业务方法
func (c *connection) StartReader() {
	fmt.Println("[Reader Goroutine is running]")
	defer fmt.Println("connID = ", c.ConnID, "[Reader is exit]")
	defer c.Stop()

	for {
		//读取客户端的数据到buf中
		//buf := make([]byte, utils.GlobalObject.MaxPackageSize)
		//_, err := c.Conn.Read(buf)
		//if err != nil {
		//	fmt.Println("receive buf err:>>", err)
		//	continue
		//}

		//创建一个拆包解包对象
		dp := NewDataPack()

		//读取客户端的Mdg Head 二进制流的 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err:>>", err)
			break
		}

		//拆包，得到msgID  和 msgDatelen 放在msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack err:>>", err)
			break
		}

		//根据Datalen 再次读取Data 放在msg.Data中
		var data []byte
		if msg.GetMsgLen() > 0 {
			data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data err:>>", err)
				break
			}
		}
		msg.SetData(data)

		//得到当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.WorkerPoolSize > 0 {
			//已经开启了工作池机制，将消息发送给worker工作池处理即可
			c.MsgHandler.SendMsgTOTaskQueue(&req)
		} else {
			//从路由中找到注册绑定的Conn对应的router调用
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 写消息GOroutine，向客户端发消息的模块
func (c *connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit]")

	//不断的阻塞等待channel的消息
	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error:>>", err)
				return
			}
		case <-c.EixtChan:
			//代表Reader已经退出，此时Writer也要退出
			return
		}

	}
}

// 启动链接：让当前的链接准备开始工作
func (c *connection) Start() {
	fmt.Println("Conn start()... ConnId = ", c.ConnID)
	//启动从当前链接的读数据业务
	go c.StartReader()
	//启动从当前链接写数据的业务
	go c.StartWriter()
	//调用创建链接之后的函数
	c.TCPServer.CallAfterConnCreate(c)
}

// 停止链接：结束当前链接的工作
func (c *connection) Stop() {
	fmt.Println("Conn stop()... Conn = ", c.ConnID)

	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}

	c.isClosed = true

	//调用销毁链接之前的函数
	c.TCPServer.CallBeforeConnDestory(c)

	//关闭socket链接
	c.Conn.Close()

	//告知writer关闭
	c.EixtChan <- true

	//将当前链接从ConnManager中删除
	c.TCPServer.GetConnManager().Remove(c)

	close(c.EixtChan)
	close(c.msgChan)
}

// 获取当前链接的绑定的socket conn
func (c *connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// 获取当前链接模块的链接ID
func (c *connection) GetConnID() uint32 {
	return c.ConnID
}

// 获取远程客户端的 TCP状态 IP Port
func (c *connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 提供一个SendMsg方法 将我们要发送给客户端的数据，先进行封包，再发送
func (c *connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connectino closed when send msg")
	}

	//将data进行封包 MsgDatalen|MsgID|data
	dp := NewDataPack()

	binaryMsg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack err msg id = ", msgId)
		return errors.New("Pack error msg")
	}

	c.msgChan <- binaryMsg

	return nil
}

// 设置链接属性
func (c *connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

// 获取链接属性
func (c *connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("No Property FOUND!")
	}
}

// 移除链接属性
func (c *connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
