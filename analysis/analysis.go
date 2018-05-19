package main

import "flag"

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

func main() {
	logFilePath := flag.String("logFilePath", "/usr/local/nginx/logs/dig.log", "log file path")
	routineNum := flag.Int("routineNum", 5, "consumer number by goroutine")
	l := flag.String("l", "/tmp/log", "this programe runtime log target file path")
	flag.Parse()

	params := cmdParams{
		logFilePath: *logFilePath,
		routineNum:  *routineNum,
	}

	var logChannel = make(chan string, 3 * *routineNum)
	var pvChannel = make(chan urlData, *routineNum)
	var uvChannel = make(chan urlData, *routineNum)
	var storageChannel = make(chan storageData, *routineNum)

	go readFileLineByLine(params, logChannel)

	for i := 0; i < params.routineNum; i ++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

}

func logConsumer(logChannel chan string, pvChannel, uvChannel chan urlData) {

}

func readFileLineByLine(parmes cmdParams, logChannel chan string) {

}
