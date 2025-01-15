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

func WriteJson(w *http.ResponseWriter, data interface{}) {

	(*w).Header().Set("Content-Type", "application/json")
	json.NewEncoder((*w)).Encode(data)
}

func KernelInfo(w http.ResponseWriter, r *http.Request) {

	k := global.GetRunner("kernel")
	data := k.GetData()
	WriteJson(&w, data)
}

func ProcessInfo(w http.ResponseWriter, r *http.Request) {
	k := global.GetRunner("process")
	data := k.GetData()
	WriteJson(&w, data)
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
	mux.Handle("/processes", http.HandlerFunc(ProcessInfo))

	log.Printf("[*] Server started at http://localhost:%d\n", PORT)
	http.Serve(listener, mux)

}
