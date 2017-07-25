package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{} // use default options

func echoWSSHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}
		fmt.Printf("recv: %s\n", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			fmt.Println("write:", err)
			break
		}
	}
}

func main() {
	main_wg := &sync.WaitGroup{}
	main_wg.Add(1)
	http.HandleFunc("/", echoWSSHandler)
	//fmt.Fatal(http.ListenAndServe("localhost:6666", nil))
	fmt.Println(http.ListenAndServeTLS("localhost:6666", "cert.pem", "key.pem", nil))
	main_wg.Wait()
}
