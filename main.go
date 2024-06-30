package main

import (
	"github.com/Jack-Gitter/tunes/server"
	"github.com/joho/godotenv"
)

func main() {

    godotenv.Load()

    r := server.InitializeHttpServer()
    r.Run(":2000")

}
