package main

import (
	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/daos"
	"github.com/Jack-Gitter/tunes/server"
	"github.com/Jack-Gitter/tunes/server/auth"
	"github.com/Jack-Gitter/tunes/server/comments"
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
    usersDAO := daos.UsersDAO{DB: db}
    userService := users.UserService{UsersDAO: &usersDAO}
    postsDAO := daos.PostsDAO{DB: db}
    postsService := posts.PostsService{PostsDAO: &postsDAO}
    commentsDAO := daos.CommentsDAO{DB: db}
    commentsService := comments.CommentsService{CommentsDAO: &commentsDAO}
    authService := auth.AuthService{UsersDAO: &usersDAO}


    defer db.Close()

	r := server.InitializeHttpServer(&userService, &postsService, &commentsService, &authService)
	r.Run(":2000")

}
