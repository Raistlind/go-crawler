package main

import (
	"flag"
	"fmt"
	"strings"
	"strconv"
)

type resource struct {
	url    string
	target string
	start  int
	end    int
}

func ruleResource() []resource {
	var res []resource
	r1 := resource{
		url:    "http://localhost:8888/",
		target: "",
		start:  0,
		end:    0,
	}

	r2 := resource{
		url:    "http://localhost:8888/list/{$id}.html",
		target: "{$id}",
		start:  1,
		end:    21,
	}

	r3 := resource{
		url:    "http://localhost:8888/movie/{$id}.html",
		target: "{$id}",
		start:  1,
		end:    12924,
	}

	res = append(res, r1, r2, r3)
	return res
}

func buildUrl(res []resource) []string {
	var list []string

	for _, resItem := range res {
		if len(resItem.target) == 0 {
			list = append(list, resItem.url)
		} else {
			for i := resItem.start; i <= resItem.end; i++ {
				urlStr := strings.Replace(resItem.url, resItem.target, strconv.Itoa(i), -1)
				list = append(list, urlStr)
			}
		}
	}
	return list
}

func main() {
	total := flag.Int("total", 100, "how many")
	filePath := flag.String("filePath", "/usr/local/nginx/logs/dig.log", "file path")
	flag.Parse()

	res := ruleResource()
	list := buildUrl(res)

	fmt.Println(total, filePath, list)
	fmt.Println("done. ")
}
