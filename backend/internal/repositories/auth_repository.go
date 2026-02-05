package repositories

import (
	"echo-chat-app-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// Interface defining the methods for AuthRepository
type AuthRepository interface {
	SyncUser(uid, email, name, avatarURL string) (*models.User, error)
}

// Semacam class yang mengimplementasikan AuthRepository
type authRepository struct{
	mysqlDB *gorm.DB
}

func NewAuthRepository(mysqlDB *gorm.DB) AuthRepository {
	return &authRepository{mysqlDB: mysqlDB}
}

func (ar *authRepository) SyncUser(uid, email, name, avatarURL string) (*models.User, error) {
	// Inisialisasi user dengan FirebaseUID yang akan dicari
	user := models.User{
        FirebaseUID: uid,
    }

	// Mencari berdasarkan FirebaseUID, jika tidak ada baru buat record baru dengan data tambahan
	err := ar.mysqlDB.Where("firebase_uid = ?", uid).
        Attrs(models.User{ // Attrs hanya dipakai jika record tidak ada
            Email:     email,
            FullName:  name,
            AvatarURL: avatarURL,
            Status:    "offline",
            LastSeen:  time.Now(),
        }).
        FirstOrCreate(&user).Error

	return &user, err
}
