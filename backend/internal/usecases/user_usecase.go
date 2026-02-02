package usecases

import "echo-chat-app-backend/internal/repositories"

type UserUseCase struct {
	userRepo repositories.UserRepository
}