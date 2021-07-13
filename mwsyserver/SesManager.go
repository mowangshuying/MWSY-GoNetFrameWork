package mwsyserver

import (
	"sync"
	"sync/atomic"
)

type MwsySesManager struct {
	// [会话id : 会话]
	m_Id2Ses sync.Map
	// 最大会话id
	m_nMaxId int64
	// 会话数量
	m_nCount int64
}

func (this *MwsySesManager) SetMaxId(id int64) {
	atomic.StoreInt64(&this.m_nMaxId, 1)
}

func (this *MwsySesManager) getCount() int {
	return int(atomic.LoadInt64(&this.m_nCount))
}

func (this *MwsySesManager) AddSes(ses *MwsyServerSes) {
	id := atomic.AddInt64(&this.m_nMaxId, 1)
	atomic.AddInt64(&this.m_nCount,1)
	ses.SetId(id)
	this.m_Id2Ses.Store(id, ses)
}

func (this *MwsySesManager) DelSes(ses *MwsyServerSes) {
	this.m_Id2Ses.Delete(ses.GetId())
	atomic.AddInt64(&this.m_nCount,-1)
}

func (this *MwsySesManager) GetSesById(id int64) *MwsyServerSes {
	v, ok := this.m_Id2Ses.Load(id)
	if ok {
		return v.(*MwsyServerSes)
	}
	return nil
}
