package service

import (
	"log"
	"strings"

	"github.com/gomodule/redigo/redis"
)

const blacklistRedisKey = "blacklist:ip"

type BlacklistService struct {
	pool *redis.Pool
}

func NewBlacklistService(pool *redis.Pool) *BlacklistService {
	return &BlacklistService{pool: pool}
}

func (s *BlacklistService) IsBlocked(ip string) bool {
	trimmed := strings.TrimSpace(ip)
	if trimmed == "" {
		return false
	}
	conn := s.pool.Get()
	defer conn.Close()

	blocked, err := redis.Bool(conn.Do("SISMEMBER", blacklistRedisKey, trimmed))
	if err != nil {
		log.Printf("check blacklist failed: ip=%s err=%v", trimmed, err)
		return false
	}
	return blocked
}

func (s *BlacklistService) AddToBlacklist(ip, reason string) error {
	trimmed := strings.TrimSpace(ip)
	if trimmed == "" {
		return nil
	}
	conn := s.pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SADD", blacklistRedisKey, trimmed); err != nil {
		return err
	}
	log.Printf("ip added to blacklist: ip=%s reason=%s", trimmed, reason)
	return nil
}

func (s *BlacklistService) RemoveFromBlacklist(ip string) error {
	trimmed := strings.TrimSpace(ip)
	if trimmed == "" {
		return nil
	}
	conn := s.pool.Get()
	defer conn.Close()

	_, err := conn.Do("SREM", blacklistRedisKey, trimmed)
	if err == nil {
		log.Printf("ip removed from blacklist: ip=%s", trimmed)
	}
	return err
}

func (s *BlacklistService) ListAll() ([]string, error) {
	conn := s.pool.Get()
	defer conn.Close()

	return redis.Strings(conn.Do("SMEMBERS", blacklistRedisKey))
}

func (s *BlacklistService) Count() (int64, error) {
	conn := s.pool.Get()
	defer conn.Close()

	return redis.Int64(conn.Do("SCARD", blacklistRedisKey))
}
