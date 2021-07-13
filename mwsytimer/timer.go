package mwsytimer

import (
	"fmt"
	"time"
)


type MwsyTimer struct {
	m_ticker *time.Ticker
	m_dur    time.Duration
	m_callback func()
	m_bRuning bool
}

func (this *MwsyTimer)Start(){
	this.m_bRuning = true
	for i := range this.m_ticker.C {
		if this.m_bRuning {
			this.m_callback()
		}else{
			//this.Stop()
			this.m_ticker.Stop()
			fmt.Println("tickstop and i.Unix() = ",i.Unix())
			break
		}
	}
}

func (this *MwsyTimer) Stop(){
	//this.m_ticker.Stop()
	this.m_bRuning = false
}

func NewMwsyTimer(dur time.Duration,callback func()) *MwsyTimer{
	this := &MwsyTimer{}
	this.m_dur = dur
	this.m_ticker = time.NewTicker(this.m_dur)
	this.m_callback = callback
	this.m_bRuning = false
	return this
}
