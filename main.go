package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/Jack-Gitter/tunes/db"
	"github.com/Jack-Gitter/tunes/models/daos"
	"github.com/Jack-Gitter/tunes/models/dtos/responses"
	"github.com/Jack-Gitter/tunes/models/services/auth"
	"github.com/Jack-Gitter/tunes/models/services/cache"
	"github.com/Jack-Gitter/tunes/models/services/comments"
	"github.com/Jack-Gitter/tunes/models/services/jwt"
	"github.com/Jack-Gitter/tunes/models/services/posts"
	"github.com/Jack-Gitter/tunes/models/services/spotify"
	"github.com/Jack-Gitter/tunes/models/services/users"
	"github.com/Jack-Gitter/tunes/server"
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
    defer db.Close()





    redisConnection := cache.GetRedisConnection()
    defer redisConnection.Close()

    //cache := cache.Cache{Redis: redisConnection}

    usersDAO := &daos.UsersDAO{}
    postsDAO := &daos.PostsDAO{}
    commentsDAO := &daos.CommentsDAO{}

    spotifyService := &spotify.SpotifyService{}
    userService := users.UserService{UsersDAO: usersDAO, DB: db}
    postsService := posts.PostsService{PostsDAO: postsDAO, UsersDAO: usersDAO, SpotifyService: spotifyService, DB: db}
    commentsService := comments.CommentsService{CommentsDAO: commentsDAO, DB: db}
    jwtService := &jwt.JWTService{}
    authService := auth.AuthService{UsersDAO: usersDAO, SpotifyService: spotifyService, JWTService: jwtService, DB: db}

	r := server.InitializeHttpServer(&userService, &postsService, &commentsService, &authService)
	r.Run(":2000")

}
