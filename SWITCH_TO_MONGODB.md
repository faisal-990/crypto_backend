# How to Switch from In-Memory to MongoDB

## Step 1: Uncomment MongoDB Service Function

In `internal/services/portfolio/service.go`, uncomment the `NewServiceWithMongo` function:

```go
func NewServiceWithMongo(cfg *config.Config, client *mongo.Client) *Service {
	repo := repository.NewMongoPortfolioRepository(client.Database(cfg.MongoDBName))
	return &Service{
		cfg:           cfg,
		repo:          repo,
		marketService: market.NewService(cfg),
	}
}
```

## Step 2: Update main.go

Replace the in-memory setup with MongoDB setup:

**Before (In-Memory):**
```go
// Using in-memory storage for development
// TODO: Switch to MongoDB when connection is ready
portfolioService := portfolio.NewService(cfg)
portfolioHandler := handlers.NewPortfolioHandler(portfolioService)
portfolioHandler.Register(api)
```

**After (MongoDB):**
```go
// Connect to MongoDB
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

mongoClient, err := db.Connect(ctx, cfg)
if err != nil {
	log.Fatalf("connect mongo: %v", err)
}
defer func() {
	_ = mongoClient.Disconnect(context.Background())
}()

// Using MongoDB storage
portfolioService := portfolio.NewServiceWithMongo(cfg, mongoClient)
portfolioHandler := handlers.NewPortfolioHandler(portfolioService)
portfolioHandler.Register(api)
```

## Step 3: Add db import back

Add this import to `cmd/server/main.go`:
```go
"github.com/faisal/crypto/backend/internal/db"
```

## That's it!

The rest of the code doesn't need to change because both repositories implement the same interface.

