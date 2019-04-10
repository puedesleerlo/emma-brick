package main

import (
	"log"
	"github.com/sacOO7/gowebsocket"
	"os"
	"os/signal"
	"encoding/json"
	// "github.com/jacobsa/go-serial/serial"
	"github.com/ev3go/ev3dev"
	// "github.com/kraman/go-firmata"
	// gobot"gobot.io/x/gobot"
	// aio"gobot.io/x/gobot/drivers/aio"
)


type Message struct {
    Sender    string `json:"sender,omitempty"`
    Recipient string `json:"recipient,omitempty"`
    Type     string `json:"type,omitempty"`
    Content   map[string]map[int]interface{} `json:"content,omitempty"`
}

var outA, _ = ev3dev.TachoMotorFor("ev3-ports:outA", "lego-ev3-m-motor")

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)


	socket := gowebsocket.New("ws://emmago.hopto.org/programs/ws")		
	
	SocketConfig(&socket)

	socket.Connect()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			// socket.Close()
			// port.Close()
			return
		}
	}
}

func SocketConfig(socket *gowebsocket.Socket) {
	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Fatal("Received connect error - ", err)
	}
  
	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server");
	}
  
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		// log.Println("Received message - " + message)
		received := Message{}
		err := json.Unmarshal([]byte(message), &received)
		check(err)
		if received.Type == "driver" {
			// angle := received.Content["variables"][0]
			outA.Command("run-forever")
		}
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	}
}

// func PrepareMessage(v int) Message {
// 	content := map[string]interface{}{}
// 	variables := map[int]interface{}{}
// 	variables[0] = v
// 	content["variables"] = variables
// 	message := Message{Type:"sensor", Content: content,}
// 	return message
// }

func check(e error) {
    if e != nil {
        panic(e)
    }
}