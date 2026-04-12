package model

import "gorm.io/gorm"

type Student struct {
	gorm.Model
	StudentID string `gorm:"size:64" json:"student_id"`
	Name      string `gorm:"size:64;not null" json:"name"`
	HouseID   uint   `gorm:"not null;index" json:"house_id"`
	House     House  `gorm:"foreignKey:HouseID" json:"house,omitempty"`
	LapCount  int    `gorm:"not null;default:0" json:"lap_count"`
}
