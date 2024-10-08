package responses

import (
	"time"
)

type PostPreview struct {
	UserIdentifer `mapstructure:",squash"`
	SongID        string
	SongName      string
	AlbumName     string
	AlbumArtURI   string
	AlbumID       string
	Rating        int
	Text          string
	Likes         []UserIdentifer
	Dislikes      []UserIdentifer
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
