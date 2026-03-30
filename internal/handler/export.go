package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"

	"github.com/keren/pao-event-panel/internal/database"
	"github.com/keren/pao-event-panel/internal/model"
)

func ExportRankings(c *gin.Context) {
	// Query students with house info (same logic as StudentRankings)
	var students []model.Student
	if err := database.DB.Preload("House").Order("lap_count DESC, name ASC").Find(&students).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query students"})
		return
	}

	// Query house rankings (same logic as HouseRankings)
	type houseResult struct {
		HouseID      uint `gorm:"column:house_id"`
		TotalLaps    int  `gorm:"column:total_laps"`
		StudentCount int  `gorm:"column:student_count"`
	}
	var houseResults []houseResult
	database.DB.Model(&model.Student{}).
		Select("house_id, COALESCE(SUM(lap_count), 0) as total_laps, COUNT(*) as student_count").
		Group("house_id").
		Order("total_laps DESC").
		Scan(&houseResults)

	var houses []model.House
	database.DB.Find(&houses)
	houseMap := make(map[uint]model.House)
	for _, h := range houses {
		houseMap[h.ID] = h
	}

	// Build workbook
	f := excelize.NewFile()
	defer f.Close()

	// Header style: bold + gray background
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Pattern: 1, Color: []string{"#E0E0E0"}},
	})

	// --- Sheet 1: Student Rankings ---
	sheet1 := "Student Rankings"
	f.SetSheetName("Sheet1", sheet1)

	headers1 := []string{"Rank", "Name", "House", "Laps"}
	widths1 := []float64{8, 25, 15, 10}
	for i, h := range headers1 {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s1", col)
		f.SetCellValue(sheet1, cell, h)
		f.SetCellStyle(sheet1, cell, cell, headerStyle)
		f.SetColWidth(sheet1, col, col, widths1[i])
	}

	currentRank := 1
	for i, s := range students {
		if i > 0 && s.LapCount < students[i-1].LapCount {
			currentRank = i + 1
		}
		row := i + 2
		f.SetCellValue(sheet1, fmt.Sprintf("A%d", row), currentRank)
		f.SetCellValue(sheet1, fmt.Sprintf("B%d", row), s.Name)
		f.SetCellValue(sheet1, fmt.Sprintf("C%d", row), s.House.Name)
		f.SetCellValue(sheet1, fmt.Sprintf("D%d", row), s.LapCount)
	}

	// --- Sheet 2: House Rankings ---
	sheet2 := "House Rankings"
	f.NewSheet(sheet2)

	headers2 := []string{"Rank", "House", "Total Laps", "Students"}
	widths2 := []float64{8, 15, 12, 12}
	for i, h := range headers2 {
		col := string(rune('A' + i))
		cell := fmt.Sprintf("%s1", col)
		f.SetCellValue(sheet2, cell, h)
		f.SetCellStyle(sheet2, cell, cell, headerStyle)
		f.SetColWidth(sheet2, col, col, widths2[i])
	}

	// Houses from aggregation results
	rank := 1
	for i, r := range houseResults {
		h := houseMap[r.HouseID]
		row := i + 2
		f.SetCellValue(sheet2, fmt.Sprintf("A%d", row), rank)
		f.SetCellValue(sheet2, fmt.Sprintf("B%d", row), h.Name)
		f.SetCellValue(sheet2, fmt.Sprintf("C%d", row), r.TotalLaps)
		f.SetCellValue(sheet2, fmt.Sprintf("D%d", row), r.StudentCount)
		rank++
	}

	// Include houses with 0 students
	for _, h := range houses {
		found := false
		for _, r := range houseResults {
			if r.HouseID == h.ID {
				found = true
				break
			}
		}
		if !found {
			row := rank + 1
			f.SetCellValue(sheet2, fmt.Sprintf("A%d", row), rank)
			f.SetCellValue(sheet2, fmt.Sprintf("B%d", row), h.Name)
			f.SetCellValue(sheet2, fmt.Sprintf("C%d", row), 0)
			f.SetCellValue(sheet2, fmt.Sprintf("D%d", row), 0)
			rank++
		}
	}

	// Write to buffer
	buf, err := f.WriteToBuffer()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate export"})
		return
	}

	filename := fmt.Sprintf("rankings_%s.xlsx", time.Now().Format("2006-01-02"))
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
