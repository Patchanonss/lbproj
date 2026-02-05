package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func main() {
	// เป้าหมาย: จะยิงไปที่ไหน?
	// ลองเปลี่ยนเป็น "http://localhost:8081" เพื่อเทสตัวเดียว
	// หรือ "http://localhost:8080" เพื่อเทสผ่าน LB
	targetURL := "http://localhost:8080" 
	
	totalRequests := 100
	var wg sync.WaitGroup

	start := time.Now() // จับเวลาเริ่ม

	fmt.Printf("Attacking %s with %d requests...\n", targetURL, totalRequests)

	for i := 0; i < totalRequests; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := http.Get(targetURL)
			if err != nil {
				fmt.Printf("Request %d failed: %v\n", id, err)
			}
		}(i)
	}

	wg.Wait() // รอจนกว่าลูกกระสุนทั้ง 100 นัดจะทำงานเสร็จ
	duration := time.Since(start)

	fmt.Println("\n----------------------------------")
	fmt.Printf("Total Time Taken: %v\n", duration)
	fmt.Println("----------------------------------")
}