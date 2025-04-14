package ws

import (
	"net/http"

	"github.com/Gergenus/StandardLib/internal/models"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Hub *Hub
}

func NewHandler(Hub *Hub) *Handler {
	return &Handler{Hub: Hub}
}

func (h *Handler) CreateRoom(c echo.Context) error {
	var req models.CreateRoomRequest

	err := c.Bind(&req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}

	h.Hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}
	return c.JSON(http.StatusCreated, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (h *Handler) JoinRoom(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{
			"error": "Invalid request payload",
		})
	}
	roomId := c.Param("roomId")
	clientId := c.Get("name").(string)
	username := clientId

	cl := &Client{
		Socket:   conn,
		Message:  make(chan *Message, 10),
		ID:       clientId,
		RoomID:   roomId,
		Username: username,
	}

	m := &Message{
		Content:  "A new user has joined the room",
		RoomID:   roomId,
		Username: username,
	}

	h.Hub.Register <- cl
	h.Hub.Broadcast <- m

	go cl.WriteMessage()
	cl.ReadMessage(h.Hub)

	return c.JSON(http.StatusOK, map[string]string{
		"status": "leaved",
	})
}

func (h *Handler) GetRooms(c echo.Context) error {
	rooms := []models.RoomsResponse{}
	for _, d := range h.Hub.Rooms {
		rooms = append(rooms, models.RoomsResponse{
			ID:   d.ID,
			Name: d.Name,
		})
	}
	return c.JSON(http.StatusOK, rooms)
}

func (h *Handler) GetClients(c echo.Context) error {
	clients := []models.ClientsResponse{}
	roomId := c.QueryParam("roomId")
	for _, d := range h.Hub.Rooms[roomId].Clients {
		clients = append(clients, models.ClientsResponse{
			ID:       d.ID,
			Username: d.Username,
		})
	}
	return c.JSON(http.StatusOK, clients)
}
