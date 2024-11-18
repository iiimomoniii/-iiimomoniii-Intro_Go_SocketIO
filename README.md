# WebSocket Server and Client Example in Go

## 1. Setting up the Go Project
go mod init socketio
go install golang.org/x/net/websocket@latest
go get golang.org/x/net/websocket


## 2. WebSocket Client in Chrome (JavaScript)
chrome
## 2.1 Client for General WebSocket Connection
let socket = new WebSocket("ws://localhost:3000/ws");
socket.onmessage = (event) => {
    console.log("Received from the server:", event.data);
};
socket.send("hello from client");
## 2.2 Client for Orderbook Feed
let socket = new WebSocket("ws://localhost:3000/orderbookfeed");
socket.onmessage = (event) => {
    console.log("Received orderbook data:", event.data);
};

## 3. run
go run main.go