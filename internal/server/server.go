package server

import (
	"net/http"
	"os"
	"time"

	"github.com/Gergenus/StandardLib/internal/handler"
	"github.com/Gergenus/StandardLib/internal/middleware"
	"github.com/Gergenus/StandardLib/internal/repository"
	"github.com/Gergenus/StandardLib/internal/service"
	"github.com/Gergenus/StandardLib/internal/ws"
	"github.com/Gergenus/StandardLib/pkg"
	"github.com/labstack/echo/v4"
)

type Server interface {
	Start()
	InitializationRouts()
}

type EchoServer struct {
	App *echo.Echo
	DB  pkg.DBwraper
}

func NewEchoServer(DB pkg.DBwraper) EchoServer {
	App := echo.New()
	return EchoServer{
		App: App,
		DB:  DB,
	}
}

func (e *EchoServer) Start() {
	s := http.Server{
		Addr:         os.Getenv("HTTPPORT"),
		Handler:      e.App,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  60 * time.Second,
	}
	s.ListenAndServe()
}

func (e *EchoServer) InitializationRouts() {

	hash := pkg.NewSHAhash(os.Getenv("SALT"))
	jwtPKG := pkg.NewJWTgo()
	userRepo := repository.PostgresUserRepo{DB: e.DB, Hash: &hash}
	authService := service.JWTauth{UserRepo: &userRepo, Hasher: &hash, Auther: &jwtPKG}
	handlerAuth := handler.NewEchoHandlerAuth(&authService)
	middle := middleware.NewEchoMiddleware(&jwtPKG)

	hub := ws.NewHub()
	hubHandler := ws.NewHandler(hub)

	e.App.POST("/SignUp", handlerAuth.SignUp)                                  // {"username": "Denis", "email": "123", "password": "123"}
	e.App.POST("/SignIn", handlerAuth.SignIn)                                  // {"username": "Denis", "email": "123", "password": "123"}
	e.App.POST("/ws/createRoom", hubHandler.CreateRoom, middle.AuthMiddleware) // { "id": "1", "name": "chatting"}
	e.App.GET("/ws/joinRoom/:roomId", hubHandler.JoinRoom, middle.WSAuthMiddleware)

	e.App.GET("ws/getRooms", hubHandler.GetRooms, middle.AuthMiddleware)
	e.App.GET("/ws/getClients", hubHandler.GetClients, middle.AuthMiddleware)

	go hub.Run()
}
