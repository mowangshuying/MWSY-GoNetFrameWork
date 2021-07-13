package main

import (
	"bufio"
	"mwsy/mwsyserver"
	//"net/http"
	//_ "net/http/pprof"
	"os"
	"strings"
)


func ReadConsoleExit() {

	for {

		// 从标准输入读取字符串，以\n为分割
		text, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			break
		}

		// 去掉读入内容的空白符
		text = strings.TrimSpace(text)

		if text == "exit" {
			break
		}
	}
}

//func RunWebSer(){
//	go http.ListenAndServe("localhost:6060", nil)
//}

func main(){
	//RunWebSer()
	sv := mwsyserver.NewMwsyServer()
	sv.Start()
	ReadConsoleExit()
}