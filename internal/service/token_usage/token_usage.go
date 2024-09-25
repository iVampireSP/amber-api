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

func (s *Service) IncrMonthTokenUsage(ctx context.Context, totalTokens int) {
	// 获取原来的值
	oldValue, _ := s.redis.Client.Get(ctx, s.cacheKey(KeyToken, s.today())).Int()
	totalTokens += oldValue

	status := s.redis.Client.Set(ctx, s.cacheKey(KeyToken, s.today()), totalTokens, 24*time.Hour)
	if status.Err() != nil {
		s.logger.Sugar.Errorf("IncrMonthTokenUsage: %v", status.Err())
	}
}

func (s *Service) IncrMonthToolCallTimes(ctx context.Context, totalTokens int) {
	oldValue, _ := s.redis.Client.Get(ctx, s.cacheKey(KeyTool, s.today())).Int()
	totalTokens += oldValue

	status := s.redis.Client.Set(ctx, s.cacheKey(KeyTool, s.today()), totalTokens, 24*time.Hour)
	if status.Err() != nil {
		s.logger.Sugar.Errorf("IncrMonthToolCallTimes: %v", status.Err())
	}
}

func (s *Service) GetMonthTokenUsage(ctx context.Context) int {
	val, err := s.redis.Client.Get(ctx, s.cacheKey(KeyToken, s.today())).Int()

	if err != nil {
		s.logger.Sugar.Errorf("GetMonthTokenUsage: %v", err)
	}

	return val
}

func (s *Service) GetLastMonthTokenUsage(ctx context.Context) int {
	val, err := s.redis.Client.Get(ctx, s.cacheKey(KeyToken, s.yesterday())).Int()
	if err != nil {
		s.logger.Sugar.Errorf("GetLastMonthTokenUsage: %v", err)
	}

	return val
}

func (s *Service) GetMonthToolCallTimes(ctx context.Context) int {
	val, err := s.redis.Client.Get(ctx, s.cacheKey(KeyTool, s.today())).Int()

	if err != nil {
		s.logger.Sugar.Errorf("GetMonthToolCallTimes: %v", err)
	}

	return val
}

func (s *Service) GetLastMonthToolCallTimes(ctx context.Context) int {
	val, err := s.redis.Client.Get(ctx, s.cacheKey(KeyTool, s.yesterday())).Int()
	if err != nil {
		s.logger.Sugar.Errorf("GetLastMonthToolCallTimes: %v", err)
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

func (s *Service) month() string {
	return time.Now().Format("01")
}

func (s *Service) lastMonth() string {
	return time.Now().AddDate(0, -1, 0).Format("01")
}
