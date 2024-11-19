package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

// 1.
// โครงสร้าง Server
type Server struct {
	//ประกาศ Server เช็ค connection นั้นยัง active อยู่หรือไม่
	conns map[*websocket.Conn]bool
}

// 2.
// สร้าง server ใหม่
func NewServer() *Server {
	return &Server{
		//สร้าง instance ใหม่ของ Server และคืนค่ากลับ
		conns: make(map[*websocket.Conn]bool),
	}
}

// 3.
// จัดการ WebSocket ที่เชื่อมต่อเข้ามาใหม่
func (s *Server) handlews(ws *websocket.Conn) {
	fmt.Println("new incomming connection from client:", ws.RemoteAddr())
	//เพิ่ม connection ลงใน conns
	s.conns[ws] = true
	//readLoop เพื่อรออ่านข้อความจาก client
	s.readLoop(ws)
}

// กรณีที่ ทดสอบส่ง ที่มีข้อมูล orderbook (payload) ไปยัง client ทุกๆ 2 วินาที
func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
	fmt.Println("new incomming connention from client to orderbook feed:", ws.RemoteAddr())

	for {
		//time.Now().UnixNano() เพื่อแสดง timestamp
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		//ทำหน้าที่ส่งข้อมูลผ่าน WebSocket connection ไปยัง client
		ws.Write([]byte(payload))
		//ส่งข้อมูลใน ทุกๆ 2 วินาที
		time.Sleep(time.Second * 2)
	}
}

// progress of process
func (s *Server) handleProgress(ws *websocket.Conn) {
	fmt.Println("New connection for progress tracking:", ws.RemoteAddr())

	total := 10 // Total steps in the process
	for i := 1; i <= total; i++ {
		// Simulate processing time for each step
		time.Sleep(time.Second)

		// Calculate progress percentage
		progress := float64(i) / float64(total) * 100

		// Create a payload containing progress details
		payload := fmt.Sprintf("{\"step\": %d, \"total\": %d, \"progress\": %.2f}", i, total, progress)

		// Send the progress payload to the client
		if _, err := ws.Write([]byte(payload)); err != nil {
			fmt.Println("Write error:", err)
			return
		}
		time.Sleep(time.Second * 2)
	}

	// Notify the client that processing is complete
	completionMessage := "{\"status\": \"complete\", \"message\": \"Processing complete!\"}"
	if _, err := ws.Write([]byte(completionMessage)); err != nil {
		fmt.Println("Write error on completion:", err)
	}
}

// 4.
// อ่านข้อความจาก client
func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error:", err)
			continue
		}
		//อ่านข้อความจาก client (ในลูป)
		msg := buf[:n]
		fmt.Println(string(msg))
		//ส่งต่อข้อความไปยังฟังก์ชัน broadcast
		s.broadcast(msg)
	}
}

// 5.
// รับข้อความจาก function readLoop
func (s *Server) broadcast(b []byte) {
	//ส่งข้อความที่รับมาจาก readLoop ผ่านตัวแปร b ไปยัง WebSocket connections ทั้งหมดใน conns ผ่านการวนลูป
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("write error:", err)
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	//กำหนด route ที่อ่านข้อความจาก client
	http.Handle("/ws", websocket.Handler(server.handlews))
	//กำหนด route ที่ server จะส่งข้อความกลับในทุกๆ 2 วินาที
	http.Handle("/orderbookfeed", websocket.Handler(server.handleWSOrderbook))
	//
	http.Handle("/progress", websocket.Handler(server.handleProgress))
	http.ListenAndServe(":3000", nil)

}
