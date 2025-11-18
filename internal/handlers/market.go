package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/faisal/crypto/backend/internal/services/market"
)

type MarketHandler struct {
	service *market.Service
}

func NewMarketHandler(service *market.Service) *MarketHandler {
	return &MarketHandler{service: service}
}

func (h *MarketHandler) Register(router *gin.RouterGroup) {
	router.GET("/market", h.getMarket)
}

func (h *MarketHandler) getMarket(c *gin.Context) {
	data, err := h.service.GetTopMarketData()
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
