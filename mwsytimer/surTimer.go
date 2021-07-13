package mwsytimer

//超级时间：提供注册函数
import (
	"fmt"
	"time"
)


type MwsySurTimer struct {
	m_ticker *time.Ticker //ticker
	m_dur    time.Duration//最原始的那个间隔
	m_TickCount2Fun map[int] []func()
	m_nTickCount int
	m_bRuning bool
}


func (this *MwsySurTimer)RegFun(nCount int,callback func()){
	if len(this.m_TickCount2Fun[nCount]) == 0 {
		funarr := []func(){}
		funarr = append(funarr, callback)
		this.m_TickCount2Fun[nCount] = funarr

	} else {
		this.m_TickCount2Fun[nCount] =  append(this.m_TickCount2Fun[nCount],callback)
	}
}

func (this *MwsySurTimer)Start(){
	this.m_bRuning = true
	for i := range this.m_ticker.C {
		if this.m_bRuning {
			this.m_nTickCount++
			for key,value := range this.m_TickCount2Fun {
				if this.m_nTickCount % key == 0 {
					for _,fun := range value {
						fun()
					}
				}
			}
		}else{
			this.m_ticker.Stop()
			fmt.Println("tickstop and i.Unix() = ",i.Unix())
			break
		}
	}
}

func (this *MwsySurTimer) Stop(){
	this.m_bRuning = false
}

func NewMwsySurTimer(dur time.Duration) *MwsySurTimer{
	this := &MwsySurTimer{}
	this.m_dur = dur
	this.m_ticker = time.NewTicker(this.m_dur)
	this.m_bRuning = false
	this.m_nTickCount = 0
	this.m_TickCount2Fun = make(map[int] []func(), 0)
	return this
}
