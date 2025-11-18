package portfolio

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/faisal/crypto/backend/internal/config"
	"github.com/faisal/crypto/backend/internal/models"
	"github.com/faisal/crypto/backend/internal/repository"
	"github.com/faisal/crypto/backend/internal/services/market"
)

type Service struct {
	cfg           *config.Config
	repo          repository.PortfolioRepository
	marketService *market.Service
}

// NewService creates a portfolio service with in-memory storage (for development)
func NewService(cfg *config.Config) *Service {
	repo := repository.NewMemoryPortfolioRepository()
	return &Service{
		cfg:           cfg,
		repo:          repo,
		marketService: market.NewService(cfg),
	}
}

// NewServiceWithMongo creates a portfolio service with MongoDB (for production)
// Use this when MongoDB connection is ready
func NewServiceWithMongo(cfg *config.Config, client *mongo.Client) *Service {
	repo := repository.NewMongoPortfolioRepository(client.Database(cfg.MongoDBName))
	return &Service{
		cfg:           cfg,
		repo:          repo,
		marketService: market.NewService(cfg),
	}
}

func (s *Service) ListHoldings(ctx context.Context, userID string) ([]models.Holding, error) {
	return s.repo.ListHoldings(ctx, userID)
}

type HoldingWithValue struct {
	models.Holding
	CurrentPrice float64 `json:"currentPrice"`
	CurrentValue float64 `json:"currentValue"`
}

func (s *Service) GetHoldingsWithValue(ctx context.Context, userID string) ([]HoldingWithValue, float64, error) {
	holdings, err := s.ListHoldings(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	marketData, err := s.marketService.GetTopMarketData()
	if err != nil {
		return nil, 0, err
	}
	priceIndex := make(map[string]market.CoinMarket, len(marketData))
	for _, coin := range marketData {
		priceIndex[coin.ID] = coin
	}

	var enriched []HoldingWithValue
	var total float64
	for _, holding := range holdings {
		coin, ok := priceIndex[holding.CoinID]
		if !ok {
			continue
		}
		value := holding.Amount * coin.CurrentPrice
		enriched = append(enriched, HoldingWithValue{
			Holding:      holding,
			CurrentPrice: coin.CurrentPrice,
			CurrentValue: value,
		})
		total += value
	}
	return enriched, total, nil
}

func (s *Service) CreateHolding(ctx context.Context, holding models.Holding) (*models.Holding, error) {
	if holding.UserID == "" || holding.CoinID == "" || holding.Amount <= 0 {
		return nil, errors.New("invalid holding payload")
	}
	return s.repo.CreateHolding(ctx, holding)
}

func (s *Service) DeleteHolding(ctx context.Context, id string, userID string) error {
	return s.repo.DeleteHolding(ctx, id, userID)
}

func (s *Service) ListSnapshots(ctx context.Context, userID string) ([]models.Snapshot, error) {
	return s.repo.ListSnapshots(ctx, userID)
}

func (s *Service) CreateSnapshot(ctx context.Context, snapshot models.Snapshot) (*models.Snapshot, error) {
	if snapshot.UserID == "" || snapshot.TotalValue < 0 {
		return nil, errors.New("invalid snapshot payload")
	}
	if snapshot.Timestamp == 0 {
		snapshot.Timestamp = models.ToPrimitiveDateTime(time.Now())
	}
	return s.repo.CreateSnapshot(ctx, snapshot)
}
