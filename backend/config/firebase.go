package config

import (
	"os"
	"context"
	"fmt"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

var FirebaseAuth *auth.Client

func InitFirebase() error {
    ctx := context.Background()
    var app *firebase.App
    var err error

    // 1. Cek apakah sedang menggunakan Emulator (untuk local development)
    emulatorHost := os.Getenv("FIREBASE_AUTH_EMULATOR_HOST")
    
    if emulatorHost != "" {
        // Jika emulator aktif, kita inisialisasi tanpa file key
        app, err = firebase.NewApp(ctx, &firebase.Config{ProjectID: "your-project-id"})
    } else {
        // 2. Jika tidak ada emulator, cari path file JSON dari .env
        serviceAccountPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
        
        if serviceAccountPath != "" {
            // Gunakan file JSON jika path tersedia
            opt := option.WithServiceAccountFile(serviceAccountPath)
            app, err = firebase.NewApp(ctx, nil, opt)
        } else {
            // 3. Mode Production Sejati (Google Cloud Run/App Engine)
            // Di server Google, kita tidak butuh file JSON sama sekali
            app, err = firebase.NewApp(ctx, nil)
        }
    }

    if err != nil {
        return fmt.Errorf("error initializing app: %v", err)
    }

	client, err := app.Auth(context.Background())
	if err != nil {
		return fmt.Errorf("error getting Auth client: %v", err)
	}

	FirebaseAuth = client
	return nil
}