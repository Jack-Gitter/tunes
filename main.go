package main

import (
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/server"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
    _ "github.com/Jack-Gitter/tunes/docs"
)

func main() {

	godotenv.Load()

	db.ConnectToDB()
	defer db.DB.Driver.Close()

	r := server.InitializeHttpServer()
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":2000")

}
