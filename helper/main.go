package main

import (
	"fmt"
	"helper/toolx"
	"os"
)

func main() {

	base, err := os.Getwd()
	if err == nil {
		err = toolx.FormatReadMe(base+"/README.md.template", base+"/README_NEW.md")
	}
	fmt.Println(err)
}
