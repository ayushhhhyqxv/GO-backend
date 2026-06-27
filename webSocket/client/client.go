package main

import (
    "log"
    "time"
    "github.com/gorilla/websocket"
)

func main() {

	dialer := websocket.Dialer{}
	conn,_,err := dialer.Dial("ws://localhost:8080/ws",nil)

	if err!=nil{
		log.Fatal(err)
	}

	defer conn.Close()

	req := map[string]string{"action":"ping"}
	if err:= conn.WriteJSON(req);err!=nil{
		log.Fatal(err)
	}

	var resp map[string]string 

	if err:=conn.ReadJSON(&resp);err!=nil{
		log.Fatal(err)
	}
	
	log.Printf("Received: %v", resp)

	conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))

	time.Sleep(100 * time.Millisecond)

}