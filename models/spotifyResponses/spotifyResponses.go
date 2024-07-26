package spotifyresponses

type AccessTokenResponnse struct {
	Access_token  string
	Token_type    string
	Scope         string
	Expires_in    int
	Refresh_token string
}

type ProfileResponse struct {
	Id           string
	Display_name string
}

type RefreshTokenResponse struct {
	Access_token  string
	Expires_in    int
	Refresh_token string
}

type SongResponse struct {
	Id    string
	Name  string
	Album AlbumResponse
}

type AlbumResponse struct {
	Id     string
	Name   string
	Images []Images
}

type Images struct {
	Url string
}
