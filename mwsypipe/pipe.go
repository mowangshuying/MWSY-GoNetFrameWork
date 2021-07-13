package mwsypipe

import "sync"

type MwsyPipe struct {
	m_PipeElemList []interface{}
	m_Mutex    *sync.Mutex
	m_Cond     *sync.Cond
}

func (this *MwsyPipe) Add(elem interface{}){
	this.m_Mutex.Lock()
	this.m_PipeElemList = append(this.m_PipeElemList, elem)
	this.m_Mutex.Unlock()
	this.m_Cond.Signal()
}

func (this *MwsyPipe) Count() int {
	this.m_Mutex.Lock()
	defer this.m_Mutex.Unlock()
	return len(this.m_PipeElemList)
}

func (this *MwsyPipe) Reset() {
	this.m_Mutex.Lock()
	this.m_PipeElemList = this.m_PipeElemList[0:0]
	this.m_Mutex.Unlock()
} 

func (this *MwsyPipe) Pick(retList *[]interface{})(exit bool){
	this.m_Mutex.Lock()

	for len(this.m_PipeElemList) == 0 {
		this.m_Cond.Wait()
	}

	this.m_Mutex.Unlock()
	this.m_Mutex.Lock()

	for _,data :=  range this.m_PipeElemList {
		if data == nil {
			exit = true
			break
		}else {
			*retList = append(*retList, data)
		}
	}

	this.m_PipeElemList = this.m_PipeElemList[0:0]
	this.m_Mutex.Unlock()
	return
}

func NewMwsyPipe() *MwsyPipe {
	this := &MwsyPipe{}
	this.m_Mutex = &sync.Mutex{}
	this.m_Cond = sync.NewCond(this.m_Mutex)
	return this
}


