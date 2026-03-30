package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var allowedMimeTypes = map[string]bool{
	"image/png":     true,
	"image/jpeg":    true,
	"image/webp":    true,
	"image/svg+xml": true,
}

func UploadLogo(c *gin.Context) {
	uploadFile(c, "logo")
}

func UploadBackground(c *gin.Context) {
	uploadFile(c, "bg")
}

func uploadFile(c *gin.Context, prefix string) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}

	if file.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large (max 5MB)"})
		return
	}

	mimeType := file.Header.Get("Content-Type")
	if !allowedMimeTypes[mimeType] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		switch mimeType {
		case "image/png":
			ext = ".png"
		case "image/jpeg":
			ext = ".jpg"
		case "image/webp":
			ext = ".webp"
		case "image/svg+xml":
			ext = ".svg"
		}
	}
	ext = strings.ToLower(ext)

	filename := fmt.Sprintf("%s-%d%s", prefix, time.Now().Unix(), ext)
	uploadPath := filepath.Join("uploads", filename)

	if err := os.MkdirAll("uploads", 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"path": "/" + uploadPath})
}
