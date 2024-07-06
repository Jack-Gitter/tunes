package models

type Post struct {
    PostMetaData
    UserIdentifer
    // []Comments Comments
}

type PostMetaData struct {
    SongID string
    SongName string
    AlbumName string
    AlbumArtURI string
    AlbumID string
    Rating int
    Text string
    Likes []UserIdentifer
    Dislikes []UserIdentifer
}
