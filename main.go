package main

import "github.com/Jack-Gitter/tunes/handlers"

func main() {

    r := handlers.InitializeHttpServer()
    r.Run(":2000")

}
