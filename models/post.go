package models

type Post struct {
    PostInformationForUser
    UserIdentifer
    // []Comments Comments
}

type PostInformationForUser struct {
    SongID string
    SongName string
    AlbumName string
    AlbumArtURI string
    AlbumID string
    Rating int
    Text string
}
