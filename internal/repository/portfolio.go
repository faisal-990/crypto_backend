package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/faisal/crypto/backend/internal/models"
)

type PortfolioRepository interface {
	ListHoldings(ctx context.Context, userID string) ([]models.Holding, error)
	CreateHolding(ctx context.Context, holding models.Holding) (*models.Holding, error)
	DeleteHolding(ctx context.Context, id string, userID string) error
	ListSnapshots(ctx context.Context, userID string) ([]models.Snapshot, error)
	CreateSnapshot(ctx context.Context, snapshot models.Snapshot) (*models.Snapshot, error)
}

type MongoPortfolioRepository struct {
	holdings *mongo.Collection
	history  *mongo.Collection
}

func NewMongoPortfolioRepository(db *mongo.Database) *MongoPortfolioRepository {
	return &MongoPortfolioRepository{
		holdings: db.Collection("holdings"),
		history:  db.Collection("snapshots"),
	}
}

func (r *MongoPortfolioRepository) ListHoldings(ctx context.Context, userID string) ([]models.Holding, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cur, err := r.holdings.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var holdings []models.Holding
	if err := cur.All(ctx, &holdings); err != nil {
		return nil, err
	}
	return holdings, nil
}

func (r *MongoPortfolioRepository) CreateHolding(ctx context.Context, holding models.Holding) (*models.Holding, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.holdings.InsertOne(ctx, holding)
	if err != nil {
		return nil, err
	}
	holding.ID = res.InsertedID.(primitive.ObjectID)
	return &holding, nil
}

func (r *MongoPortfolioRepository) DeleteHolding(ctx context.Context, id string, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.holdings.DeleteOne(ctx, bson.M{"_id": objID, "user_id": userID})
	return err
}

func (r *MongoPortfolioRepository) ListSnapshots(ctx context.Context, userID string) ([]models.Snapshot, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cur, err := r.history.Find(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var snapshots []models.Snapshot
	if err := cur.All(ctx, &snapshots); err != nil {
		return nil, err
	}
	return snapshots, nil
}

func (r *MongoPortfolioRepository) CreateSnapshot(ctx context.Context, snapshot models.Snapshot) (*models.Snapshot, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	res, err := r.history.InsertOne(ctx, snapshot)
	if err != nil {
		return nil, err
	}
	snapshot.ID = res.InsertedID.(primitive.ObjectID)
	return &snapshot, nil
}
