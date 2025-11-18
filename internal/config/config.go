package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Port        string
	MongoURI    string
	MongoDBName string

	CoinGeckoBaseURL string
	CoinGeckoAPIKey  string

	CacheTTLSeconds int
	MarketDataLimit int // Number of coins to fetch (for dev/testing)
	AllowedOrigins  []string
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	origin := []string{
		"http://localhost:5173",
		"http://localhost:5174",
	}

	cacheTTL := getEnvAsInt("CACHE_TTL_SECONDS", 300)  // Default 5 minutes for dev
	marketLimit := getEnvAsInt("MARKET_DATA_LIMIT", 5) // Default 5 coins for dev
	cfg := &Config{
		Port:             getEnv("PORT", "8080"),
		MongoURI:         getEnv("MONGO_URI", "mongodb://localhost:27017/crypto"),
		MongoDBName:      getEnv("MONGO_DB_NAME", "crypto"),
		CoinGeckoBaseURL: getEnv("COINGECKO_BASE_URL", "https://api.coingecko.com/api/v3"),
		CoinGeckoAPIKey:  getEnv("COINGECKO_API_KEY", ""),
		CacheTTLSeconds:  cacheTTL,
		MarketDataLimit:  marketLimit,
		AllowedOrigins:   origin,
	}
	return cfg, nil
}

func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return fallback
}

func CORSMiddleware(origins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(origins))
	for _, origin := range origins {
		allowed[strings.TrimSpace(origin)] = struct{}{}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if _, ok := allowed[origin]; ok {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
