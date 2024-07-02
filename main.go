package main

import (
	//"encoding/json"
	//"fmt"

	"fmt"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/server"
	"github.com/joho/godotenv"
	//"github.com/mitchellh/mapstructure"
)

func main() {

    godotenv.Load()
    db.ConnectToDB()
    result, _ := db.GetUserFromDbBySpotifyID("id")
    fmt.Println(result)
    //mapstructure.Decode()

    r := server.InitializeHttpServer()
    r.Run(":2000")

}
