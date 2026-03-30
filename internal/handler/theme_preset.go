package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/keren/pao-event-panel/internal/database"
	"github.com/keren/pao-event-panel/internal/model"
)

func ListThemePresets(c *gin.Context) {
	var presets []model.ThemePreset
	if err := database.DB.Order("created_at asc").Find(&presets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, presets)
}

type createPresetRequest struct {
	Name string `json:"name" binding:"required"`
}

func CreateThemePreset(c *gin.Context) {
	var req createPresetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Load current activity settings
	var settings model.ActivitySettings
	if err := database.DB.First(&settings, 1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "settings not found"})
		return
	}

	// Create preset from current settings
	preset := model.ThemePresetFromSettings(settings, req.Name)
	if err := database.DB.Create(&preset).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "preset name already exists"})
		return
	}

	c.JSON(http.StatusCreated, preset)
}

func DeleteThemePreset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	result := database.DB.Delete(&model.ThemePreset{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "preset not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func ApplyThemePreset(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// Load the preset
	var preset model.ThemePreset
	if err := database.DB.First(&preset, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "preset not found"})
		return
	}

	// Apply preset to activity settings
	var settings model.ActivitySettings
	if err := database.DB.First(&settings, 1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "settings not found"})
		return
	}

	updates := preset.ToSettingsUpdates()
	if err := database.DB.Model(&settings).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload and return updated settings
	database.DB.First(&settings, 1)
	c.JSON(http.StatusOK, settings)
}

type applyBuiltinRequest struct {
	Name string `json:"name" binding:"required"`
}

func ApplyBuiltinPreset(c *gin.Context) {
	var req applyBuiltinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	preset, ok := model.BuiltinPresets[req.Name]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown builtin preset"})
		return
	}

	// Apply preset to activity settings
	var settings model.ActivitySettings
	if err := database.DB.First(&settings, 1).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "settings not found"})
		return
	}

	updates := preset.ToSettingsUpdates()
	if err := database.DB.Model(&settings).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload and return updated settings
	database.DB.First(&settings, 1)
	c.JSON(http.StatusOK, settings)
}
