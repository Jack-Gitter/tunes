package responses


type Role string

const (
    BASIC_USER Role = "BASIC"
    MODERATOR Role = "MODERATOR"
    ADMIN Role = "ADMIN"
)

func IsValidRole(role Role) bool {
    return role == BASIC_USER || role == MODERATOR || role == ADMIN
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
