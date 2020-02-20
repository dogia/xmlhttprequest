package xmlhttprequest

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//New ...
func New() *XMLHttpRequest {
	var xhr = &XMLHttpRequest{}
	xhr.request = &http.Request{}
	xhr.request.Header = defaultHeaders()
	xhr.readyState = unsent
	xhr.transport = &http.Transport{}
	xhr.client = &http.Client{
		Transport: xhr.transport,
	}
	xhr.async = false
	xhr.sendFlag = false
	xhr.ResponseText = ""
	xhr.ResponseXML = nil
	xhr.Status = 0
	xhr.StatusText = ""
	xhr.WithCredentials = false
	xhr.eventListener = &eventListener{"", func() {}}
	return xhr
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

//EventListener ...
func (xhr *XMLHttpRequest) EventListener(event string, action func()) {
	xhr.eventListener = &eventListener{event, action}
}

//Open ...
func (xhr *XMLHttpRequest) Open(method, URL string, async bool, user, password string) {
	xhr.abortOpen()

	if isOnMap(forbiddenRequestMethods, method) {
		log.Fatalf("Method %v is not supported by xhr", method)
	}

	sURL, err := url.Parse(URL)
	if err != nil {
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
	if strings.EqualFold(xhr.eventListener.event, "readystatechange") {
		xhr.eventListener.mananger()
	}
	if strings.EqualFold(xhr.eventListener.event, "onprogress") {
		xhr.eventListener.mananger()
	}
}

//Send ...
func (xhr *XMLHttpRequest) Send(data string) {
	if xhr.readyState != opened {
		log.Fatalln("Cannot use send if readyState is diferent of OPENED")
	} else if xhr.sendFlag {
		log.Fatalln("Send has already been called")
	}

	var local bool

	if strings.EqualFold(xhr.request.URL.Scheme, "file") {
		local = true
	}

	if local {
		if !strings.EqualFold(xhr.request.Method, "GET") {
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
				if strings.EqualFold(xhr.eventListener.event, "readystatechange") {
					xhr.eventListener.mananger()
				}
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
			if strings.EqualFold(xhr.eventListener.event, "readystatechange") {
				xhr.eventListener.mananger()
			}
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

		if xhr.GetRequestHeader("Content-Type") == nil {
			xhr.forceHeader("Content-Type", "text/plain;charset=UTF-8")
		}
	} else if !strings.EqualFold(xhr.request.Method, "POST") {
		xhr.forceHeader("Content-Length", "0")
	}
	xhr.errorFlag = false

	var err error
	xhr.response, err = xhr.client.Do(xhr.request)
	if err != nil {
		panic(err)
	}
	xhr.readyState = headersRecived
	if strings.EqualFold(xhr.eventListener.event, "readystatechange") {
		xhr.eventListener.mananger()
	}

	dt, err := ioutil.ReadAll(xhr.response.Body)
	if err != nil {
		panic(err)
	}
	xhr.ResponseText = string(dt)
	xhr.readyState = loading
	if strings.EqualFold(xhr.eventListener.event, "readystatechange") {
		xhr.eventListener.mananger()
	}

	err = xml.Unmarshal(dt, xhr.ResponseXML)
	if err != nil {
		xhr.ResponseXML = nil
	}

	xhr.readyState = done
	if strings.EqualFold(xhr.eventListener.event, "readystatechange") {
		xhr.eventListener.mananger()
	}
	fmt.Println(xhr.response)
	xhr.sendFlag = true
}

//Abort ...
func (xhr *XMLHttpRequest) Abort() {
	if strings.EqualFold(xhr.eventListener.event, "onabort") {
		xhr.eventListener.mananger()
	}

	if xhr.request != nil {
		xhr.transport.CancelRequest(xhr.request)
	}
	New()
}

//ReadyState ...
func (xhr *XMLHttpRequest) ReadyState() uint8 {
	return xhr.readyState
}
