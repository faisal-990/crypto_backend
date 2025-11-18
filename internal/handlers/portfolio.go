package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/faisal/crypto/backend/internal/models"
	"github.com/faisal/crypto/backend/internal/services/portfolio"
)

type PortfolioHandler struct {
	service *portfolio.Service
}

func NewPortfolioHandler(service *portfolio.Service) *PortfolioHandler {
	return &PortfolioHandler{service: service}
}

func (h *PortfolioHandler) Register(router *gin.RouterGroup) {
	router.GET("/portfolio", h.getPortfolio)
	router.POST("/portfolio", h.createHolding)
	router.DELETE("/portfolio/:id", h.deleteHolding)

	router.GET("/portfolio/history", h.getHistory)
	router.POST("/portfolio/history", h.createSnapshot)
}

type createHoldingRequest struct {
	UserID string  `json:"userId" binding:"required"`
	CoinID string  `json:"coinId" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}

func (h *PortfolioHandler) getPortfolio(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		userID = "1"
	}
	data, total, err := h.service.GetHoldingsWithValue(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"totalValue": total,
		"holdings":   data,
	})
}

func (h *PortfolioHandler) createHolding(c *gin.Context) {
	var req createHoldingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	holding := models.Holding{
		UserID: req.UserID,
		CoinID: req.CoinID,
		Amount: req.Amount,
	}
	res, err := h.service.CreateHolding(c.Request.Context(), holding)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *PortfolioHandler) deleteHolding(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		userID = "1"
	}
	id := c.Param("id")
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.service.DeleteHolding(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *PortfolioHandler) getHistory(c *gin.Context) {
	userID := c.Query("userId")
	if userID == "" {
		userID = "1"
	}
	data, err := h.service.ListSnapshots(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

type createSnapshotRequest struct {
	UserID     string  `json:"userId" binding:"required"`
	TotalValue float64 `json:"totalValue" binding:"required"`
}

func (h *PortfolioHandler) createSnapshot(c *gin.Context) {
	var req createSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	snapshot := models.Snapshot{
		UserID:     req.UserID,
		TotalValue: req.TotalValue,
		Timestamp:  models.ToPrimitiveDateTime(time.Now()),
	}
	res, err := h.service.CreateSnapshot(c.Request.Context(), snapshot)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}
