package mwsymsg

import "fmt"

//1
type ChatStruct struct{
	Msg string
	Value int32
}

func (self *ChatStruct) String() string { return fmt.Sprintf("%+v", *self) }

//2