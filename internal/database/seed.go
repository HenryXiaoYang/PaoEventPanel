package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/keren/pao-event-panel/internal/model"
)

func Seed(superAdminUser, superAdminPass string) {
	seedHouses()
	seedSuperAdmin(superAdminUser, superAdminPass)
	seedActivitySettings()
}

func seedHouses() {
	for _, h := range model.DefaultHouses {
		var count int64
		DB.Model(&model.House{}).Where("code = ?", h.Code).Count(&count)
		if count == 0 {
			if err := DB.Create(&h).Error; err != nil {
				log.Printf("Failed to seed house %s: %v", h.Name, err)
			}
		}
	}
}

func seedSuperAdmin(username, password string) {
	var count int64
	DB.Model(&model.User{}).Where("role = ?", model.RoleSuperAdmin).Count(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash super admin password: %v", err)
	}

	user := model.User{
		Username: username,
		Password: string(hash),
		Role:     model.RoleSuperAdmin,
	}
	if err := DB.Create(&user).Error; err != nil {
		log.Fatalf("Failed to create super admin: %v", err)
	}
	log.Printf("Super admin '%s' created", username)
}

func seedActivitySettings() {
	var count int64
	DB.Model(&model.ActivitySettings{}).Count(&count)
	if count > 0 {
		return
	}

	settings := model.DefaultActivitySettings
	if err := DB.Create(&settings).Error; err != nil {
		log.Printf("Failed to seed activity settings: %v", err)
	}
}
