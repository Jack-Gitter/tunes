package responses

type Role string

const (
    BASIC_USER Role = "1"
    MODERATOR
    ADMIN
)

type User struct {
    UserIdentifer `mapstructure:",squash"`
    Bio string
    Role Role
    Posts []PostPreview
    Followers []UserIdentifer
    Following []UserIdentifer
}

type UserIdentifer struct {
    Username string
    SpotifyID string
}

