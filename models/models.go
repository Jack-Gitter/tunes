package models

type Role int

const (
    BASIC_USER Role = iota
    MODERATOR
    ADMIN
)

type User struct {
    Username string
    SpotifyID string
    Bio string
    Role Role
}

type Post struct {
    songID string
    songName string
    albumName string
    albumArtURI string
    albumID string
    rating int
    text string
}

type Comment struct {
    text string
}
