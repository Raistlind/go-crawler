package main

import (
	"flag"
	"time"
	"github.com/sirupsen/logrus"
	"os"
)

type cmdParams struct {
	logFilePath string
	routineNum  int
}

type digData struct {
	time  string
	url   string
	refer string
	ua    string
}

type urlData struct {
	data digData
	uid  string
}

type urlNode struct {
}

type storageBlock struct {
	counterType  string
	storageModel string
	unode        urlNode
}

var log = logrus.New()

func init() {
	log.Out = os.Stdout
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	logFilePath := flag.String("logFilePath", "/usr/local/nginx/logs/dig.log", "log file path")
	routineNum := flag.Int("routineNum", 5, "consumer number by goroutine")
	l := flag.String("l", "/tmp/log", "this programe runtime log target file path")
	flag.Parse()

	params := cmdParams{
		logFilePath: *logFilePath,
		routineNum:  *routineNum,
	}

	logFd, err := os.OpenFile(*l, os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		log.Out = logFd
		defer logFd.Close()
	}
	log.Infof("Exec start.")
	log.Infof("Params: logFilePath=%s, routineNum=%d", params.logFilePath, params.routineNum)

	var logChannel = make(chan string, 3*params.routineNum)
	var pvChannel = make(chan urlData, params.routineNum)
	var uvChannel = make(chan urlData, params.routineNum)
	var storageChannel = make(chan storageBlock, params.routineNum)

	go readFileLineByLine(params, logChannel)

	for i := 0; i < params.routineNum; i ++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

	go pvCounter(pvChannel, storageChannel)
	go uvCounter(uvChannel, storageChannel)

	go dataStorage(storageChannel)

	time.Sleep(1000 * time.Second)

}

func dataStorage(storageChannel chan storageBlock) {

}

func pvCounter(pvChannel chan urlData, storageChannel chan storageBlock) {

}

func uvCounter(uvChannel chan urlData, storageChannel chan storageBlock) {

}

func logConsumer(logChannel chan string, pvChannel, uvChannel chan urlData) {

}

func readFileLineByLine(parmes cmdParams, logChannel chan string) {

}
