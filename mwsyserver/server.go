package mwsyserver

import (
	//"encoding/json"
	"fmt"
	"mwsy/mwsylog"
	//_ "mwsy/mwsylog"
	"net"
	"net/http"

	//"strconv"
	"encoding/binary"
	"encoding/json"
	"mwsy/mwsymsg"
	"mwsy/mwsypipe"

	"github.com/gorilla/websocket"
)

type MsgServerFun func(*MwsyServer,*MwsyServerSes, interface{})

type MsgIDFun func(msg interface{}) 

type MwsyServerMsg struct {
	M_Ses  *MwsyServerSes
	M_Data interface{}
}


type MwsyServer struct {
	m_Upgrader websocket.Upgrader
	m_Listener net.Listener
	m_httpSv   *http.Server
	//会话管理器
	m_SesManerger MwsySesManager
	//发送队列 [会话id:消息]
	m_SendQueue *mwsypipe.MwsyPipe
	//接收队列 [会话id:消息]
	m_RecvQueue *mwsypipe.MwsyPipe
	//
	//处理接收到的消息
	m_OnMsgFun MsgServerFun
	//[消息id:处理消息函数]
	m_MsgIDFunMap map[int32]MsgIDFun
}

//
func (this *MwsyServer) RegMsgFun(msgid int32,msgIDFun MsgIDFun){
	this.m_MsgIDFunMap[msgid] = msgIDFun
}

func (this *MwsyServer) UnRegMsgFun(msgid int32,msgIDFun MsgIDFun){
	delete(this.m_MsgIDFunMap,msgid)
}

func (this *MwsyServer) ProcMsgFun(msgid int32,msgObj interface{}){
	if this.m_MsgIDFunMap[msgid] != nil {
		this.m_MsgIDFunMap[msgid](msgObj)
	}
}

//
func (this *MwsyServer) AddSes(ses *MwsyServerSes) {
	this.m_SesManerger.AddSes(ses)
}

func (this *MwsyServer) DelSes(ses *MwsyServerSes) {
	this.m_SesManerger.DelSes(ses)
}

//添加到队列中
func (this *MwsyServer) Send(ses *MwsyServerSes, msgObj interface{}) {
	msg := MwsyServerMsg{M_Ses: ses, M_Data: msgObj}
	this.m_SendQueue.Add(msg)
}

func (this *MwsyServer) Recv(ses *MwsyServerSes, msgObj interface{}) {
	msg := MwsyServerMsg{M_Ses: ses, M_Data: msgObj}
	this.m_RecvQueue.Add(msg)
}

//处理将发送的消息
func (this *MwsyServer) SendLoop() {
	//取出队列中数据
	var writeList []interface{}
	for {
		writeList = writeList[0:0]
		this.m_SendQueue.Pick(&writeList)
		for _, msg := range writeList {
			ses := msg.(MwsyServerMsg).M_Ses
			data := msg.(MwsyServerMsg).M_Data
			bytearr := data.([]byte)
			if bytearr != nil {
				webcon := ses.GetWebCon()
				webcon.WriteMessage(websocket.BinaryMessage, bytearr)
			}
		}
	}
}

//处理接收到的消息
func (this *MwsyServer) RecvLoop() {
	var writeList []interface{}
	for {
		writeList = writeList[0:0]
		this.m_RecvQueue.Pick(&writeList)
		//fmt.Println("writelist:",len(writeList))
		for _, msg := range writeList {
			ses := msg.(MwsyServerMsg).M_Ses
			data := msg.(MwsyServerMsg).M_Data
			this.m_OnMsgFun(this,ses, data)
		}
	}
}

//
func (this *MwsyServer) Start() {
	//开启服务器
	var err error
	listener, err := net.Listen("tcp", "127.0.0.1:18802")
	this.m_Listener = listener
	if err != nil {
		return
	}

	//
	mux := http.NewServeMux()
	this.m_Upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//升级成功后将获取到WebSocket.Conn 利用这个Conn可进行消息收发
		webcon, err := this.m_Upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Println("upgrade error")
			return
		}

		//创建ses
		pSes := NewMwsyServerSession(this, webcon)
		pSes.Start()

	})

	//
	this.m_httpSv = &http.Server{Addr: "127.0.0.1:18802", Handler: mux}

	//等待链接
	go this.m_httpSv.Serve(this.m_Listener)

	//发送循环
	go this.SendLoop()

	//接收循环
	go this.RecvLoop()
}


//创建一个服务
func NewMwsyServer() *MwsyServer {
	this := &MwsyServer{}
	this.m_OnMsgFun = func(sv *MwsyServer,ses *MwsyServerSes, msg interface{}) {
		ba := msg.([]byte)
		msgID := binary.LittleEndian.Uint32(ba)
		msgBody := ba[4:]
		sv.ProcMsgFun(int32(msgID),msgBody)
	}
	this.m_RecvQueue = mwsypipe.NewMwsyPipe()
	this.m_SendQueue = mwsypipe.NewMwsyPipe()
	this.m_MsgIDFunMap = make(map[int32]MsgIDFun)

	this.RegMsgFun(1234,func(msgx interface{}){
		var msgObj mwsymsg.ChatStruct
		json.Unmarshal(msgx.([]byte),&msgObj)
		logstr := fmt.Sprintf("ses cont = %v msg = %v",this.m_SesManerger.getCount(),msgObj)
		mwsylog.LogInst.MwsyLogNormalPrint(logstr)
	})

	return this
}