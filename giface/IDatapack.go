package giface

/*
	封包、拆包 模块
	直接面向TCP链接中的数据流，用于处理TCP粘包问题
*/

type IDataPack interface {
	//获取包的头的长度 方法
	GetHeadLen() uint32
	//封包方法
	Pack(msg IMessage) ([]byte, error)
	//拆包方法
	Unpack([]byte) (IMessage, error)
}
