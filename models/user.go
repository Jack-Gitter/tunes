package models

type Role int

const (
    BASIC_USER Role = iota
    MODERATOR
    ADMIN
)

type User struct {
    UserIdentifer
    Bio string
    Role Role
    Posts []PostInformationForUser
    Followers []UserIdentifer
    Following []UserIdentifer
}

type UserIdentifer struct {
    Username string
    SpotifyID string
}

