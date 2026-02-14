package service

import (
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	rateLimitAPIKeyPrefix = "rate_limit:api_key:"
	rateLimitIPPrefix     = "rate_limit:ip:"
)

var acquirePermitScript = redis.NewScript(1, `
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("PEXPIRE", KEYS[1], ARGV[2])
end
if current <= tonumber(ARGV[1]) then
  return 1
end
return 0
`)

type RateLimitService struct {
	pool *redis.Pool
}

func NewRateLimitService(pool *redis.Pool) *RateLimitService {
	return &RateLimitService{pool: pool}
}

func (s *RateLimitService) TryAcquire(key string, limit int, window time.Duration) (bool, error) {
	conn := s.pool.Get()
	defer conn.Close()

	windowMillis := int(window / time.Millisecond)
	allowed, err := redis.Int(acquirePermitScript.Do(conn, key, limit, windowMillis))
	if err != nil {
		return false, err
	}
	return allowed == 1, nil
}

func BuildAPIKeyRateLimitKey(apiKey string) string {
	return rateLimitAPIKeyPrefix + apiKey
}

func BuildIPRateLimitKey(ip string) string {
	return rateLimitIPPrefix + ip
}
