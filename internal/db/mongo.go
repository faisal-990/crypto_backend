package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/faisal/crypto/backend/internal/config"
)

func Connect(ctx context.Context, cfg *config.Config) (*mongo.Client, error) {
	// Set connection options with longer timeouts
	opts := options.Client().ApplyURI(cfg.MongoURI)
	opts.SetServerSelectionTimeout(30 * time.Second)
	opts.SetConnectTimeout(30 * time.Second)
	opts.SetSocketTimeout(30 * time.Second)
	
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	
	// Ping to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}
	
	return client, nil
}
