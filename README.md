# xmlhttprequest

This is a library to simulate xmlhttprequest javascript´s object on golang. To import it put in your OS console. 
->go get github.com/dogia/xmlhttprequest

<h3>How to use:</h3>

package main

import xhr "github.com/dogia/xmlhttprequest"

func main(){

  xhr := &xhr.XMLHttpRequest{}
  
  xhr.Open(method, URL #string, async #bool, user, password #string)
  
  xhr.EventListener("onload", func() {
		fmt.Println("READY")
	})
  
  xhr.Send(bodyData #string)
  
}

Functions

xhr.EventListener(event #string, mananger #func())

xhr.SetHeader(Name, Value #string)

xhr.Abort()

xhr.ReadyState() //returns uint8 readyState

