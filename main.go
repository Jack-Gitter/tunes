package main

import (
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/server"
	"github.com/joho/godotenv"
)

// @title           Tunes backend API
// @version         1.0
// @description     The backend REST API for Tunes

// @contact.name   Jack Gitter
// @contact.email  jack.a.gitter@gmail.com

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	godotenv.Load()

	db.ConnectToDB()
	defer db.DB.Driver.Close()

	r := server.InitializeHttpServer()
	r.Run(":2000")

}
