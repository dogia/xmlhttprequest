package xmlhttprequest

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

type eventListener struct {
	event    string
	mananger func()
}

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

//SetHeader ...
func (xhr *XMLHttpRequest) SetHeader(key string, value string) {
	if xhr.readyState != opened {
		if strings.EqualFold(xhr.eventListener.event, "onerror") {
			xhr.eventListener.mananger()
		}
		log.Fatalln("Cannot set headers if readyState is diferent of OPENED")
		return
	}

	if !isOnMap(forbiddenRequestHeaders, key) {
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
	} else {
		if debug {
			log.Printf("Header %v is not supported by xhr \n", key)
		}
	}
}

//GetResponseHeader ...
func (xhr *XMLHttpRequest) GetResponseHeader(header string) []string {
	if xhr.readyState > opened && xhr.response != nil && xhr.response.Header != nil && !xhr.errorFlag {
		for k, v := range xhr.response.Header {
			if k == header {
				return v
			}
		}
	}
	return nil
}

//GetAllResponseHeader ...
func (xhr *XMLHttpRequest) GetAllResponseHeader(header string) map[string][]string {
	if xhr.readyState > headersRecived && !xhr.errorFlag {
		retHeaders := make(map[string][]string)
		for k, v := range xhr.response.Header {
			if !strings.EqualFold(k, "set-cookie") && !strings.EqualFold(k, "set-cookie2") {
				retHeaders[k] = v
			}
		}
		return retHeaders
	}
	return nil
}

//GetRequestHeader ...
func (xhr *XMLHttpRequest) GetRequestHeader(header string) []string {
	if xhr.readyState > unsent {
		for k, v := range xhr.request.Header {
			if k == header {
				return v
			}
		}
	}
	return nil
}

func defaultHeaders() map[string][]string {
	headers := make(map[string][]string)
	headers["User-Agent"] = []string{"go-XMLHttpRequest"}
	headers["Accept"] = []string{"*/*"}

	return headers
}

//New ...
func (xhr *XMLHttpRequest) New() {
	xhr.eventListener = &eventListener{"", func() {}}
	xhr.request = &http.Request{}
	xhr.request.Header = defaultHeaders()
	xhr.readyState = 0
	xhr.transport = &http.Transport{}
	xhr.client = &http.Client{
		Transport: xhr.transport,
	}
	xhr.async = false
	xhr.sendFlag = false
	xhr.readyState = 0
	xhr.ResponseText = ""
	xhr.ResponseXML = nil
	xhr.Status = 0
	xhr.StatusText = ""
	xhr.WithCredentials = false
}

//EventListener ...
func (xhr *XMLHttpRequest) EventListener(event string, action func()) {
	xhr.eventListener = &eventListener{event, action}
}

func isAllowedHTTPHeader(r *http.Request) {

}

func isAllowedHTTPMethod(r *http.Request) {

}

//Open ...
func (xhr *XMLHttpRequest) Open(method, URL string, async bool, user, password string) {
	xhr.Abort()

	if isOnMap(forbiddenRequestMethods, method) {
		if strings.EqualFold(xhr.eventListener.event, "onabort") {
			xhr.eventListener.mananger()
		}
		log.Fatalf("Method %v is not supported by xhr", method)
	}

	sURL, err := url.Parse(URL)
	if err != nil {
		if strings.EqualFold(xhr.eventListener.event, "onabort") {
			xhr.eventListener.mananger()
		}
		log.Fatalln(err)
	}

	xhr.request.Method = method
	xhr.request.URL = sURL
	xhr.async = async

	if !strings.EqualFold(user, "") {
		auth := base64.StdEncoding.EncodeToString([]byte(user + ":" + password))
		xhr.SetHeader("Authorization", "Basic "+auth)
	}

	xhr.readyState = opened
	if strings.EqualFold(xhr.eventListener.event, "onprogress") {
		xhr.eventListener.mananger()
	}
}

//Send ...
func (xhr *XMLHttpRequest) Send(data string) {
	if xhr.readyState != opened {
		if strings.EqualFold(xhr.eventListener.event, "onabort") {
			xhr.eventListener.mananger()
		}
		log.Fatalln("Cannot use send if readyState is diferent of OPENED")
	} else if xhr.sendFlag {
		if strings.EqualFold(xhr.eventListener.event, "onabort") {
			xhr.eventListener.mananger()
		}
		log.Fatalln("Send has already been called")
	}

	var local bool

	if strings.EqualFold(xhr.request.URL.Scheme, "file") {
		local = true
	}

	if local {
		if !strings.EqualFold(xhr.request.Method, "GET") {
			if strings.EqualFold(xhr.eventListener.event, "onabort") {
				xhr.eventListener.mananger()
			}
			log.Fatalln("Only method get supported.")
		}

		if xhr.async {
			go func() {
				data, err := ioutil.ReadFile(xhr.request.URL.Host)
				if err != nil {
					panic(err)
				}
				xhr.ResponseText = string(data)
				xhr.readyState = done
				xhr.Status = 200
				xhr.StatusText = "OK"
				if strings.EqualFold(xhr.eventListener.event, "onready") || strings.EqualFold(xhr.eventListener.event, "onreadystatechange") {
					xhr.eventListener.mananger()
				}
			}()
		} else {
			data, err := ioutil.ReadFile(xhr.request.URL.Host)
			if err != nil {
				panic(err)
			}
			xhr.ResponseText = string(data)
			xhr.readyState = done
			xhr.Status = 200
			xhr.StatusText = "OK"
			if strings.EqualFold(xhr.eventListener.event, "onready") || strings.EqualFold(xhr.eventListener.event, "onreadystatechange") {
				xhr.eventListener.mananger()
			}
		}
	}

	if strings.EqualFold(xhr.request.Method, "GET") || strings.EqualFold(xhr.request.Method, "HEAD") {
		data = ""
	} else if !strings.EqualFold(data, "") {
		length := strconv.Itoa(len(data))
		xhr.forceHeader("Content-Length", length)
		//TODO
	}

	xhr.sendFlag = true
}

//Abort ...
func (xhr *XMLHttpRequest) Abort() {
	if xhr.request != nil {
		xhr.transport.CancelRequest(xhr.request)
	}

	xhr.New()
	if strings.EqualFold(xhr.eventListener.event, "onabort") {
		xhr.eventListener.mananger()
	}
}

//ReadyState ...
func (xhr *XMLHttpRequest) ReadyState() uint8 {
	return xhr.readyState
}
