package giface

/*
链接管理模块抽象层
*/
type IConnManager interface {
	//添加链接
	Add(conn IConnection)
	//删除链接
	Remove(conn IConnection)
	//根据connID获取链接
	Get(connID uint32) (IConnection, error)
	//得到当前链接总数
	Len() int
	//清除并终止所有的链接
	ClearConn()
}
