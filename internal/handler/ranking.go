package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/keren/pao-event-panel/internal/database"
	"github.com/keren/pao-event-panel/internal/model"
)

type HouseRanking struct {
	HouseID      uint   `json:"house_id"`
	Name         string `json:"name"`
	Code         string `json:"code"`
	Color        string `json:"color"`
	TotalLaps    int    `json:"total_laps"`
	StudentCount int    `json:"student_count"`
}

type RankedStudent struct {
	Rank     int         `json:"rank"`
	ID       uint        `json:"id"`
	Name     string      `json:"name"`
	House    model.House `json:"house"`
	LapCount int         `json:"lap_count"`
}

type Stats struct {
	TotalLaps         int `json:"total_laps"`
	TotalParticipants int `json:"total_participants"`
}

func HouseRankings(c *gin.Context) {
	type result struct {
		HouseID      uint `gorm:"column:house_id"`
		TotalLaps    int  `gorm:"column:total_laps"`
		StudentCount int  `gorm:"column:student_count"`
	}
	var results []result

	database.DB.Model(&model.Student{}).
		Select("house_id, COALESCE(SUM(lap_count), 0) as total_laps, COUNT(*) as student_count").
		Group("house_id").
		Order("total_laps DESC").
		Scan(&results)

	var houses []model.House
	database.DB.Find(&houses)
	houseMap := make(map[uint]model.House)
	for _, h := range houses {
		houseMap[h.ID] = h
	}

	rankings := make([]HouseRanking, 0, len(results))
	for _, r := range results {
		h := houseMap[r.HouseID]
		rankings = append(rankings, HouseRanking{
			HouseID:      r.HouseID,
			Name:         h.Name,
			Code:         h.Code,
			Color:        h.Color,
			TotalLaps:    r.TotalLaps,
			StudentCount: r.StudentCount,
		})
	}

	// Include houses with 0 students
	for _, h := range houses {
		found := false
		for _, r := range rankings {
			if r.HouseID == h.ID {
				found = true
				break
			}
		}
		if !found {
			rankings = append(rankings, HouseRanking{
				HouseID:      h.ID,
				Name:         h.Name,
				Code:         h.Code,
				Color:        h.Color,
				TotalLaps:    0,
				StudentCount: 0,
			})
		}
	}

	c.JSON(http.StatusOK, rankings)
}

func StudentRankings(c *gin.Context) {
	var students []model.Student
	database.DB.Preload("House").Order("lap_count DESC, name ASC").Find(&students)

	rankings := make([]RankedStudent, len(students))
	currentRank := 1
	for i, s := range students {
		if i > 0 && s.LapCount < students[i-1].LapCount {
			currentRank = i + 1
		}
		rankings[i] = RankedStudent{
			Rank:     currentRank,
			ID:       s.ID,
			Name:     s.Name,
			House:    s.House,
			LapCount: s.LapCount,
		}
	}

	c.JSON(http.StatusOK, rankings)
}

func RankingStats(c *gin.Context) {
	var stats Stats

	database.DB.Model(&model.Student{}).
		Select("COALESCE(SUM(lap_count), 0) as total_laps, COUNT(*) as total_participants").
		Scan(&stats)

	c.JSON(http.StatusOK, stats)
}
