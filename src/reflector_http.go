package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func echoHTTPHandler(w http.ResponseWriter, r *http.Request) {
	if payload, err := ioutil.ReadAll(r.Body); err == nil {
		fmt.Println(payload)
		w.Write(payload)
	} else {
		fmt.Println(err)
	}
}

func main() {
	main_wg := &sync.WaitGroup{}
	main_wg.Add(1)
	http.HandleFunc("/", echoHTTPHandler)
	http.ListenAndServe(":8888", nil)
	main_wg.Wait()
}
