package service

// AuthService handles authentication business logic (Argon2id hashing, session management).
type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// TODO: Login, Logout, ValidateSession, HashPassword, VerifyPassword
