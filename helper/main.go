package main

import (
	"fmt"
	"helper/toolx"
	"os"
)

func main() {

	base, err := os.Getwd()
	if err == nil {
		err = toolx.FormatReadMe(base+"/README.md.template", base+"/README.md")
	}

	if err == nil {
		fmt.Println("生成成功!")
	} else {
		fmt.Println("出错" + err.Error())
	}
}
