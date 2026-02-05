package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil" // พระเอกของเรา (Reverse Proxy Utility)
	"net/url"
	"sync"
)

// --- Step 1: The Setup (เตรียมสมุดรายชื่อ) ---
var (
	// รายชื่อ Backend Server ที่เราเปิดทิ้งไว้ (8081, 8082, 8083)
	servers = []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	
	currentIndex = 0          // ตัวนับคิว (Counter)
	mutex        sync.Mutex   // กุญแจล็อค (กันแย่งกันนับเลข)
)
func main() {
	// --- Step 2: Handle (เปิดหน้าร้านรับแขก) ---
	// บอกว่าถ้ามีใครเข้า "/" ให้มาเรียกฟังก์ชัน loadBalancerHandler นะ
	http.HandleFunc("/", loadBalancerHandler)

	fmt.Println("Load Balancer started at :8080")
	
	// เริ่ม Run Server ที่ Port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func loadBalancerHandler(w http.ResponseWriter, r *http.Request) {
	// --- Step 3: Select (เลือกเหยื่อด้วย Round Robin) ---
	// เรียกฟังก์ชันเพื่อขอ URL ของ Server ตัวถัดไป
	targetServer := getNextServer()

	// Parse URL จาก string ให้เป็น Object ที่ Go เข้าใจ
	url, _ := url.Parse(targetServer)

	// --- Step 4 & 5: Director & Proxy (เตรียมส่งต่อ) ---
	// สร้าง Reverse Proxy ชี้ไปที่เป้าหมาย (targetServer)
	// NewSingleHostReverseProxy คือฟังก์ชันวิเศษที่จัดการเรื่องแก้ Header/URL ให้เอง
	proxy := httputil.NewSingleHostReverseProxy(url)

	// Custom: แอบแก้ Header นิดหน่อย ให้ Backend รู้ว่า request นี้ผ่าน LB มานะ
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	fmt.Printf("Redirecting request to: %s\n", targetServer)

	// --- Step 6: Response (ส่งออกไปและรอรับของกลับ) ---
	// ServeHTTP จะทำการ "ยิง Request" ไปหา Backend
	// และเมื่อ Backend ตอบกลับมา มันจะเขียนใส่ w (ResponseWriter) ส่งคืน User ให้เองอัตโนมัติ
	proxy.ServeHTTP(w, r)
}

// ฟังก์ชันช่วยคำนวณ Round Robin (Step 3 ขยายความ)
func getNextServer() string {
	// ล็อคประตูห้องนับเลข ห้ามใครเข้ามายุ่ง
	mutex.Lock()
	defer mutex.Unlock() // ทำงานเสร็จแล้วค่อยปลดล็อค (defer คือทำตอนจบฟังก์ชัน)

	// หยิบ server ตามคิวปัจจุบัน
	server := servers[currentIndex]

	// ขยับตัวนับไปช่องถัดไป
	currentIndex++
	
	// ถ้าตัวนับเกินจำนวน server (เช่น 3) ให้วนกลับไปเป็น 0 ใหม่ (Modulo Logic)
	if currentIndex >= len(servers) {
		currentIndex = 0
	}

	return server
}