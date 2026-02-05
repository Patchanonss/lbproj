package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)
var mutex sync.Mutex
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Please specify a port (e.g., go run main.go 8081)")
		return
	}
	port := os.Args[1]

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 1. ล็อคประตู! ห้ามใครแซงคิว (รับได้ทีละ request)
		mutex.Lock()
		defer mutex.Unlock()

		// 2. แกล้งทำงานหนัก 100 ms
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("Received request from %s\n", r.RemoteAddr)
		fmt.Fprintf(w, "Hello from Backend Server! Running on Port: %s\n", port)
	})

	fmt.Printf("Server is starting on port %s...\n", port)
	
	// ListenAndServe จะรันค้างไว้ ไม่จบการทำงานจนกว่าจะ Error
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}