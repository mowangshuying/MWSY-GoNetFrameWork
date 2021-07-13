package main

import (
	"bufio"
	"fmt"
	"mwsy/mwsytimer"
	"os"
	"strings"
	"time"
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
	surTimer := mwsytimer.NewMwsySurTimer(time.Second)
	go surTimer.Start()
	surTimer.RegFun(1,func(){
		fmt.Println("this is 1")
	})

	//surTimer.RegFun(1,func(){
	//	fmt.Println("this is 1 and 2")
	//})

	surTimer.RegFun(3,func(){
		fmt.Println("this is 3")
	})

	surTimer.RegFun(5,func(){
		fmt.Println("this is 5")
	})

	ReadConsoleExit()
}