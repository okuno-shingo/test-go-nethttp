package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const MaxConcurrentRequests = 5

var semaphore = make(chan struct{}, MaxConcurrentRequests)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		semaphore <- struct{}{}
		defer func() { <-semaphore }()

		if rand.Intn(10) == 0 {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Hello from server at: %s\n", r.Host)

		if rand.Intn(20) == 0 {
			time.Sleep(5000 * time.Millisecond)
		} else {
			time.Sleep(10 * time.Millisecond)
		}
	})

	fmt.Println("Server is running...")
	http.ListenAndServe(":8080", nil)
}
