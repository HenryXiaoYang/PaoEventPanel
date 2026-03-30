package database

import "github.com/keren/pao-event-panel/internal/model"

func Migrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.House{},
		&model.Student{},
		&model.LapRecord{},
		&model.ActivitySettings{},
		&model.ThemePreset{},
	)
}
