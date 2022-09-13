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

func doBRReq(timeout time.Duration, out chan<- Result) {

	client := http.Client{
		Timeout: timeout,
	}

	resp, resErr := client.Get(blumrycURL)

	if resErr != nil {
		out <- Result{
			isValid:    false,
			errMessage: "resErr " + resErr.Error(),
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
			errMessage: "bodyErr" + parseRawBodyErr.Error(),
		}
		return
	}

	ServiceDataResponse := ServiceDataResponse{}

	bodyParsingErr := json.Unmarshal(bodyRaw, &ServiceDataResponse)

	if bodyParsingErr != nil {
		out <- Result{
			isValid:    false,
			errMessage: "bodyParsingErr: " + bodyParsingErr.Error(),
		}
		return
	}

	// validate JSON object schema
	// TODO: check if it is working OK
	if ServiceDataResponse.RequestId == "" || ServiceDataResponse.Timeout == 0 {
		out <- Result{
			isValid:    false,
			errMessage: "invalid JSON schema",
		}
		return
	}

	out <- Result{
		isValid: true,
		message: string(bodyRaw),
	}

}

func callDownstreamService(w http.ResponseWriter, req *http.Request) {
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
	c1 := make(chan Result)
	c2 := make(chan Result)
	c3 := make(chan Result)

	// this variable is available across gorutine is that proper behavior?
	var someReqAlreadySucceeded = false

	// 1st req
	go doBRReq(qTimeoutMs, c1)

	// fetch 2. API call
	go func() {
		time.Sleep(DOWNSTREAM_SERVICE_TIMEOUT_MS)
		if someReqAlreadySucceeded {
			return
		}
		doBRReq(qTimeoutMs-DOWNSTREAM_SERVICE_TIMEOUT_MS, c2)
	}()

	// fetch 3. API call
	go func() {
		time.Sleep(DOWNSTREAM_SERVICE_TIMEOUT_MS)
		if someReqAlreadySucceeded {
			return
		}
		doBRReq(qTimeoutMs-DOWNSTREAM_SERVICE_TIMEOUT_MS, c3)
	}()

	count := 0

	for i := 0; i < 3; i++ {

		if someReqAlreadySucceeded {
			break
		}

		select {
		case res1 := <-c1:
			if !res1.isValid {
				fmt.Println("req 1 err: ", res1.errMessage)
			} else {
				fmt.Println("req 1 ok : ", res1.message)
				someReqAlreadySucceeded = true
				fmt.Fprintf(w, res1.message)
			}

		case res2 := <-c2:
			if !res2.isValid {
				fmt.Println("req 2 err: ", res2.errMessage)
			} else {
				fmt.Println("req 2 ok : ", res2.message)
				someReqAlreadySucceeded = true
				fmt.Fprintf(w, res2.message)
			}

		case res3 := <-c3:
			if !res3.isValid {
				fmt.Println("req 3 err: ", res3.errMessage)
			} else {
				fmt.Println("req 3 ok : ", res3.message)
				someReqAlreadySucceeded = true
				fmt.Fprintf(w, res3.message)
			}
		}

		count++
	}

	if !someReqAlreadySucceeded {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "downstream services is not working properly")
	}

	fmt.Println("<~~~~~~~ req ended")
	fmt.Println("")
}

func main() {

	fmt.Println("server is running on port 8090")
	http.HandleFunc("/hello", callDownstreamService)

	// TODO; extract port into process envs
	http.ListenAndServe(":8090", nil)
}
