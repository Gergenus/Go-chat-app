package models

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateRoomRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type RoomsResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClientsResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
