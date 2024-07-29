package responses

import "time"

type Comment struct {

    CommentID int
    Likes int
    Dislikes int
	CommentText string
    CommentorID string
    CommentorUsername string
    PostSpotifyID string
    CreatedAt time.Time
    UpdatedAt time.Time
    SongID string
    

}
