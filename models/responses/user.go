package responses

type Role string

const (
    BASIC_USER Role = "BASIC"
    MODERATOR Role = "MODERATOR"
    ADMIN Role = "ADMIN"
)

func IsValidRole(role string) bool {
    return role == string(BASIC_USER) || role == string(MODERATOR) || role == string(ADMIN)
}

type User struct {
    UserIdentifer `mapstructure:",squash"`
    Bio string
    Role Role
}

type UserIdentifer struct {
    Username string
    SpotifyID string
}

