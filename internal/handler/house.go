package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/keren/pao-event-panel/internal/database"
	"github.com/keren/pao-event-panel/internal/model"
)

func ListHouses(c *gin.Context) {
	var houses []model.House
	database.DB.Find(&houses)
	c.JSON(http.StatusOK, houses)
}

func UpdateHouseColor(c *gin.Context) {
	id := c.Param("id")
	var body struct {
		Color string `json:"color" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "color is required"})
		return
	}

	var house model.House
	if err := database.DB.First(&house, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "house not found"})
		return
	}

	house.Color = body.Color
	database.DB.Save(&house)

	c.JSON(http.StatusOK, house)
}
