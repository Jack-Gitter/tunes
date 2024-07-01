package main

import (
	//"encoding/json"
	//"fmt"

	"fmt"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/server"
	"github.com/joho/godotenv"
)

func main() {

    godotenv.Load()
    db.ConnectToDB()
    result := db.GetUserFromDbBySpotifyID("id")
    fmt.Println(result.Records[0].Values)


    r := server.InitializeHttpServer()
    r.Run(":2000")

}
