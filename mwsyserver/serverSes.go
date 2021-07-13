package mwsyserver

import (
	//"encoding/json"
	"fmt"
	"sync"
	"mwsy/mwsylog"
	"github.com/gorilla/websocket"
)

type MwsyServerSes struct {
	//websocket真实链接
	m_WebConn *websocket.Conn
	//等待退出
	m_WaitExit sync.WaitGroup
	//会话对应的id
	m_nId int64
	//服务器指针：只有作为服务器时候不为空
	m_Server *MwsyServer
}

func (this *MwsyServerSes) GetWebCon() *websocket.Conn {
	return this.m_WebConn
}

func (this *MwsyServerSes) GetId() int64 {
	return this.m_nId
}

func (this *MwsyServerSes) SetId(id int64) {
	this.m_nId = id
}

func (this *MwsyServerSes) Close() {
}

func (this *MwsyServerSes) Send(msg interface{}) {
	this.m_Server.Send(this, msg)
}

func (this *MwsyServerSes) RecvLoop() {
	for {
		//读取消息，将消息放入管道中
		msgtype, msg, err := this.m_WebConn.ReadMessage()
		if err != nil {
			break
		}

		if msgtype == websocket.BinaryMessage {
			this.m_Server.Recv(this, msg)
		}
	}
	this.m_WaitExit.Done()
}

func (this *MwsyServerSes) Start() {
	this.m_WaitExit.Add(1)
	//添加一个会话
	this.m_Server.AddSes(this)
	//添加日志打印服务器接收到链接
	logstr := fmt.Sprintf("accept ses suc sesid = %v nowcount = %v",this.m_nId,this.m_Server.m_SesManerger.getCount())
	mwsylog.LogInst.MwsyLogNormalPrint(logstr)
	//开启接收循坏
	go this.RecvLoop()
	//等待接收循环
	go this.Wait()
}

func (this *MwsyServerSes) Wait() {
	this.m_WaitExit.Wait()
	if this.m_Server != nil {
		this.m_Server.DelSes(this)
	}
	fmt.Println("ses closed id == ", this.m_nId)
}

func NewMwsyServerSession(server *MwsyServer, webcon *websocket.Conn) *MwsyServerSes {
	this := &MwsyServerSes{}
	this.m_WebConn = webcon
	this.m_Server = server
	return this
}
