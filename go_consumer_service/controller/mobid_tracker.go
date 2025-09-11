package controller

import (
	"context"
	"go_consumer_service/service"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type MobidTrackerController struct {
	// inject the service depencies
	mobidTrackerService service.MobidTrackerService
}

func NewMobidTrackerController(mobidTrackerService service.MobidTrackerService) *MobidTrackerController {
	return &MobidTrackerController{
		mobidTrackerService: mobidTrackerService,
	}
}

func (m *MobidTrackerController) RegisterRoutes(router *gin.Engine) {
	router.GET("/mobid-trackers", m.getMobidTracker)
}

func (m *MobidTrackerController) getMobidTracker(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	// Parse pagination params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit
	mobidTrackers, total, err := m.mobidTrackerService.GetAllMobidTracker(ctx, limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	// Response with metadata
	c.JSON(http.StatusOK, gin.H{
		"data":       mobidTrackers,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": int(math.Ceil(float64(total) / float64(limit))),
	})
}
