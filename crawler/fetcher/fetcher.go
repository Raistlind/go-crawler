package fetcher

import (
	"net/http"
	"fmt"
	"golang.org/x/text/transform"
	"io/ioutil"
	"golang.org/x/text/encoding"
	"bufio"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/unicode"
	"log"
	"time"
	"encoding/json"
	"GolandProjects/goexercises/crawler/config"
)

var rateLimiter = time.Tick(time.Second / config.Qps)

func Fetch(url string) ([]byte, error) {
	<-rateLimiter
	log.Printf("Fetching url %s", url)
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)

	request.Header.Add("http-equiv", "Content-Type")
	request.Header.Add("Content", "text/html")
	request.Header.Add("charset", "gbk")
	request.Header.Add(
		"User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) "+
			"AppleWebKit/537.36 (KHTML, like Gecko) Chrome/34.0.1847.137 Safari/537.36 LBBROWSER")

	if err != nil {
		panic(err)
	}
	resp, err := client.Do(request)
	if err != nil {
		log.Println("there is a error：clinet.do break:: ",err)
		log.Println(json.Marshal(resp))
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error: status code", resp.StatusCode)
		return nil, fmt.Errorf("wrong status code: %d", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := determineEncoding(bodyReader)
	utf8Reader := transform.NewReader(resp.Body, e.NewDecoder())

	defer resp.Body.Close()
	return ioutil.ReadAll(utf8Reader)
}

func determineEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		log.Printf("Fetcher error: %v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
