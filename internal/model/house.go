package model

import "gorm.io/gorm"

type House struct {
	gorm.Model
	Name     string    `gorm:"size:16;uniqueIndex;not null" json:"name"`
	Code     string    `gorm:"size:16;uniqueIndex;not null" json:"code"`
	Color    string    `gorm:"size:7;not null" json:"color"`
	Students []Student `gorm:"foreignKey:HouseID" json:"students,omitempty"`
}

var DefaultHouses = []House{
	{Name: "Spring", Code: "spring", Color: "#81C784"},
	{Name: "Summer", Code: "summer", Color: "#FFB74D"},
	{Name: "Autumn", Code: "autumn", Color: "#E57373"},
	{Name: "Winter", Code: "winter", Color: "#64B5F6"},
}
