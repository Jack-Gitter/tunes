package main

import (
	"github.com/Jack-Gitter/tunes/handlers"
	"github.com/joho/godotenv"
)

func main() {

    godotenv.Load()

    r := handlers.InitializeHttpServer()
    r.Run(":2000")

}
