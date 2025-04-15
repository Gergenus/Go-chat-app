package main

import (
	"os"

	"github.com/Gergenus/StandardLib/internal/server"
	"github.com/Gergenus/StandardLib/pkg"
)

func main() {
	pkg.Load()
	db := pkg.InitDB(os.Getenv("USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("HOST"), os.Getenv("PORT"), os.Getenv("DBNAME"), os.Getenv("SSLMODE"))

	Server := server.NewEchoServer(db)
	Server.InitializationRouts()

	Server.Start()

}
