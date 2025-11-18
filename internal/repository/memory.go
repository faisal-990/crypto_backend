package repository

import (
	"context"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/faisal/crypto/backend/internal/models"
)

// MemoryPortfolioRepository is an in-memory implementation for development/testing
type MemoryPortfolioRepository struct {
	holdings map[string]models.Holding // key: holding ID
	snapshots map[string]models.Snapshot // key: snapshot ID
	mu       sync.RWMutex
}

func NewMemoryPortfolioRepository() *MemoryPortfolioRepository {
	return &MemoryPortfolioRepository{
		holdings:  make(map[string]models.Holding),
		snapshots: make(map[string]models.Snapshot),
	}
}

func (r *MemoryPortfolioRepository) ListHoldings(ctx context.Context, userID string) ([]models.Holding, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []models.Holding
	for _, holding := range r.holdings {
		if holding.UserID == userID {
			result = append(result, holding)
		}
	}
	return result, nil
}

func (r *MemoryPortfolioRepository) CreateHolding(ctx context.Context, holding models.Holding) (*models.Holding, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate ID if not set
	if holding.ID.IsZero() {
		holding.ID = primitive.NewObjectID()
	}

	idStr := holding.ID.Hex()
	r.holdings[idStr] = holding
	return &holding, nil
}

func (r *MemoryPortfolioRepository) DeleteHolding(ctx context.Context, id string, userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	holding, exists := r.holdings[id]
	if !exists {
		return nil // Already deleted or doesn't exist
	}

	if holding.UserID != userID {
		return nil // Not the user's holding, but don't error
	}

	delete(r.holdings, id)
	return nil
}

func (r *MemoryPortfolioRepository) ListSnapshots(ctx context.Context, userID string) ([]models.Snapshot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []models.Snapshot
	for _, snapshot := range r.snapshots {
		if snapshot.UserID == userID {
			result = append(result, snapshot)
		}
	}
	return result, nil
}

func (r *MemoryPortfolioRepository) CreateSnapshot(ctx context.Context, snapshot models.Snapshot) (*models.Snapshot, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Generate ID if not set
	if snapshot.ID.IsZero() {
		snapshot.ID = primitive.NewObjectID()
	}

	// Set timestamp if not set
	if snapshot.Timestamp == 0 {
		snapshot.Timestamp = models.ToPrimitiveDateTime(time.Now())
	}

	idStr := snapshot.ID.Hex()
	r.snapshots[idStr] = snapshot
	return &snapshot, nil
}

