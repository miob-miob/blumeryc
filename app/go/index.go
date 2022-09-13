package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// ~~~~~~~~~ config ~~~~~~~~~

var blumrycURL = "https://1yaq2zrc91.execute-api.eu-central-1.amazonaws.com/default/blumeryc-downstream-service-dominik-tilp"
var restRequestsDelayMs = 300
var DOWNSTREAM_SERVICE_TIMEOUT_MS = time.Duration(restRequestsDelayMs) * time.Millisecond

// ~~~~~~~~~~~~~~~~~~~~~~~~~~

type ServiceDataResponse struct {
	RequestId string `json:"requestId"`
	Timeout   int    `json:"timeout"`
}

type Result struct {
	isValid    bool
	message    string
	errMessage string
}

func getDownstreamData(timeout time.Duration, out chan<- Result) {

	client := http.Client{
		Timeout: timeout,
	}

	resp, resErr := client.Get(blumrycURL)

	if resErr != nil {
		out <- Result{
			isValid:    false,
			errMessage: "response error: " + resErr.Error(),
		}
		return
	}

	if resp.StatusCode != http.StatusOK {
		out <- Result{
			isValid:    false,
			errMessage: "status not ok",
		}
		return
	}

	bodyRaw, parseRawBodyErr := ioutil.ReadAll(resp.Body)
	if parseRawBodyErr != nil {
		out <- Result{
			isValid:    false,
			errMessage: "body parser error: " + parseRawBodyErr.Error(),
		}
		return
	}

	ServiceDataResponse := ServiceDataResponse{}

	bodyParsingErr := json.Unmarshal(bodyRaw, &ServiceDataResponse)

	if bodyParsingErr != nil {
		out <- Result{
			isValid:    false,
			errMessage: "json body object parser error: " + bodyParsingErr.Error(),
		}
		return
	}

	// validate JSON object schema
	// TODO: check if it is working OK
	if ServiceDataResponse.RequestId == "" || ServiceDataResponse.Timeout == 0 {
		out <- Result{
			isValid:    false,
			errMessage: "invalid JSON schema object error",
		}
		return
	}

	out <- Result{
		isValid: true,
		message: string(bodyRaw),
	}

}

func callDownstreamService(w http.ResponseWriter, req *http.Request) {
	fmt.Println("")
	fmt.Println("~~~~~~>")

	// ----------------------------
	// parse & validate HTTP inputs

	queryParams := req.URL.Query()

	qTimeout, err := strconv.Atoi(queryParams.Get("timeout"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 bad request, query param 'timeout' is not int"))
		return
	}

	if qTimeout <= restRequestsDelayMs {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 bad request, <= " + strconv.Itoa(restRequestsDelayMs)))
		return
	}

	// -----------------------------------
	// do the downstream service api calls

	var qTimeoutMs = time.Duration(qTimeout) * time.Millisecond

	// create channels to do the communications between async API calls
	channel1 := make(chan Result)
	channel2 := make(chan Result)
	channel3 := make(chan Result)

	// this variable is available across gorutine is that proper behavior?
	var someReqAlreadySucceeded = false

	// 1st req
	go getDownstreamData(qTimeoutMs, channel1)

	// fetch 2. API call
	// TODO: bug if 1st req fail in 10sec, 2nd and 3rd will not be exec right after 1st fail
	go func() {
		time.Sleep(DOWNSTREAM_SERVICE_TIMEOUT_MS)
		if someReqAlreadySucceeded {
			return
		}
		getDownstreamData(qTimeoutMs-DOWNSTREAM_SERVICE_TIMEOUT_MS, channel2)
	}()

	// fetch 3. API call
	go func() {
		time.Sleep(DOWNSTREAM_SERVICE_TIMEOUT_MS)
		if someReqAlreadySucceeded {
			return
		}
		getDownstreamData(qTimeoutMs-DOWNSTREAM_SERVICE_TIMEOUT_MS, channel3)
	}()

	for i := 0; i < 3; i++ {

		if someReqAlreadySucceeded {
			break
		}

		var result = Result{}
		var order = 0

		select {
		case res1 := <-channel1:
			result = res1
			order = 1

		case res2 := <-channel2:
			result = res2
			order = 2

		case res3 := <-channel3:
			result = res3
			order = 3
		}

		if !result.isValid {
			fmt.Println("req "+strconv.Itoa(order)+" err: ", result.errMessage)
		} else {
			fmt.Println("req "+strconv.Itoa(order)+" ok : ", result.message)
			someReqAlreadySucceeded = true
			fmt.Fprintf(w, result.message)
			return
		}

	}

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "downstream services is not working properly")
}

func main() {

	fmt.Println("server is running on port 8090")
	http.HandleFunc("/hello", callDownstreamService)

	// TODO; extract port into process envs
	http.ListenAndServe(":8090", nil)
}
