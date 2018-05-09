package parser

import (
	"GolandProjects/goexercises/crawler/engine"
	"regexp"
)

var profileRe = regexp.MustCompile(`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`)
var cityUrlRe = regexp.MustCompile(`href="(http://www.zhenai.com/zhenghun/shanghai/[^"]+)"`)

//`<a href="(http://album.zhenai.com/u/[0-9]+)"[^>]*>([^<]+)</a>`

func ParseCity(contents []byte) engine.ParseResult {
	matches := profileRe.FindAllSubmatch(contents, -1)

	result := engine.ParseResult{}
	for _, m := range matches {
		url := string(m[1])
		name := string(m[2])
		//result.Items = append(result.Items, "User "+name)
		result.Requests = append(result.Requests, engine.Request{
			Url:        url,
			ParserFunc: ProfileParser(name, url),
		})
	}

	matches = cityUrlRe.FindAllSubmatch(contents, -1)
	for _, m := range matches {
		result.Requests = append(result.Requests, engine.Request{
			Url:        string(m[1]),
			ParserFunc: ParseCity,
		})
	}
	return result
}
