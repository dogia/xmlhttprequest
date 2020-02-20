# xmlhttprequest

This is a library to simulate xmlhttprequest javascript´s object on golang. To import it put in your OS console. 
->go get github.com/dogia/xmlhttprequest

<h3>How to use:</h3>

package main

import xhr "github.com/dogia/xmlhttprequest"

	func main(){

		xhr := xhr.New()

		xhr.EventListener("readystatechange", func() {
			fmt.Println("Ready state: " + strconv.Itoa(int(xhr.ReadyState())))
		})
		xhr.Open("POST", "http://localhost", false, "", "")
		xhr.Send("data=abcde")

		fmt.Println(xhr.ResponseText)

Functions

xhr.EventListener(event #string, mananger #func())

xhr.SetHeader(Name, Value #string)

xhr.Abort()

xhr.ReadyState() //returns uint8 readyState

