package main

import (
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/server"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	db.ConnectToDB()
    defer db.DB.Driver.Close()

	r := server.InitializeHttpServer()
	r.Run(":2000")

}
