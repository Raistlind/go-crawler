package main

import (
	"GolandProjects/goexercises/crawler/engine"
	"GolandProjects/goexercises/crawler/scheduler"
	"GolandProjects/goexercises/crawler/persist"
	"GolandProjects/goexercises/crawler/zhenai/parser"
)

func main() {
	itemChan, err := persist.ItemSaver()
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler:   &scheduler.QueuedScheduler{},
		WorkerCount: 100,
		ItemChan:    itemChan,
	}

	e.Run(engine.Request{
		Url:        "http://www.zhenai.com/zhenghun",
		ParserFunc: parser.ParseCityList,
	})

	//e.Run(engine.Request{
	//	Url:        "http://www.zhenai.com/zhenghun/shanghai",
	//	ParserFunc: parser.ParseCity,
	//})
}
