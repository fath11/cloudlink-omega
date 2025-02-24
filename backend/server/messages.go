package cloudlinkOmega

import (
	"log"

	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofrs/uuid"
)

func JSONDump(message any) []byte {
	payload, _ := json.Marshal(message)
	return payload
}

// MulticastMessage broadcasts a payload to multiple clients.
func MulticastMessage(clients map[uuid.UUID]*Client, message []byte, ignoreOrigin *Client) {
	for _, client := range clients {
		// Spawn goroutines to multicast the payload
		UnicastMessage(client, message, ignoreOrigin)
	}
}

// MulticastMessageArray broadcasts a payload to multiple clients.
func MulticastMessageArray(clients []*Client, message []byte, ignoreOrigin *Client) {
	for _, client := range clients {
		// Spawn goroutines to multicast the payload
		UnicastMessage(client, message, ignoreOrigin)
	}
}

// UnicastMessageAny broadcasts a payload to a singular client.
func UnicastMessage(client *Client, message []byte, ignoreOrigin *Client) {
	// Lock state for the websocket connection to prevent accidental concurrent writes to websocket
	client.connectionMutex.Lock()

	if (ignoreOrigin != nil) && (ignoreOrigin == client) {
		client.connectionMutex.Unlock()
		return
	}

	// Attempt to send message to client
	if err := client.connection.WriteMessage(websocket.TextMessage, message); err != nil {
		log.Printf("[%s] Client %s TX ERR: %s", client.manager.Name, client.id, err)
	} else {
		log.Printf("[%s] Client %s TX OK: %s", client.manager.Name, client.id, message)
	}
	client.connectionMutex.Unlock()
}
