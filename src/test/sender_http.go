package test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type Result struct {
	RequestNum   uint
	ResponseTime time.Duration
}

func sendRequest(destUrl string, payload []byte, reqNo uint, client *http.Client, results []Result) {
	var req *http.Request
	var err error
	if payload != nil {
		req, err = http.NewRequest("POST", destUrl, bytes.NewBuffer(payload))
	} else {
		req, err = http.NewRequest("GET", destUrl, nil)
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	start := time.Now()
	resp, _ := client.Do(req)
	/*if resPayload, err := ioutil.ReadAll(resp.Body); err == nil {
		fmt.Println("Response payload Len - ", len(resPayload))
	}*/
	io.Copy(ioutil.Discard, resp.Body)
	elapsed := time.Since(start)
	resp.Body.Close()

	results[reqNo].RequestNum = reqNo
	results[reqNo].ResponseTime = elapsed
	fmt.Printf("Transaction %d, Elapsed time - %s\n", reqNo, elapsed)
}

func prepareRequest(numOfTransactions, payloadSize uint) (*http.Client, []byte, []Result, error) {
	var err error
	var payload []byte

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
		},
	}

	if 0 == payloadSize {
		payload = nil
	} else {
		payload = make([]byte, payloadSize)
		for i, _ := range payload {
			payload[i] = 'a'
		}
		fmt.Println("Payload - ", payload, " Len - ", len(payload))
	}
	if err != nil {
		fmt.Println(err)
		return nil, nil, nil, err
	}

	results := make([]Result, numOfTransactions)

	return client, payload, results, err
}
