package utils

import (
	"encoding/json"
	"ginx/giface"
	"os"
)

/*
	存储一些有关ginx框架的全局参数
	一些参数是可以通过ginx.json由用户进行配置
*/

type GlobalObj struct {
	/*
		Server
	*/
	TcpServer giface.IServer //当前Ginx全局的Server对象
	Host      string         //当前服务器主机监听的IP
	TcpPort   int            //当前服务器主机监听的端口号
	Name      string         //当前服务器的名称

	/*
		Ginx
	*/
	Version           string //当前Ginx的版本号
	MaxConn           int    //当前服务器主机允许的最大链接数
	MaxPackageSize    uint32 //当前Ginx框架数据包的最大值
	WorkerPoolSize    uint32 //当前业务工作Worker池的Goroutine数量
	MaxWorkerTaskSize uint32 //Ginx框架允许用户最多开辟多少个Worker
}

/*
定义一个全局的对外GlobalObj
*/
var GlobalObject *GlobalObj

/*
初始化GlobalObject
*/
func init() {
	//如果配置文件没有加载，默认值
	GlobalObject = &GlobalObj{
		Name:              "GinxServerApp",
		Version:           "V0.9",
		TcpPort:           8999,
		Host:              "0.0.0.0",
		MaxConn:           1000,
		MaxPackageSize:    4096,
		WorkerPoolSize:    10,   //Worker工作池的队列的个数
		MaxWorkerTaskSize: 1024, //每个Worker对应的消息队列的任务的数量最大值
	}

	//从conf/ginx.json加载用户自定义参数
	GlobalObject.Reload()
}

// 加载ginx.json中的参数
func (g *GlobalObj) Reload() {
	data, err := os.ReadFile("conf/ginx.json")
	if err != nil {
		panic(err)
	}
	//将json文件数据解析到struct中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}

}
