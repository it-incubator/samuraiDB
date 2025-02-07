package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"path/filepath"
	"strings"

	fa "samurai-db/internal/file_adapter"
	im "samurai-db/internal/index_manager"
	sdb "samurai-db/internal/samuraidb"
	sm "samurai-db/internal/segment_manager"

	"github.com/google/uuid"
)

type RequestAction struct {
	Type    string                 `json:"type"`
	Payload map[string]interface{} `json:"payload"`
	UUID    string                 `json:"uuid"`
}

func main() {
	dir := filepath.Join("db")
	fileAdapter := fa.NewAdapter(dir)
	segmentManager := sm.NewSegmentManager(fileAdapter)
	indexManager := im.NewIndexManager(fileAdapter)
	db := sdb.NewSamuraiDB(segmentManager, indexManager)

	// Инициализация базы данных
	if err := db.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Создание TCP сервера
	listener, err := net.Listen("tcp", ":4001")
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer listener.Close()

	log.Println("Server listening on port 4001")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v\n", err)
			continue
		}

		go handleConnection(conn, db)
	}
}

func handleConnection(conn net.Conn, db *sdb.SamuraiDB) {
	defer conn.Close()
	log.Println("Client connected")

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		data := scanner.Text()
		log.Printf("Received from client: %s", data)

		var requestAction RequestAction
		if err := json.Unmarshal([]byte(data), &requestAction); err != nil {
			log.Printf("Failed to parse request: %v", err)
			fmt.Fprintf(conn, "Invalid request format\n")
			continue
		}

		log.Printf("Unmarshal from client: %s", requestAction)

		switch strings.ToUpper(requestAction.Type) {
		case "SET":
			id := uuid.New().String()
			requestAction.Payload["id"] = id

			if err := db.Set(id, requestAction.Payload); err != nil {
				log.Printf("Failed to set value: %v", err)
				fmt.Fprintf(conn, "Error setting value\n")
			}

			response := map[string]interface{}{
				"id":   id,
				"uuid": requestAction.UUID,
			}
			for k, v := range requestAction.Payload {
				response[k] = v
			}
			responseData, _ := json.Marshal(response)
			conn.Write(responseData)
			log.Printf("Response: %s", responseData)

		case "GET":
			id, ok := requestAction.Payload["id"].(string)
			if !ok {
				fmt.Fprintf(conn, "Invalid id format\n")
			}
			data, err := db.Get(id)
			if err != nil || data == nil {
				fmt.Fprintf(conn, "Data not found\n")
			}

			response := map[string]interface{}{
				"uuid": requestAction.UUID,
			}
			for k, v := range data {
				response[k] = v
			}

			responseData, _ := json.Marshal(response)
			conn.Write(responseData)
			log.Printf("Response: %s", responseData)

		default:
			log.Printf("Unknown request type: %s", requestAction.Type)
			conn.Write([]byte("Unknown request type\n"))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Client error: %v", err)
	}

	log.Println("Client disconnected")
}
