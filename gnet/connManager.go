package gnet

import (
	"errors"
	"fmt"
	"ginx/giface"
	"sync"
)

/*
	链接管理模块
*/

type ConnManager struct {
	connections map[uint32]giface.IConnection //管理的链接集合
	connLock    sync.RWMutex                  //保护链接集合的读写锁
}

// 创建当前链接的方法
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]giface.IConnection),
	}
}

func (cm *ConnManager) Add(conn giface.IConnection) {
	//保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//将conn加入到ConnManager中
	cm.connections[conn.GetConnID()] = conn

	fmt.Println("connID = ", conn.GetConnID(), ", connection Add to ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Remove(conn giface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//删除链接信息
	delete(cm.connections, conn.GetConnID())

	fmt.Println("connID = ", conn.GetConnID(), ", connection Remove from ConnManager successfully: conn num = ", cm.Len())
}

func (cm *ConnManager) Get(connID uint32) (giface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()

	if conn, ok := cm.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not FOUND!")
	}
}

func (cm *ConnManager) Len() int {
	return len(cm.connections)
}

func (cm *ConnManager) ClearConn() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	//删除conn并停止conn的工作
	for connID, conn := range cm.connections {
		//停止
		conn.Stop()
		//删除
		delete(cm.connections, connID)
	}

	fmt.Println("Clear All connections success! conn num = ", cm.Len())
}
