package usecases

import (
	"echo-chat-app-backend/internal/models"
	"echo-chat-app-backend/internal/repositories"
)

type AuthUseCase struct {
	authRepo repositories.AuthRepository
}

func NewAuthUseCase(authRepo repositories.AuthRepository) *AuthUseCase {
	return &AuthUseCase{
		authRepo: authRepo,
	}
}

func (au *AuthUseCase) SyncUser(uid, email, name, avatarURL string) (*models.User, error) {
	// Panggil metode dari authRepo untuk menyinkronkan data pengguna
	return au.authRepo.SyncUser(uid, email, name, avatarURL)
}
