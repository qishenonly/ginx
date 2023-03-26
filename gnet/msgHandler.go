package gnet

import (
	"fmt"
	"ginx/giface"
	"ginx/utils"
	"strconv"
)

/*
消息处理模块的实现
*/
type MsgHandle struct {
	Apis map[uint32]giface.IRouter
	//负责Worker取任务的消息队列
	TaskQueue []chan giface.IRequest
	//业务工作Worker池的worker数量
	WorkerPoolSize uint32
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]giface.IRouter),
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize, //从全局配置中获取
		TaskQueue:      make([]chan giface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}

// 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request giface.IRequest) {
	//1 从Request中找到msgID
	handle, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is NOT FOUND! Need Register!")
	}
	//2 根据MsgID 调度对应的router业务
	handle.PreHandle(request)
	handle.Handle(request)
	handle.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router giface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		//id已经注册了
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	//2 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add API msgID = ", msgID, ", success!")
}

// 启动一个Worker工作池（开启工作池的动作只能发生一次，一个ginx框架只能有一个工作池）
func (mh *MsgHandle) StartWorkerPool() {
	//根据workerPoolSize 分别开启Worker， 每个Worker用一个go来承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//1 当前的worker对应的channel消息队列 开辟空间--第0个worker对应第0个channel...
		mh.TaskQueue[i] = make(chan giface.IRequest, utils.GlobalObject.MaxWorkerTaskSize)
		//2 启动当前的worker，阻塞等待消息从channel传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个worker工作流程
func (mh *MsgHandle) StartOneWorker(worerID int, taskQueue chan giface.IRequest) {
	fmt.Println("Worker ID = ", worerID, " is started ...")

	//不断的阻塞等待对应消息队列的消息
	for {
		select {
		//如果有消息过来，出列的就是一个客户端的request，执行当前request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 将消息交给TaskQueue， 由worker进行处理
func (mh *MsgHandle) SendMsgTOTaskQueue(request giface.IRequest) {
	//1 将消息平均分配给不同的worker
	//根据客户端建立的ConnID来进行分配
	//TODO 不是正规的负载均衡算法，后续可以对每一个request加requestID来替换ConnID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID = ", request.GetConnection().GetConnID(),
		" request MsgID = ", request.GetMsgID(),
		" to WorkerID = ", workerID)

	//2 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- request
}
