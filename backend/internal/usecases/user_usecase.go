package usecases

import (
	"echo-chat-app-backend/internal/models"
	"echo-chat-app-backend/internal/repositories"
)

type UserUseCase struct {
	userRepo repositories.UserRepository
}

func NewUserUseCase(userRepo repositories.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) Me(uid string) (*models.User, error) {
	return uc.userRepo.Me(uid)
}

func (uc *UserUseCase) SearchUserByUsername(username string) (*models.User, error) {
	return uc.userRepo.SearchUserByUsername(username)
}

func (uc *UserUseCase) UpdateProfile(uid, name, username, avatar_url string) (*models.User, error) {
	return uc.userRepo.UpdateProfile(uid, name, username, avatar_url)
}
