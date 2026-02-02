package repositories

import (
	"gorm.io/gorm"
	
)

// Interface defining the methods for AuthRepository
type AuthRepository interface {
	SyncUser(uid, email string) error
}

// Semacam class yang mengimplementasikan AuthRepository
type authRepository struct{
	mysqlDB *gorm.DB
}

func NewAuthRepository(mysqlDB *gorm.DB) AuthRepository {
	return &authRepository{mysqlDB: mysqlDB}
}

func (ar *authRepository) SyncUser(uid, email string) error {
	// TODO: Implement the logic to sync user data with the database
	return nil
}
