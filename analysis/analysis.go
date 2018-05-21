package main

import (
	"flag"
	"time"
	"github.com/sirupsen/logrus"
	"os"
	"bufio"
	"io"
	"strings"
	"github.com/mgutz/str"
	"net/url"
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"github.com/mediocregopher/radix.v2/pool"
)

const HANDLE_DIG = " /dig?"
const HANDLE_MOVIE = "/movie/"
const HANDLE_LIST = "/list/"
const HANDLE_HTML = ".html"

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
	data  digData
	uid   string
	unode urlNode
}

type urlNode struct {
	unType string
	unRid  int
	unUrl  string
	unTime string
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
	routineNum := flag.Int("routineNum", 10, "consumer number by goroutine")
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

	redisPool, err := pool.New("tcp", "localhost:6379", 2*params.routineNum)
	if err != nil {
		log.Fatalln("Redis pool created failed.")
		panic(err)
	} else {
		go func() {
			for {
				redisPool.Cmd("PING")
				time.Sleep(3 * time.Second)
			}
		}()
	}

	go readFileLineByLine(params, logChannel)

	for i := 0; i < params.routineNum; i ++ {
		go logConsumer(logChannel, pvChannel, uvChannel)
	}

	go pvCounter(pvChannel, storageChannel)
	go uvCounter(uvChannel, storageChannel, redisPool)

	go dataStorage(storageChannel, redisPool)

	time.Sleep(1000 * time.Second)

}

func dataStorage(storageChannel chan storageBlock, redisPool *pool.Pool) {
	for block := range storageChannel {
		prefix := block.counterType + "_"
		setKeys := []string{
			prefix + "day_" + getTime(block.unode.unTime, "day"),
			prefix + "hour_" + getTime(block.unode.unTime, "hour"),
			prefix + "min_" + getTime(block.unode.unTime, "min"),
			prefix + block.unode.unType + "_day_" + getTime(block.unode.unTime, "day"),
			prefix + block.unode.unType + "_hour_" + getTime(block.unode.unTime, "hour"),
			prefix + block.unode.unType + "_min_" + getTime(block.unode.unTime, "min"),
		}

		rowId := block.unode.unRid

		for _, key := range setKeys {
			ret, err := redisPool.Cmd(block.storageModel, key, 1, rowId).Int()
			if ret <= 0 || err != nil {
				log.Errorln("DataStorage redis storage error.", block.storageModel, key, rowId)
			}
		}
	}
}

func pvCounter(pvChannel chan urlData, storageChannel chan storageBlock) {
	for data := range pvChannel {
		sItem := storageBlock{"pv", "ZINCRBY", data.unode,}
		storageChannel <- sItem
	}
}

func uvCounter(uvChannel chan urlData, storageChannel chan storageBlock, redisPool *pool.Pool) {
	for data := range uvChannel {

		hyperLogLogKey := "uv_hpll_" + getTime(data.data.time, "day")
		ret, err := redisPool.Cmd("PFADD", hyperLogLogKey, data.uid, "EX", 86400).Int()
		if err != nil {
			log.Warningln("UvCounter check redia hyperloglog failed, ", err)
		}
		if ret != 1 {
			continue
		}

		sItem := storageBlock{"uv", "ZINCRBY", data.unode,}
		storageChannel <- sItem
	}
}

func logConsumer(logChannel chan string, pvChannel, uvChannel chan urlData) error {
	for logStr := range logChannel {
		data := cutLogFetchData(logStr)

		hasher := md5.New()
		hasher.Write([]byte(data.refer + data.ua))
		uid := hex.EncodeToString(hasher.Sum(nil))

		uData := urlData{data, uid, formatUrl(data.url, data.time)}

		log.Infoln("logConsumer:", data, uid)

		pvChannel <- uData
		uvChannel <- uData
	}
	return nil
}

func cutLogFetchData(logStr string) digData {
	logStr = strings.TrimSpace(logStr)
	pos1 := str.IndexOf(logStr, HANDLE_DIG, 0)
	if pos1 == -1 {
		return digData{}
	}
	pos1 += len(HANDLE_DIG)
	pos2 := str.IndexOf(logStr, " HTTP/", pos1)
	d := str.Substr(logStr, pos1, pos2-pos1)

	urlInfo, err := url.Parse("http://localhost/?" + d)
	if err != nil {
		return digData{}
	}
	data := urlInfo.Query()
	return digData{
		time:  data.Get("time"),
		url:   data.Get("url"),
		refer: data.Get("refer"),
		ua:    data.Get("ua"),
	}
}

func readFileLineByLine(parames cmdParams, logChannel chan string) error {
	fd, err := os.Open(parames.logFilePath)
	if err != nil {
		log.Warningf("ReadFileLineByLine can't open file:%s", parames.logFilePath)
		return err
	}
	defer fd.Close()

	count := 0
	bufferRead := bufio.NewReader(fd)
	for {
		line, err := bufferRead.ReadString('\n')
		logChannel <- line
		log.Infoln("line:", line)
		count ++

		if count%(1000*parames.routineNum) == 0 {
			log.Infof("ReadFileLineByLien line: %d", count)
		}

		if err != nil {
			if err == io.EOF {
				time.Sleep(3 * time.Second)
				log.Infof("ReadFileLineByLine wait, readline:%d", count)
			} else {
				log.Warningf("ReadFileLineByLine read log error:%s", err)
			}
		}
	}
	return nil
}

func formatUrl(url, t string) urlNode {
	pos1 := str.IndexOf(url, HANDLE_MOVIE, 0)
	if pos1 != -1 {
		pos1 += len(HANDLE_MOVIE)
		pos2 := str.IndexOf(url, HANDLE_HTML, 0)
		idStr := str.Substr(url, pos1, pos2-pos1)
		id, _ := strconv.Atoi(idStr)
		return urlNode{"movie", id, url, t}
	} else {
		pos1 = str.IndexOf(url, HANDLE_LIST, 0)
		if pos1 != -1 {
			pos1 += len(HANDLE_LIST)
			pos2 := str.IndexOf(url, HANDLE_HTML, 0)
			idStr := str.Substr(url, pos1, pos2-pos1)
			id, _ := strconv.Atoi(idStr)
			return urlNode{"list", id, url, t}
		} else {
			return urlNode{"home", 1, url, t}
		}
	}
}

func getTime(logTime, timeType string) string {
	var item string
	switch timeType {
	case "day":
		item = "2006-01-02"
		break
	case "hour":
		item = "2006-01-02 15"
		break
	case "min":
		item = "2006-01-02 15:03"
		break
	}
	t, _ := time.Parse(item, time.Now().Format(item))
	return strconv.FormatInt(t.Unix(), 10)
}
