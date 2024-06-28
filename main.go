package main

import "github.com/Jack-Gitter/tunes/server"

func main() {

    r := server.InitializeHttpServer()
    r.Run(":2000")

}
