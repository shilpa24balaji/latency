package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

func echoHTTPSHandler(w http.ResponseWriter, r *http.Request) {
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
	http.HandleFunc("/", echoHTTPSHandler)
	http.ListenAndServeTLS("127.0.0.1:9999", "cert.pem", "key.pem", nil)
	main_wg.Wait()
}
