package request

import (
	"net/http"

	"github.com/nmeji/tstr/response"
)

type HttpHeader map[string]string

type HttpBody []byte

type MimeType int

const (
	MIMETYPE_JSON = 1 + iota
	MIMETYPE_XML
	MIMETYPE_HTML
	MIMETYPE_PLAIN_TEXT
)

type Request struct {
	Header HttpHeader
	Body   HttpBody
}

func (r Request) Send(url string) {

}

func New() *Request {
	return &Request{}
}

func Get(url string) (*response.Checker, error) {
	httpresp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return response.NewChecker(httpresp), nil
}
