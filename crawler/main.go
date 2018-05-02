package main

import (
	"GolandProjects/goexercises/crawler/engine"
	"GolandProjects/goexercises/crawler/zhenai/parser"
)

func main() {
	engine.Run(engine.Request{
		Url: "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})
}
