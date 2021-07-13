package main

import (
	"bufio"
	"mwsy/mwsyclient"
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

func main(){
	for i:= 1;i<=100;i++ {
		sv := mwsyclient.NewMwsyClient()
		sv.Start()
	}
	ReadConsoleExit()
}