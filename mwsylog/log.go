package mwsylog

import (
	"bufio"
	"fmt"
	"mwsy/mwsypipe"
	"mwsy/mwsytimer"
	"os"
	"time"
)


type MwsyLog_Level int

//日志级别
const (
	MwsyLog_Normal MwsyLog_Level = iota
	MwsyLog_Warnning
	MwsyLog_Error
)

// [时间] [日志级别] [日志详细信息]

type MwsyLog struct{
	//管道信息
	m_Pipe *mwsypipe.MwsyPipe
	//文件路径名
	m_FilePath string
	//文件
	m_File *os.File
	//
	m_Writer *bufio.Writer
	//
	m_Timer *mwsytimer.MwsyTimer
}

func (this *MwsyLog) GetTime() string {
	NowTime := time.Now()
	year := NowTime.Year()     //年
	month := NowTime.Month()   //月
	day := NowTime.Day()       //天
	hour := NowTime.Hour()     //时
	minute := NowTime.Minute() //分
	second := NowTime.Second() //秒
	return fmt.Sprintf("[%d-%02d-%02d %02d:%02d:%02d] ",year,month,day,hour,minute,second)
}

func (this *MwsyLog) MwsyLogNormalPrint(logcontext string){
	this.m_Pipe.Add(func(){
		sLog := this.GetTime()
		sLog += " [Normal] "
		sLog += logcontext
		sLog += "\n"
		this.m_Writer.WriteString(sLog)
		this.m_Writer.Flush()
		fmt.Printf(sLog)
	})
}

func (this *MwsyLog) MwsyLogWarnningPrint(logcontext string){
	this.m_Pipe.Add(func(){
		sLog := this.GetTime()
		sLog += " [Warnning] "
		sLog += logcontext
		sLog += "\n"
		this.m_Writer.WriteString(sLog)
		this.m_Writer.Flush()
		fmt.Printf(sLog)
	})
}

func (this *MwsyLog) MwsyLogErrorPrint(logcontext string){
	this.m_Pipe.Add(func(){
		sLog := this.GetTime()
		sLog += " [Error] "
		sLog += logcontext
		sLog += "\n"
		this.m_Writer.WriteString(sLog)
		fmt.Printf(sLog)
	})
}

func (this *MwsyLog) Start(){
	var writeList []interface{}
	for {
		writeList = writeList[0:0]
		this.m_Pipe.Pick(&writeList)
		for _, callback := range writeList {
			callback.(func()())()
		}
	}
	this.m_File.Close()
}


func NewMwsyLog(filePath string) *MwsyLog{
	this := &MwsyLog{}
	this.m_FilePath = filePath
	var err error
	this.m_File, err = os.OpenFile(this.m_FilePath, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0666)
	if err!=nil{
		fmt.Println("file open error:",err)
		return nil
	}
	this.m_Writer = bufio.NewWriter(this.m_File)
	this.m_Pipe = mwsypipe.NewMwsyPipe()
	this.m_Timer = mwsytimer.NewMwsyTimer(60*time.Second,func() {
		if this.m_Pipe!=nil {
			this.m_Writer.Flush()
		}
	})

	return this
}


var LogInst *MwsyLog
func init(){
	LogInst=NewMwsyLog("./example.log")
	go LogInst.Start()
	go LogInst.m_Timer.Start()
}