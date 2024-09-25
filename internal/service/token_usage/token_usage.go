package token_usage

import (
	"context"
	"fmt"
	"time"
)

type Key string

const (
	KeyToken Key = "token"
	KeyTool  Key = "tool"
)

func (s *Service) IncrTodayTokenUsage(ctx context.Context, totalTokens int) {
	// 获取原来的值
	oldValue, _ := s.redis.Client.Get(ctx, s.cacheKey(KeyToken, s.today())).Int()
	totalTokens += oldValue

	status := s.redis.Client.Set(ctx, s.cacheKey(KeyToken, s.today()), totalTokens, 24*time.Hour)
	if status.Err() != nil {
		s.logger.Sugar.Errorf("IncrTodayTokenUsage: %v", status.Err())
	}
}

func (s *Service) IncrTodayToolCallTimes(ctx context.Context, totalTokens int) {
	oldValue, _ := s.redis.Client.Get(ctx, s.cacheKey(KeyTool, s.today())).Int()
	totalTokens += oldValue

	status := s.redis.Client.Set(ctx, s.cacheKey(KeyTool, s.today()), totalTokens, 24*time.Hour)
	if status.Err() != nil {
		s.logger.Sugar.Errorf("IncrTodayToolCallTimes: %v", status.Err())
	}
}

func (s *Service) GetTodayTokenUsage(ctx context.Context) int {
	val, err := s.redis.Client.Get(ctx, s.cacheKey(KeyToken, s.today())).Int()

	if err != nil {
		s.logger.Sugar.Errorf("GetTodayTokenUsage: %v", err)
	}

	return val
}

func (s *Service) GetYesterdayTokenUsage(ctx context.Context) int {
	val, err := s.redis.Client.Get(ctx, s.cacheKey(KeyToken, s.yesterday())).Int()
	if err != nil {
		s.logger.Sugar.Errorf("GetYesterdayTokenUsage: %v", err)
	}

	return val
}

func (s *Service) GetTodayToolCallTimes(ctx context.Context) int {
	val, err := s.redis.Client.Get(ctx, s.cacheKey(KeyTool, s.today())).Int()

	if err != nil {
		s.logger.Sugar.Errorf("GetTodayToolCallTimes: %v", err)
	}

	return val
}

func (s *Service) GetYesterdayToolCallTimes(ctx context.Context) int {
	val, err := s.redis.Client.Get(ctx, s.cacheKey(KeyTool, s.yesterday())).Int()
	if err != nil {
		s.logger.Sugar.Errorf("GetYesterdayToolCallTimes: %v", err)
	}

	return val
}

func (s *Service) cacheKey(key Key, day string) string {
	return fmt.Sprintf("token_usage:%s:%s", key, day)
}

func (s *Service) today() string {
	return time.Now().Format("02")
}

func (s *Service) yesterday() string {
	return time.Now().AddDate(0, 0, -1).Format("02")
}
