package main

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

func main() {
	e := echo.New()
	s := http.Server{
		Addr:         ":8080",
		Handler:      e,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  60 * time.Second,
	}

	pkg.Load()
	db := pkg.InitDB(os.Getenv("USER"), os.Getenv("PASSWORD"), os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("DBNAME"), os.Getenv("SSLMODE"))

	hash := pkg.NewSHAhash(os.Getenv("SALT"))
	jwtPKG := pkg.NewJWTgo()
	userRepo := repository.PostgresUserRepo{DB: db, Hash: &hash}
	authService := service.JWTauth{UserRepo: &userRepo, Hasher: &hash, Auther: &jwtPKG}
	handlerAuth := handler.NewEchoHandlerAuth(&authService)
	middle := middleware.NewEchoMiddleware(&jwtPKG)

	hub := ws.NewHub()
	hubHandler := ws.NewHandler(hub)

	e.POST("/SignUp", handlerAuth.SignUp)
	e.POST("/SignIn", handlerAuth.SignIn)
	e.GET("/", handler.Access, middle.AuthMiddleware)
	e.POST("/ws/createRoom", hubHandler.CreateRoom, middle.AuthMiddleware)
	e.GET("/ws/joinRoom/:roomId", hubHandler.JoinRoom, middle.WSAuthMiddleware)

	go hub.Run()

	s.ListenAndServe()

}
