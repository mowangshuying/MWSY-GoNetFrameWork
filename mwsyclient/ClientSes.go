package mwsyclient

import (
	//"fmt"
	"sync"
	"github.com/gorilla/websocket"
	"mwsy/mwsylog"
)

type  MwsyClientSes struct {
	//websocket真实链接
	m_WebConn *websocket.Conn
	//等待退出
	m_WaitExit sync.WaitGroup
	//客户端指针
	m_Client   *MwsyClient
}

func (this *MwsyClientSes) GetWebCon() *websocket.Conn {
	return this.m_WebConn
}

func (this *MwsyClientSes) Send(msg interface{}) {
	this.m_Client.Send(this, msg)
}

func (this *MwsyClientSes) RecvLoop() {
	for {
		//读取消息，将消息放入管道中
		msgtype, msg, err := this.m_WebConn.ReadMessage()
		if err != nil {
			//接收到错误消息，断开与服务器的链接
			mwsylog.LogInst.MwsyLogNormalPrint("ses close!")
			break
		}

		if msgtype == websocket.BinaryMessage {
			this.m_Client.Recv(this, msg)
		}
	}
	this.m_WaitExit.Done()
}

func (this *MwsyClientSes) Start() {
	this.m_WaitExit.Add(1)
	//开启接收循坏
	go this.RecvLoop()
	//等待接收循环
	go this.Wait()
}

func (this *MwsyClientSes) Wait() {
	this.m_WaitExit.Wait()
}

func NewMwsyClientSes(client *MwsyClient,webcon *websocket.Conn) *MwsyClientSes {
	this := &MwsyClientSes{}
	this.m_WebConn = webcon
	this.m_Client = client
	return this
}
