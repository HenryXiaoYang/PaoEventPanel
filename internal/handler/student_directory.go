package handler

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"github.com/keren/pao-event-panel/internal/model"
)

const studentDirectoryPath = "data/students.db"

func SearchStudentDirectory(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusOK, []model.StudentDirectoryEntry{})
		return
	}

	db, err := openDirectoryDB()
	if err != nil {
		c.JSON(http.StatusOK, []model.StudentDirectoryEntry{})
		return
	}
	defer db.Close()

	pattern := "%" + q + "%"
	rows, err := db.Query(
		`SELECT id, COALESCE(english_name, ''), COALESCE(chinese_name, ''), system_name
		 FROM students
		 WHERE id LIKE ? OR english_name LIKE ? OR chinese_name LIKE ? OR system_name LIKE ?
		 LIMIT 10`,
		pattern, pattern, pattern, pattern,
	)
	if err != nil {
		c.JSON(http.StatusOK, []model.StudentDirectoryEntry{})
		return
	}
	defer rows.Close()

	results := make([]model.StudentDirectoryEntry, 0)
	for rows.Next() {
		var entry model.StudentDirectoryEntry
		if err := rows.Scan(&entry.ID, &entry.EnglishName, &entry.ChineseName, &entry.SystemName); err != nil {
			continue
		}
		results = append(results, entry)
	}
	c.JSON(http.StatusOK, results)
}

func GetDirectoryStatus(c *gin.Context) {
	db, err := openDirectoryDB()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"loaded": false, "count": 0})
		return
	}
	defer db.Close()

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM students").Scan(&count); err != nil {
		c.JSON(http.StatusOK, gin.H{"loaded": false, "count": 0})
		return
	}
	c.JSON(http.StatusOK, gin.H{"loaded": true, "count": count})
}

func DownloadDirectoryTemplate(c *gin.Context) {
	tmpFile, err := os.CreateTemp("", "template-*.db")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create template"})
		return
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	db, err := sql.Open("sqlite3", tmpPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create template"})
		return
	}
	_, err = db.Exec(`CREATE TABLE students (
		id TEXT PRIMARY KEY,
		english_name TEXT,
		chinese_name TEXT,
		system_name TEXT NOT NULL
	)`)
	db.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create template"})
		return
	}

	c.FileAttachment(tmpPath, "template.db")
}

func UploadStudentDirectory(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}

	if file.Size > 50*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 50MB)"})
		return
	}

	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}
	defer src.Close()

	if err := os.MkdirAll("data", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create data directory"})
		return
	}

	tmpPath := studentDirectoryPath + ".tmp"
	dst, err := os.Create(tmpPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	if _, err := io.Copy(dst, src); err != nil {
		dst.Close()
		os.Remove(tmpPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	dst.Close()

	// Validate it's a valid SQLite DB with the expected schema
	testDB, err := sql.Open("sqlite3", tmpPath)
	if err != nil {
		os.Remove(tmpPath)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid database file"})
		return
	}
	var count int
	err = testDB.QueryRow("SELECT COUNT(*) FROM students").Scan(&count)
	testDB.Close()
	if err != nil {
		os.Remove(tmpPath)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid database: missing 'students' table")})
		return
	}

	if err := os.Rename(tmpPath, studentDirectoryPath); err != nil {
		os.Remove(tmpPath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "uploaded", "count": count})
}

func openDirectoryDB() (*sql.DB, error) {
	if _, err := os.Stat(studentDirectoryPath); os.IsNotExist(err) {
		return nil, err
	}
	return sql.Open("sqlite3", filepath.Clean(studentDirectoryPath)+"?mode=ro")
}
