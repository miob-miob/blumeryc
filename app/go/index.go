// Writing a basic HTTP server is easy using the
// `net/http` package.
package main

import (
	"encoding/json"
	"io/ioutil"

	"strconv"

	"fmt"
	"net/http"
	"time"
)

var blumrycURL = "https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp"

type ResponseServiceData struct {
	RequestId string `json:"requestId"`
	Timeout   int    `json:"timeout"`
}

type Result struct {
	isValid    bool
	message    string
	errMessage string
}

func doBRReq(timeout time.Duration, out chan<- Result) {
	output := new(Result)

	client := http.Client{
		Timeout: timeout,
	}

	resp, resErr := client.Get(blumrycURL)

	// defer resp.Body.Close()

	if resErr != nil {
		output.errMessage = "resErr " + resErr.Error()
		output.isValid = false
		out <- *output
		return
	}

	if resp.StatusCode != http.StatusOK {
		output.errMessage = "status not ok"
		output.isValid = false
		out <- *output
		return
	}

	bodyRaw, parseRawBodyErr := ioutil.ReadAll(resp.Body)
	if parseRawBodyErr != nil {
		output.errMessage = "bodyErr" + parseRawBodyErr.Error()
		output.isValid = false
		out <- *output
		return
	}

	responseServiceData := ResponseServiceData{}

	bodyParsingErr := json.Unmarshal(bodyRaw, &responseServiceData)

	if bodyParsingErr != nil {
		output.errMessage = "bodyParsingErr: " + bodyParsingErr.Error()
		output.isValid = false
		out <- *output
		return
	}

	// validate JSON schema
	if responseServiceData.RequestId == "" || responseServiceData.Timeout == 0 {
		output.errMessage = "invalid JSON schema"
		output.isValid = false
		out <- *output
		return
	}

	output.isValid = true
	output.message = string(bodyRaw)

	out <- *output

}

func callDownstreamService(w http.ResponseWriter, req *http.Request) {
	fmt.Println("------->")

	// ----------------------
	// parse & validate inputs

	queryParams := req.URL.Query()

	qTimeout, err := strconv.Atoi(queryParams.Get("timeout"))

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

	// -----------------------------------
	// do the downstream service api calls

	var userTimeout = time.Duration(qTimeout) * time.Millisecond
	var initRequestTimeout = userTimeout
	// TODO change to 300
	var restRequestDelay = time.Duration(300) * time.Millisecond
	var restRequestTimeout = userTimeout - restRequestDelay

	c1 := make(chan Result)
	c2 := make(chan Result)
	c3 := make(chan Result)

	// TODO: merge with sent variable
	var someReqAlreadySucceeded = false

	// 1st req
	go doBRReq(initRequestTimeout, c1)

	go func() {
		time.Sleep(restRequestDelay)
		fmt.Println("is 1st done", someReqAlreadySucceeded)
		if someReqAlreadySucceeded {
			return
		}
		// 2nd req
		doBRReq(restRequestTimeout, c2)
	}()

	go func() {
		time.Sleep(restRequestDelay)
		fmt.Println("is 1st done", someReqAlreadySucceeded)
		if someReqAlreadySucceeded {
			return
		}
		// 3rd req
		doBRReq(restRequestTimeout, c3)
	}()

	count := 0

	for !someReqAlreadySucceeded {

		// iterate 1 to 3 times in max scenario
		if count == 3 {
			break
		}
		// for i := 0; i < 3; i++ {
		// Await both of these values
		// simultaneously, printing each one as it arrives.
		select {
		case res1 := <-c1:
			count++

			if !res1.isValid {
				fmt.Println("req 1 err: ", res1.errMessage)
			} else {
				fmt.Println("req 1 ok : ", res1.message)
				someReqAlreadySucceeded = true
				fmt.Fprintf(w, res1.message)
				// will return break those channels?
				// return
			}

		case res2 := <-c2:
			count++

			if !res2.isValid {
				fmt.Println("req 2 err: ", res2.errMessage)
			} else {
				fmt.Println("req 2 ok : ", res2.message)
				someReqAlreadySucceeded = true
				fmt.Fprintf(w, res2.message)
				// return
			}

		case res3 := <-c3:
			count++

			if !res3.isValid {
				fmt.Println("req 3 err: ", res3.errMessage)
			} else {
				fmt.Println("req 3 ok : ", res3.message)
				someReqAlreadySucceeded = true
				fmt.Fprintf(w, res3.message)
				// return
			}
		}
	}

	// fmt.Println("someReqAlreadySucceeded: ", someReqAlreadySucceeded)

	// TODO: is 500 proper error status code?
	if !someReqAlreadySucceeded {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "all 3 req failed")
	}

	fmt.Println("<------- req ended")

}

func main() {

	http.HandleFunc("/hello", callDownstreamService)

	// TODO; extract port into process envs
	http.ListenAndServe(":8090", nil)
}
