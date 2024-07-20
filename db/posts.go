package db

import (
	"time"
	"github.com/Jack-Gitter/tunes/models/responses"
//	"github.com/mitchellh/mapstructure"
)

/* ===================== CREATE =====================  */

func CreatePost(spotifyID string, songID string, songName string, albumID string, albumName string, albumImage string, rating int, text string, createdAt time.Time) (*responses.Post, error) {
    return nil, nil
}

/* ===================== READ =====================  */

func GetUserPostByID(postID string, spotifyID string) (*responses.Post, bool, error) {
    return nil, false, nil
}

// make this method get the posts with id offset -> offset+limit-1
func GetUserPostsPreviewsByUserID(spotifyID string, createdAt time.Time) (*responses.PaginationResponse[[]responses.PostPreview, time.Time], error) {
    return nil, nil
}

func GetUserPostPreviewByID(songID string, spotifyID string) (*responses.PostPreview, bool, error){
    return nil, false, nil
}


/* ===================== DELETE =====================  */

func DeletePost(songID string, spotifyID string) (bool, bool, error) {
    return false, false, nil
}


/* PROPERTY UPDATES */
func UpdatePost(spotifyID string, songID string, text *string, rating *int) (*responses.PostPreview, bool, error) {
    return nil, false, nil
}
