package main

import (
	"regexp"
	"fmt"
)

const text = `
My email is ccmouse@gmail.com
email is abc@def.org
email2 is    kkk@qq.com
email3 is  ddd@abc.com.cn
`

func main() {
	re := regexp.MustCompile(`([a-zA-Z0-9]+)@([a-zA-Z0-9.]+)(\.[a-zA-Z0-9]+)`)
	match := re.FindAllStringSubmatch(text, -1)
	fmt.Println(match)
	for i := range match {
		fmt.Println(match[i])
	}
}
