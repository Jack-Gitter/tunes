package models

type Post struct {
    PostPreview `mapstructure:",squash"`
    // []Comments Comments
}

type PostPreview struct {
    UserIdentifer `mapstructure:",squash"`
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
