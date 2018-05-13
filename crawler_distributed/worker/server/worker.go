package main

import (
	"log"
	"GolandProjects/goexercises/crawler_distributed/rpcsupport"
	"fmt"
	"GolandProjects/goexercises/crawler_distributed/worker"
	"flag"
)

var port = flag.Int("port", 0, "the port for me to listen on")

func main() {
	flag.Parse()
	if *port == 0 {
		fmt.Println("must specify a port")
		return
	}
	log.Fatal(rpcsupport.ServeRpc(
		fmt.Sprintf(":%d", *port),
		worker.CrawlService{}))
}
