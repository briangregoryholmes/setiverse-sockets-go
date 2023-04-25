package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Room struct {
    Name        string
    Connections map[*websocket.Conn]bool
}

type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

var rooms = make(map[string]*Room)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
        // Only allow connections from localhost
        return true
    },
}

func main() {
    http.HandleFunc("/", wsHandler)

    fmt.Println("Starting server on :4000")
    err := http.ListenAndServe(":4000", nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
    // Parse the room name from the URL path
    roomName := r.URL.Path[len("/"):]

	// If no room name, refuse connection
	if roomName == "" {
		http.Error(w, "No room name specified", http.StatusBadRequest)
		return
	}

	// Print the room name to the console
	fmt.Println("Room name:", roomName)

    // Check if a room with this name already exists
    room, ok := rooms[roomName]

    if !ok {
        // If it doesn't, create a new room
        room = &Room{
            Name:        roomName,
            Connections: make(map[*websocket.Conn]bool),
        }
        rooms[roomName] = room
    }

    // Upgrade the HTTP connection to a WebSocket connection
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println(err)
        return
    }

    // Add the new connection to the room
    room.Connections[conn] = true

	// Ping the client every 30 seconds to keep the connection alive
	go func() {
		for {
			// Send a ping message to all clients in the room
			for conn := range room.Connections {
				log.Println("Sending ping")
				err := conn.WriteMessage(websocket.PingMessage, []byte{})
				if err != nil {
					log.Println("Ping error:", err)
					delete(room.Connections, conn)
				}
			}

			// Wait for 30 seconds before pinging again
			time.Sleep(30 * time.Second)
		}
	}()
	

    // Start a goroutine to read messages from this connection and broadcast them to the room
     go func() {
        for {
            _, msg, err := conn.ReadMessage()

            if err != nil {
                log.Println("Read error:", err)
                delete(room.Connections, conn)
                break
            }

            // Parse the message as JSON
            var parsedMessage Message
            err = json.Unmarshal(msg, &parsedMessage)

            // Console log the event parameter
            fmt.Println("Event Received:", parsedMessage.Data)

            // Set up a switch case to handle different events
            switch parsedMessage.Event {
            case "CARD":
                // Do something when a user joins
                fmt.Println(parsedMessage.Data)
                fmt.Println(room.Connections)
            case "SET":
                // Do something when a user sends a message
                fmt.Println("Called set")
            }

            if err != nil {
                log.Println("Read error:", err)
                delete(room.Connections, conn)
                break
            }

            // Broadcast the message to all clients in the room, except the sender
            for connOther := range room.Connections {
                if connOther == conn {
                    continue
                }

                if err := connOther.WriteMessage(websocket.TextMessage, msg); err != nil {
                    log.Println("Write error:", err)
                    delete(room.Connections, connOther)
                    break
                }
            }
        }
    }()
}