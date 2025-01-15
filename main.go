package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

const PORT int = 5555

func KernelInfo(w http.ResponseWriter, r *http.Request) {

	k := global.getRunner("kernel")
	data := k.GetData()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func main() {

	mux := http.NewServeMux()
	go StartRunners()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", PORT))

	if err != nil {
		log.Fatalf("[!] Couldn't connect to %d", PORT)
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second)
		w.Write([]byte("Hello"))

	})

	mux.Handle("/kernel", http.HandlerFunc(KernelInfo))

	log.Printf("[*] Server started at http://localhost:%d\n", PORT)
	http.Serve(listener, mux)

}
