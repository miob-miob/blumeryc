// Writing a basic HTTP server is easy using the
// `net/http` package.
package main

import (
	"io/ioutil"
	"log"

	"strconv"

	"fmt"
	"net/http"
	"time"
)

var blumrycURL = "https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp"

var boolTrue bool = true
var boolFalse bool = false

// isValid boolean
func doBRReq(timeout time.Duration, out chan<- string) (bool, string) {

	client := http.Client{
		Timeout: timeout,
	}

	fmt.Println("1")

	resp, resErr := client.Get(blumrycURL)

	// defer resp.Body.Close()

	if resErr != nil {
		fmt.Println("resErr ahoj")
		fmt.Println(resErr)
		return false, ""
	}
	fmt.Println("2")

	if resp.StatusCode != http.StatusOK {
		return false, ""
	}

	fmt.Println("3")
	// TODO: validate, that body includes data

	body, bodyErr := ioutil.ReadAll(resp.Body)

	if bodyErr != nil {
		fmt.Println("bodyErr ahoj")
		fmt.Println(resErr)
		// log.Fatalln(bodyErr)
		return false, ""
	}

	log.Printf("error")
	// Exception will be handled...

	return true, string(body)

}

// A fundamental concept in `net/http` servers is
// *handlers*. A handler is an object implementing the
// `http.Handler` interface. A common way to write
// a handler is by using the `http.HandlerFunc` adapter
// on functions with the appropriate signature.
func hello(w http.ResponseWriter, req *http.Request) {

	x := req.URL.Query()

	qTimeout, err := strconv.Atoi(x.Get("timeout"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 bad request, timeout qParam is not int"))
		return
	}

	if qTimeout <= 300 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 bad request, <= 300"))
		return
	}

	// -------------

	var userTimeout = time.Duration(qTimeout) * time.Millisecond
	var initRequestTimeout = userTimeout
	// TODO: implement
	var restRequestDelay = time.Duration(300) * time.Millisecond
	var restRequestTimeout = userTimeout - restRequestDelay

	fmt.Println(
		restRequestDelay,
		restRequestTimeout,
		qTimeout,
	)

	c1 := make(chan string, 1)
	// c2 := make(chan string)
	// c2 := make(chan string)

	isValid, body := doBRReq(initRequestTimeout, c1)

	fmt.Println("--------------")
	fmt.Println(isValid)
	fmt.Println(body)

	fmt.Fprintf(w, body)
}

func main() {

	// We register our handlers on server routes using the
	// `http.HandleFunc` convenience function. It sets up
	// the *default router* in the `net/http` package and
	// takes a function as an argument.
	http.HandleFunc("/hello", hello)
	// http.HandleFunc("/headers", headers)

	// Finally, we call the `ListenAndServe` with the port
	// and a handler. `nil` tells it to use the default
	// router we've just set up.
	http.ListenAndServe(":8090", nil)
}
