package cache

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"smart-fish/back_end/config"

	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	ctx = context.Background()
)

// Connect 初始化 Redis 连接
func Connect() {
	cfg := config.AppConfig.Redis
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Printf("Warning: Redis connection failed: %v (cache disabled, falling back to DB)", err)
		RDB = nil
		return
	}
	log.Println("Redis connected successfully")
}

// Enabled 检查 Redis 是否可用
func Enabled() bool {
	return RDB != nil
}

// Get 从缓存获取值并反序列化到 dest
func Get(key string, dest interface{}) error {
	if !Enabled() {
		return redis.Nil
	}
	val, err := RDB.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Set 序列化并缓存值
func Set(key string, value interface{}, ttl time.Duration) error {
	if !Enabled() {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return RDB.Set(ctx, key, data, ttl).Err()
}

// Del 删除缓存
func Del(keys ...string) error {
	if !Enabled() {
		return nil
	}
	return RDB.Del(ctx, keys...).Err()
}

// DelByPattern 按模式删除缓存（谨慎使用）
func DelByPattern(pattern string) error {
	if !Enabled() {
		return nil
	}
	iter := RDB.Scan(ctx, 0, pattern, 100).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(keys) > 0 {
		return RDB.Del(ctx, keys...).Err()
	}
	return nil
}

// ==================== 缓存 Key 常量 ====================

const (
	// 天气缓存: weather:{location}
	KeyWeather = "weather:%s"
	WeatherTTL = 30 * time.Minute

	// 帖子列表缓存: posts:list:{tag}:{search}:{user}:{page}
	KeyPostsList = "posts:list:%s:%s:%s:%d"
	PostsListTTL = 5 * time.Minute

	// 帖子详情缓存: posts:detail:{id}
	KeyPostDetail = "posts:detail:%d"
	PostDetailTTL = 10 * time.Minute

	// 帖子点赞数缓存: posts:likes:{id}
	KeyPostLikes = "posts:likes:%d"
	PostLikesTTL = 3 * time.Minute

	// 系统概览缓存
	KeySummary = "summary"
	SummaryTTL = 5 * time.Minute

	// 区域环境数据缓存
	KeyRegionEnv     = "region:env"
	RegionEnvTTL     = 10 * time.Minute
	KeyRegionEnvHist = "region:env:hist:%d:%d"
	RegionEnvHistTTL = 10 * time.Minute

	// 垂钓统计缓存: fishing:stats:{user_id}
	KeyFishingStats = "fishing:stats:%d"
	FishingStatsTTL = 5 * time.Minute

	// 水域热门缓存
	KeyPopularSpots = "spots:popular"
	PopularSpotsTTL = 15 * time.Minute

	// 水域历史/环境数据缓存
	KeySpotHistorical = "spots:hist:%d"
	SpotHistoricalTTL = 10 * time.Minute
	KeySpotEnv        = "spots:env:%d"
	SpotEnvTTL        = 10 * time.Minute

	// 垂钓建议缓存
	KeySuggestionsLatest = "suggestions:latest"
	SuggestionsLatestTTL = 10 * time.Minute
)
