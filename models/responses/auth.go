package responses

type AuthResponse struct {
    AccessToken string
    RefreshToken string
    Data any
}
