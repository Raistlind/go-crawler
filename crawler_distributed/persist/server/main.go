package main

import (
	"gopkg.in/olivere/elastic.v5"
	"GolandProjects/goexercises/crawler_distributed/rpcsupport"
	"GolandProjects/goexercises/crawler_distributed/persist"
	"github.com/gpmgo/gopm/modules/log"
	"fmt"
	"GolandProjects/goexercises/crawler_distributed/config"
)

func main() {
	log.Fatal("", serveRpc(fmt.Sprintf(":d", config.ItemSaverPort), config.ElasticIndex))
}

func serveRpc(host, index string) error {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return err
	}

	return rpcsupport.ServeRpc(host, &persist.ItemSaverService{
		Client: client,
		Index:  index,
	})
}
