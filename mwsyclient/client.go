package mwsyclient

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"mwsy/mwsylog"
	"mwsy/mwsymsg"
	"mwsy/mwsypipe"
	"mwsy/mwsytimer"

	//"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)


type MsgClientFun func(*MwsyClient,*MwsyClientSes, interface{})

type MsgIDFun func(msg interface{}) 

type MwsyClientMsg struct {
	M_Ses  *MwsyClientSes
	M_Data interface{}
}

type MwsyClient struct {
	m_Ses *MwsyClientSes
	//发送队列 [会话id:消息]
	m_SendQueue *mwsypipe.MwsyPipe
	//接收队列 [会话id:消息]
	m_RecvQueue *mwsypipe.MwsyPipe
	//处理接收到的消息
	m_OnMsgFun MsgClientFun
	//[消息id:处理消息函数]
	m_MsgIDFunMap map[int32]MsgIDFun
}

func (this *MwsyClient) RegMsgFun(msgid int32,msgIDFun MsgIDFun){
	this.m_MsgIDFunMap[msgid] = msgIDFun
}

func (this *MwsyClient) UnRegMsgFun(msgid int32,msgIDFun MsgIDFun){
	delete(this.m_MsgIDFunMap,msgid)
}

func (this *MwsyClient) ProcMsgFun(msgid int32,msgObj interface{}){
	if this.m_MsgIDFunMap[msgid] != nil {
		this.m_MsgIDFunMap[msgid](msgObj)
	}
}

func (this *MwsyClient) Send(ses *MwsyClientSes,msgObj interface{}){
	msg := MwsyClientMsg{M_Ses: ses,M_Data: msgObj}
	this.m_SendQueue.Add(msg)
}

func (this *MwsyClient) Recv(ses *MwsyClientSes,msgObj interface{}){
	msg := MwsyClientMsg{M_Ses: ses,M_Data: msgObj}
	this.m_RecvQueue.Add(msg)
}

//处理将发送的消息
func (this *MwsyClient) SendLoop() {
	//取出队列中数据
	var writeList []interface{}
	for {
		writeList = writeList[0:0]
		this.m_SendQueue.Pick(&writeList)
		for _, msg := range writeList {
			ses := msg.(MwsyClientMsg).M_Ses
			data := msg.(MwsyClientMsg).M_Data
			bytearr := data.([]byte)
			if bytearr != nil {
				webcon := ses.GetWebCon()
				webcon.WriteMessage(websocket.BinaryMessage, bytearr)
			}
		}
	}
}

//处理接收到的消息
func (this *MwsyClient) RecvLoop() {
	var writeList []interface{}
	for {
		writeList = writeList[0:0]
		this.m_RecvQueue.Pick(&writeList)
		for _, msg := range writeList {
			ses := msg.(MwsyClientMsg).M_Ses
			data := msg.(MwsyClientMsg).M_Data
			this.m_OnMsgFun(this,ses, data)
		}
	}
}

//
func (this *MwsyClient) Start(){
	//
	dialer := websocket.Dialer{}
	dialer.Proxy = http.ProxyFromEnvironment
	dialer.HandshakeTimeout = 60 * time.Second

	
	//conn, err := net.Dial("tcp", "127.0.0.1:18802")
	
	conn, _, err := dialer.Dial("ws://127.0.0.1:18802", nil)
	
	//
	if err != nil {
		fmt.Println("connect server error")
		return
	}

	this.m_Ses = NewMwsyClientSes(this,conn)

	go func(){
		timer := mwsytimer.NewMwsyTimer(time.Second*2,func() {
			msgObj := mwsymsg.ChatStruct{Msg:"i am msg",Value:  1234}
			bytearr,_ := json.Marshal(msgObj)
			pkt := make([]byte, 4+len(bytearr))
			binary.LittleEndian.PutUint32(pkt, 1234)
			copy(pkt[4:], bytearr)
		   this.Send(this.m_Ses,pkt)
		})
		timer.Start()
	}()
	
	go this.SendLoop()
	go this.RecvLoop()
}


//创建一个服务
func NewMwsyClient() *MwsyClient {
	this := &MwsyClient{}
	this.m_OnMsgFun = func(sv *MwsyClient,ses *MwsyClientSes, msg interface{}) {
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
		mwsylog.LogInst.MwsyLogNormalPrint(msgObj.String())
	})

	return this
}





