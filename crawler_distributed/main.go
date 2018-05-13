package main

import (
	"GolandProjects/goexercises/crawler/engine"
	"GolandProjects/goexercises/crawler/scheduler"
	"GolandProjects/goexercises/crawler/zhenai/parser"
	itemsaver "GolandProjects/goexercises/crawler_distributed/persist/client"
	"GolandProjects/goexercises/crawler_distributed/config"
	worker "GolandProjects/goexercises/crawler_distributed/worker/client"
	"net/rpc"
	"GolandProjects/goexercises/crawler_distributed/rpcsupport"
	"log"
	"flag"
	"strings"
)

var (
	itemSaverHost = flag.String("itemsaver_host", "", "itemsaver host")
	workerHosts   = flag.String("worker_hosts", "", "worker hosts (comma separated)")
)

func main() {
	flag.Parse()

	itemChan, err := itemsaver.ItemSaver(*itemSaverHost)
	if err != nil {
		panic(err)
	}

	pool := createClientPool(strings.Split(*workerHosts, ","))

	processor := worker.CreateProcessor(pool)

	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan,
		RequestProcessor: processor,
	}

	e.Run(engine.Request{
		Url:    "http://www.zhenai.com/zhenghun",
		Parser: engine.NewFuncParser(parser.ParseCityList, config.ParseCityList),
	})
}

func createClientPool(hosts []string) chan *rpc.Client {
	var clients []*rpc.Client
	for _, h := range hosts {
		client, err := rpcsupport.NewClient(h)
		if err == nil {
			clients = append(clients, client)
			log.Printf("Connected to %s", h)
		} else {
			log.Printf("Error connecting to %s: %v", h, err)
		}
	}

	out := make(chan *rpc.Client)
	go func() {
		var count int = 0
		for {
			for i, client := range clients {
				out <- client
				log.Printf("append a CLIENT(%d) to out channel<<<<<<<<"+
					"<<<<<<<<<<<<< %d, clients len: %d", i, count, len(clients))
				count++
			}
		}
	}()
	return out
}
