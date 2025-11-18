package market

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/faisal/crypto/backend/internal/config"
)

type Service struct {
	cfg    *config.Config
	client *http.Client
	cache  *cache.Cache
}

func NewService(cfg *config.Config) *Service {
	return &Service{
		cfg:    cfg,
		client: &http.Client{Timeout: 10 * time.Second},
		cache:  cache.New(time.Duration(cfg.CacheTTLSeconds)*time.Second, time.Minute),
	}
}

type CoinMarket struct {
	ID                       string    `json:"id"`
	Symbol                   string    `json:"symbol"`
	Name                     string    `json:"name"`
	CurrentPrice             float64   `json:"current_price"`
	PriceChangePercentage24h float64   `json:"price_change_percentage_24h"`
	SparklineIn7D            Sparkline `json:"sparkline_in_7d"`
}

type Sparkline struct {
	Price []float64 `json:"price"`
}

type CoinGeckoMarketResponse struct {
	ID                       string  `json:"id"`
	Symbol                   string  `json:"symbol"`
	Name                     string  `json:"name"`
	CurrentPrice             float64 `json:"current_price"`
	PriceChangePercentage24h float64 `json:"price_change_percentage_24h"`
	SparklineIn7D            struct {
		Price []float64 `json:"price"`
	} `json:"sparkline_in_7d"`
}

func (s *Service) GetTopMarketData() ([]CoinMarket, error) {
	if cached, found := s.cache.Get("market"); found {
		return cached.([]CoinMarket), nil
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/coins/markets", s.cfg.CoinGeckoBaseURL), nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Set("vs_currency", "usd")
	q.Set("order", "market_cap_desc")
	q.Set("per_page", fmt.Sprintf("%d", s.cfg.MarketDataLimit))
	q.Set("page", "1")
	q.Set("sparkline", "true")
	req.URL.RawQuery = q.Encode()
	if s.cfg.CoinGeckoAPIKey != "" {
		// CORRECT: This is for the free "Demo" plan
		req.Header.Set("x-cg-demo-api-key", s.cfg.CoinGeckoAPIKey)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Read error body for debugging
		bodyBytes := make([]byte, 512)
		n, _ := resp.Body.Read(bodyBytes)
		return nil, fmt.Errorf("coingecko returned status %d: %s", resp.StatusCode, string(bodyBytes[:n]))
	}

	var rawPayload []CoinGeckoMarketResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawPayload); err != nil {
		return nil, err
	}

	// Transform to our format
	payload := make([]CoinMarket, 0, len(rawPayload))
	for _, coin := range rawPayload {
		payload = append(payload, CoinMarket{
			ID:                       coin.ID,
			Symbol:                   coin.Symbol,
			Name:                     coin.Name,
			CurrentPrice:             coin.CurrentPrice,
			PriceChangePercentage24h: coin.PriceChangePercentage24h,
			SparklineIn7D: Sparkline{
				Price: coin.SparklineIn7D.Price,
			},
		})
	}

	s.cache.Set("market", payload, cache.DefaultExpiration)
	return payload, nil
}
