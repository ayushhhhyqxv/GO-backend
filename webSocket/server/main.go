package main

// server side

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Request struct {
	Action string `json:"action"`
}

type Response struct {
	Reply string `json:"reply"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin:  func(r *http.Request) bool {return true},  // Currently allows all request , usually used for CORS !
}

func wsHandler(w http.ResponseWriter, r *http.Request){
	conn,err := upgrader.Upgrade(w,r,nil)
	if err!=nil{
		log.Printf("Upgrade Error %v",err)
		return 
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(30*time.Second))
	conn.SetWriteDeadline(time.Now().Add(10*time.Second))

	// conn.SetPongHandler(func(data string) error {
	// 	log.Println("Pong Received")
	// 	conn.SetReadDeadline(time.Now().Add(30*time.Second))
	// 	return nil
	// }) 

	// Already this work is done by Gorilla ! 

	for {
		var req Request 

		err:= conn.ReadJSON(&req)
		if err!=nil {

			// Check for connection error's too

			if websocket.IsCloseError(err,websocket.CloseNormalClosure){
				log.Println("Client closed Normally")
				return
			}
			if websocket.IsUnexpectedCloseError(err,websocket.CloseNormalClosure){
				log.Printf("Client closed AbNormally %v",err)
				return
			}

			log.Printf("Read Error: %v",err)
			return
		}
		var resp Response

		switch req.Action {
		case "ping" :
			resp.Reply = "pong"
		default:
			resp.Reply = "default"
		}

		err = conn.WriteJSON(resp)
		if err!=nil {
			log.Printf("Write error %v",err)
			return
		}
	}
}

func main(){
	http.HandleFunc("/ws",wsHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080",nil))
}
