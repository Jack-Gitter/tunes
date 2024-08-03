package main

import (
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/dtos"
	"github.com/Jack-Gitter/tunes/server"
	"github.com/Jack-Gitter/tunes/server/posts"
	"github.com/Jack-Gitter/tunes/server/users"
	"github.com/joho/godotenv"
)

// @title           Tunes backend API
// @version         1.0
// @description     The backend REST API for Tunes

// @contact.name   Jack Gitter
// @contact.email  jack.a.gitter@gmail.com

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description "Authorization header value"
func main() {

	godotenv.Load()

    db := db.ConnectToDB()
    usersDTO := dtos.UsersDTO{DB: db}
    userService := users.UserService{UsersDTO: usersDTO}
    postsDTO := dtos.PostsDTO{DB: db}
    postsService := posts.PostsService{PostsDTO: postsDTO}

    defer db.Close()

	r := server.InitializeHttpServer(userService, postsService)
	r.Run(":2000")

}
