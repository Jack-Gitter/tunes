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
    SongID string
    SongName string
    AlbumName string
    AlbumArtURI string
    AlbumID string
    Rating int
    Text string
}

type Comment struct {
    text string
}
