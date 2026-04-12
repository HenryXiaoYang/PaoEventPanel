package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/keren/pao-event-panel/internal/database"
	"github.com/keren/pao-event-panel/internal/model"
)

func GetActivity(c *gin.Context) {
	var settings model.ActivitySettings
	if err := database.DB.First(&settings, 1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "settings not found"})
		return
	}
	c.JSON(http.StatusOK, settings)
}

func UpdateActivity(c *gin.Context) {
	var settings model.ActivitySettings
	if err := database.DB.First(&settings, 1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "settings not found"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	allowedFields := map[string]bool{
		"activity_name":          true,
		"logo_path":              true,
		"bg_mode":                true,
		"bg_color":               true,
		"bg_gradient_start":      true,
		"bg_gradient_end":        true,
		"bg_gradient_angle":      true,
		"bg_image_path":          true,
		"bg_blur":                true,
		"dark_bg_mode":           true,
		"dark_bg_color":          true,
		"dark_bg_gradient_start": true,
		"dark_bg_gradient_end":   true,
		"dark_bg_gradient_angle": true,
		"dark_bg_image_path":     true,
		"dark_bg_blur":           true,
		"primary_color":          true,
		"secondary_color":        true,
		"accent_color":           true,
		"muted_color":            true,
		"dark_primary_color":     true,
		"dark_secondary_color":   true,
		"dark_accent_color":      true,
		"dark_muted_color":       true,
		"card_opacity":           true,
		"card_blur":              true,
		"dark_card_opacity":      true,
		"dark_card_blur":         true,
		"weather_location":       true,
	}

	filteredUpdates := make(map[string]interface{})
	for k, v := range updates {
		if allowedFields[k] {
			filteredUpdates[k] = v
		}
	}

	if err := database.DB.Model(&settings).Updates(filteredUpdates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	database.DB.First(&settings, 1)
	c.JSON(http.StatusOK, settings)
}
