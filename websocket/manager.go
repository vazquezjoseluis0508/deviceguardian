package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Ajusta según tus necesidades de CORS
	},
}

var clients = make(map[*websocket.Conn]bool) // Conexiones de cliente activas

// HandleConnections maneja las solicitudes de conexión WebSocket
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer ws.Close()

	// Registra la nueva conexión
	clients[ws] = true

	// Este bucle mantiene viva la conexión. Puedes implementar lógica adicional según sea necesario.
	for {
		// Mantén la conexión viva o procesa mensajes entrantes
	}
}

// Broadcast envía datos a todos los clientes conectados
func Broadcast(data interface{}) {
	for client := range clients {
		err := client.WriteJSON(data)
		if err != nil {
			log.Printf("Error broadcasting: %v", err)
			client.Close()
			delete(clients, client)
		}
	}
}
