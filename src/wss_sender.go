package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"time"
)

type Result struct {
	RequestNum   uint
	StartTime    time.Time
	ResponseTime time.Duration
}

func wss_receiver(c *websocket.Conn, respTime chan<- time.Duration, quit <-chan struct{}, result chan Result) {
	defer c.Close()
	for {
		select {
		case rst := <-result:
			_, _, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}
			rst.ResponseTime = time.Since(rst.StartTime)
			//fmt.Printf("recv: %s\n", message)
			fmt.Printf("Transaction %d, Elapsed time - %s\n", rst.RequestNum, rst.ResponseTime)
			respTime <- rst.ResponseTime

		case <-quit:
			return
		}
	}
}

func setUp(noOfTrans, payloadSize uint, destUrl string) (*websocket.Conn, []byte, error) {
	// Setup HTTPS client
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	wss_dialer := &websocket.Dialer{
		TLSClientConfig: tlsConfig,
	}

	c, _, err := wss_dialer.Dial(destUrl, nil)
	if err != nil {
		fmt.Println("dial:", err)
	}

	payload := make([]byte, payloadSize)
	for i, _ := range payload {
		payload[i] = 'a'
	}
	fmt.Println("Payload - ", payload, " Len - ", len(payload))

	return c, payload, err
}

func main() {
	var noOfTrans, payloadSize uint
	var destUrl string
	var totalResTime, avgResTime int64
	var res time.Duration

	flag.UintVar(&noOfTrans, "num-of-transactions", 10, "Number of Transactions/Requests")
	flag.UintVar(&payloadSize, "payload-size", 10, "Payload Size to be sent in each request in bytes")
	flag.StringVar(&destUrl, "dest-url", "wss://localhost:6666/", "URL of the destination/relector server")
	flag.Parse()

	fmt.Println("Number of Transactions/Requests - ", noOfTrans)
	fmt.Println("Payload Size in bytes - ", payloadSize)
	fmt.Println("URL of the relector server", destUrl)

	conn, payload, err := setUp(noOfTrans, payloadSize, destUrl)
	defer conn.Close()
	if err != nil {
		fmt.Println("conn: ", err)
		return
	}

	// channel used to shutdown the go routines
	quit := make(chan struct{})
	respTime := make(chan time.Duration)
	result := make(chan Result, 1)

	go wss_receiver(conn, respTime, quit, result)

	var i uint
	for i = 0; i < noOfTrans; i++ {
		err := conn.WriteMessage(websocket.TextMessage, payload)
		if err != nil {
			fmt.Println("write:", err)
			return
		} else {
			rst := new(Result)
			rst.RequestNum = i
			rst.StartTime = time.Now()
			result <- *rst
			res = <-respTime
			totalResTime += res.Nanoseconds()
		}
	}
	close(quit)
	avgResTime = totalResTime / int64(noOfTrans)
	fmt.Printf("\nAverage Response Time - %s\n\n", time.Duration(avgResTime))
}
