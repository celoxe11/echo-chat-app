package repositories

import (
	"context"
	"echo-chat-app-backend/internal/models"

	"firebase.google.com/go/v4/auth"
	"gorm.io/gorm"
)

type UserRepository interface {
	Me(uid string) (*models.User, error)
	SearchUserByUsername(username string) (*models.User, error)
	UpdateProfile(uid, name, username, avatar_url string) (*models.User, error)
}

type userRepository struct {
	authClient *auth.Client
	mysqlDB    *gorm.DB
}

func NewUserRepository(authClient *auth.Client, mysqlDB *gorm.DB) UserRepository {
	return &userRepository{authClient: authClient, mysqlDB: mysqlDB}
}

func (ur *userRepository) Me(uid string) (*models.User, error) {
	user := models.User{}
	err := ur.mysqlDB.Where("firebase_uid = ?", uid).First(&user).Error
	return &user, err
}

func (ur *userRepository) SearchUserByUsername(username string) (*models.User, error) {
	user := models.User{}
	err := ur.mysqlDB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (ur *userRepository) UpdateProfile(uid, name, username, avatar_url string) (*models.User, error) {
	ctx := context.Background()
	user := models.User{}

	// update user di firebase
	params := (&auth.UserToUpdate{}).
		DisplayName(name).
		PhotoURL(avatar_url)

	_, err := ur.authClient.UpdateUser(ctx, uid, params)
	if err != nil {
		return nil, err
	}

	// update user di mysql
	err = ur.mysqlDB.Model(&models.User{}).
		Where("firebase_uid = ?", uid).
		Updates(models.User{
			FullName:  name,
			Username:  username,
			AvatarURL: avatar_url,
		}).Error
	if err != nil {
		return nil, err
	}

	// Fetch updated user to return
	err = ur.mysqlDB.Where("firebase_uid = ?", uid).First(&user).Error
	return &user, err
}
