package dto

type CreateAccountRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
    PlayerName string `json:"player_name"`
}

type CreateAccountResponse struct {
    AccountID string `json:"account_id"`
    PlayerName string `json:"player_name"`
    PlayerID string `json:"player_id"`
    Email string `json:"email"`
}

type LoginRequest struct {
    Email string `json:"email"`
    Password string `json:"password"`
}

type LoginResponse struct {
    Token string `json:"token"`
    ExpiresAt int64 `json:"expires_at"`
}