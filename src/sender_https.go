package main

import (
	"bytes"
	"crypto/tls"
	"flag"
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

	// Setup HTTPS client
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			TLSClientConfig:     tlsConfig,
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

func main() {
	var noOfTrans, payloadSize uint
	var destUrl string

	flag.UintVar(&noOfTrans, "num-of-transactions", 10, "Number of Transactions/Requests")
	flag.UintVar(&payloadSize, "payload-size", 10, "Payload Size to be sent in each request in bytes")
	flag.StringVar(&destUrl, "dest-url", "https://localhost:9999/", "URL of the destination/relector server")
	flag.Parse()

	fmt.Println("Number of Transactions/Requests - ", noOfTrans)
	fmt.Println("Payload Size in bytes - ", payloadSize)
	fmt.Println("URL of the relector server", destUrl)

	client, payload, results, _ := prepareRequest(noOfTrans, payloadSize)

	var totalResTime, avgResTime, firstConnResTime int64
	var i uint

	for i = 0; i < noOfTrans; i++ {

		sendRequest(destUrl, payload, i, client, results)

		if i == 0 {
			firstConnResTime = results[i].ResponseTime.Nanoseconds()
		} else {
			totalResTime += results[i].ResponseTime.Nanoseconds()
		}
		fmt.Printf("Result %d, Elapsed time - %s\n", results[i].RequestNum, results[i].ResponseTime)
	}
	avgResTime = totalResTime / int64(noOfTrans)
	fmt.Println("\nFirst Connection Response Time - ", time.Duration(firstConnResTime))
	fmt.Printf("\nAverage Connection Reuse Response Time - %s\n\n", time.Duration(avgResTime))
}
