package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/yupi/airouter/go-backend/internal/config"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
)

const chatCacheKeyPrefix = "chat:cache:"

type ChatCacheService struct {
	pool    *redis.Pool
	enabled bool
	ttl     time.Duration
}

func NewChatCacheService(pool *redis.Pool, cfg *config.Config) *ChatCacheService {
	ttl := time.Duration(cfg.AICacheTTLSeconds) * time.Second
	if ttl <= 0 {
		ttl = 300 * time.Second
	}
	return &ChatCacheService{
		pool:    pool,
		enabled: cfg.AICacheEnabled,
		ttl:     ttl,
	}
}

func (s *ChatCacheService) Get(request dto.ChatRequest) (*dto.ChatResponse, bool) {
	if !s.enabled || request.EnableReasoning != nil && *request.EnableReasoning {
		return nil, false
	}
	conn := s.pool.Get()
	defer conn.Close()

	key := chatCacheKey(buildChatCacheSource(request))
	raw, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		return nil, false
	}
	var response dto.ChatResponse
	if err = json.Unmarshal(raw, &response); err != nil {
		return nil, false
	}
	return &response, true
}

func (s *ChatCacheService) Set(request dto.ChatRequest, response dto.ChatResponse) {
	if !s.enabled || request.EnableReasoning != nil && *request.EnableReasoning {
		return
	}
	payload, err := json.Marshal(response)
	if err != nil {
		return
	}
	conn := s.pool.Get()
	defer conn.Close()
	key := chatCacheKey(buildChatCacheSource(request))
	_, _ = conn.Do("SETEX", key, int(s.ttl/time.Second), payload)
}

func buildChatCacheSource(request dto.ChatRequest) string {
	req := request
	req.Stream = nil
	data, _ := json.Marshal(req)
	return string(data)
}

func chatCacheKey(source string) string {
	sum := sha256.Sum256([]byte(source))
	return chatCacheKeyPrefix + hex.EncodeToString(sum[:])
}
