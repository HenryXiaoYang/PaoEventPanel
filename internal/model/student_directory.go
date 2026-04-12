package model

// StudentDirectoryEntry represents a student from the external directory database.
// Not managed by GORM — read from a separate SQLite file (data/students.db).
type StudentDirectoryEntry struct {
	ID          string `json:"id"`
	EnglishName string `json:"english_name"`
	ChineseName string `json:"chinese_name"`
	SystemName  string `json:"system_name"`
}
