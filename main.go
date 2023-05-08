package main

import (
	// "encoding/json"
	"fmt"
	"go-fiber-websocket/database"
	"go-fiber-websocket/models"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"

	// "github.com/google/uuid"

	// "github.com/gofrs/uuid"

	// "github.com/gofiber/adaptor/v2/fasthttpadaptor"
	gorilla "github.com/gorilla/websocket"
	// "github.com/valyala/fasthttp"
)

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define a struct to hold the client ID and WebSocket connection

func main() {
	database.ConnectDb()
	app := fiber.New()

	// Create a map to keep track of the clients
	clients := make(map[string]*models.Client)

	app.Use("/ws", func(c *fiber.Ctx) error {

		// username := c.Params("username")
		// fmt.Println("username", username)

		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// Get the client ID from the URL parameter
		// username := c.Params("username")
		// msg := c.Query("msg")
		id := c.Params("id")
		// id := uuid.New()
		// fmt.Println("id", id)
		// fmt.Println("id", msg)
		// var clientss []models.Client
		// Add the client to the map
		clients[id] = &models.Client{Id: id, Conn: c}
		// clientss = append(clientss, clients)
		// database.Database.Db.Create(&clients)

		// fmt.Println("clients", clients[id])
		// Close the connection when the client disconnects
		defer func() {
			delete(clients, id)
			err := c.Close()
			if err != nil {
				log.Println(err)
			}
		}()

		// Continuously read messages from the client
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("Received message from client %s: %s", id, string(msg))
			// msgss := string(msg)
			// uniqueID := uuid.New()
			// msgs := models.Message{
			// 	ID:      uniqueID.String(),
			// 	Message: msg,
			// 	To:      id,
			// 	From:    id,
			// }

			// ids, err := strconv.Atoi(id)
			// SendMessage(&msgs, clientss)
		}
	}))

	// Handler to send a message to a specific client
	app.Get("/send/:id", func(c *fiber.Ctx) error {
		// type payloads struct{
		// 	message  []byte
		// 	to string 
		// }
		// var payload payloads
		// if err := c.BodyParser(&payload); err != nil {
		// 	return err
		// }
		// message := payload.message
		// to := payload.to
		id := c.Params("id")
		msg := c.Query("msg")
		// id := c.Params("")
		// fmt.Println(msg)
		// var dict map[string]string
    // err := json.Unmarshal([]byte(msg), &dict)
		// if err != nil {
		// 	return err
		// }

		// fmt.Println(dict)

		// fmt.Println("id", id)
		fmt.Println("msg", msg)


		// Get the client from the map
		client, ok := clients[id]
		if !ok {
			return fiber.NewError(fiber.StatusNotFound, "Client not found")
		}

		// Send the message to the client
		err := client.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println(err)
			return err
		}

		// msgss := string(msg)
			uniqueID := uuid.New()
			roomId := uuid.New().String()
			msgs := models.Message{
				ID:      uniqueID.String(),
				Message: []byte(msg),//message
				RoomId: roomId,
				To:      id,//to
				From:    id,
			}

			database.Database.Db.Create(&msgs)

		return nil
	})

	// Convert the fiber.Handler to an http.Handler using the fiber.Handler's func(*fiber.Ctx) error method signature
	// handler := http.Handler(app.Handler)

	// Start the server
	err := app.Listen(":3000")
	if err != nil {
		log.Fatal(err)
	}
}

func SendMessage(b *models.Message, clients []models.Client) error {
	log.Println("sending message")
	// clients := client{}
	// id = 2
	client := models.Client{}
	database.Database.Db.Find(&client, "id =?", b.To)

	// if !ok {
	// 	return fiber.NewError(fiber.StatusNotFound, "Client not found")
	// }

	err := client.Conn.WriteMessage(websocket.TextMessage, []byte(b.Message))
	if err != nil {
		log.Println(err)
		return err
	}

	log.Println("done sending message")
	return nil
}
