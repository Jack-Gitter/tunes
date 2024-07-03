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
