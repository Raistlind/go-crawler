package engine

type ParseFunc func([]byte) ParseResult

type Request struct {
	Url        string
	ParserFunc ParseFunc
}

type ParseResult struct {
	Requests []Request
	Items    []Item
}

type Item struct {
	Url     string
	Type    string
	Id      string
	Payload interface{}
}

func NilParser([]byte) ParseResult {
	return ParseResult{}
}
