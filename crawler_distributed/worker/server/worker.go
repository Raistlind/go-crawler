package main

import (
	"log"
	"GolandProjects/goexercises/crawler_distributed/rpcsupport"
	"fmt"
	"GolandProjects/goexercises/crawler_distributed/config"
	"GolandProjects/goexercises/crawler_distributed/worker"
)

func main() {
	log.Fatal(rpcsupport.ServeRpc(
		fmt.Sprintf(":%d", config.WorkerPort0),
		worker.CrawlService{}))
}
