package main

import (
	"GolandProjects/goexercises/crawler/engine"
	"GolandProjects/goexercises/crawler/scheduler"
	"GolandProjects/goexercises/crawler/zhenai/parser"
	"GolandProjects/goexercises/crawler_distributed/persist/client"
	"fmt"
	"GolandProjects/goexercises/crawler_distributed/config"
)

func main() {
	itemChan, err := client.ItemSaver(fmt.Sprintf(":d", config.ItemSaverPort))
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
}
