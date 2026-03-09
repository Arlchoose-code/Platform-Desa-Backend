// Platform Desa — User Seeder
// Copyright (c) 2026 Syahril Haryono
// Licensed under MIT License

package database

import (
	"fmt"

	"github.com/Arlchoose-code/platform-desa-backend/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SeedUsers(db *gorm.DB) {
	users := []struct {
		Name     string
		Email    string
		Password string
		Role     string
		Phone    string
	}{
		{
			Name:     "Super Admin",
			Email:    "superadmin@desa.id",
			Password: "superadmin123",
			Role:     "superadmin",
			Phone:    "081234567890",
		},
		{
			Name:     "Admin Desa",
			Email:    "admin@desa.id",
			Password: "admin123",
			Role:     "admin",
			Phone:    "081234567891",
		},
		{
			Name:     "Operator Desa",
			Email:    "operator@desa.id",
			Password: "operator123",
			Role:     "operator",
			Phone:    "081234567892",
		},
	}

	for _, u := range users {
		var existing models.User
		if err := db.Where("email = ?", u.Email).First(&existing).Error; err == nil {
			continue
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println("Gagal hash password untuk:", u.Email)
			continue
		}

		phone := u.Phone
		db.Create(&models.User{
			Name:     u.Name,
			Email:    u.Email,
			Password: string(hashed),
			Role:     u.Role,
			Phone:    &phone,
			IsActive: true,
		})

		fmt.Printf("User seeded: %s (%s)\n", u.Email, u.Role)
	}
}
