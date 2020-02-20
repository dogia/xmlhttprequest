package xmlhttprequest

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"strings"
)

//XMLHttpRequest it's like XMLHttpRequest javascript object
type XMLHttpRequest struct {
	async           bool
	transport       *http.Transport
	client          *http.Client
	request         *http.Request
	response        *http.Response
	readyState      uint8
	eventListener   *eventListener
	ResponseText    string
	ResponseXML     interface{}
	Status          uint16
	StatusText      string
	WithCredentials bool
	sendFlag        bool
	errorFlag       bool
}

type eventListener struct {
	event    string
	mananger func()
}

var forbiddenRequestHeaders = []string{
	"accept-charset",
	"accept-encoding",
	"access-control-request-headers",
	"access-control-request-method",
	"connection",
	"content-length",
	"content-transfer-encoding",
	"cookie",
	"cookie2",
	"date",
	"expect",
	"host",
	"keep-alive",
	"origin",
	"referer",
	"te",
	"trailer",
	"transfer-encoding",
	"upgrade",
	"via",
}

var forbiddenRequestMethods = []string{
	"trace",
	"track",
	"connect",
}

func isOnMap(haystack, needle interface{}) bool {
	switch haystack.(type) {
	case []string:
		for _, v := range haystack.([]string) {
			if strings.EqualFold(v, needle.(string)) {
				return true
			}
		}

	default:
		switch reflect.TypeOf(haystack).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(haystack)

			for i := 0; i < s.Len(); i++ {
				if s.Index(i) == reflect.ValueOf(needle) {
					return true
				}
			}
		default:
			err := errors.New("Type doesn´t supported by isOnMap")
			panic(err)
		}
	}
	return false
}
func (xhr *XMLHttpRequest) forceHeader(key string, value string) {
	reqHeader := make(map[string][]string)
	reqHeader[key] = []string{value}
	for k, v := range xhr.request.Header {
		reqHeader[k] = v
	}
	reqHeader[key] = []string{value}
	xhr.request.Header = reqHeader
	if debug {
		log.Printf("Header added %v with value %v\n", key, value)
	}
}

func defaultHeaders() map[string][]string {
	headers := make(map[string][]string)
	headers["User-Agent"] = []string{"go-XMLHttpRequest"}
	headers["Accept"] = []string{"*/*"}

	return headers
}

func (xhr *XMLHttpRequest) abortOpen() {
	if strings.EqualFold(xhr.eventListener.event, "onabort") {
		xhr.eventListener.mananger()
	}

	if xhr.request != nil {
		xhr.transport.CancelRequest(xhr.request)
	}
}
