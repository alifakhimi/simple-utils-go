package cache

import (
	"context"
	"errors"
	"time"

	"github.com/alifakhimi/simple-service-go/utils/templates"
	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

// Error block
var (
	ErrCacheInfo = errors.New("invalid cache info")
)

// instance
var c *Cache

type Cache struct {
	Active       bool           `json:"active,omitempty"`
	RedisOptions *redis.Options `json:"redis_options,omitempty" mapstructure:"redis_options"`
	RedisClient  *redis.Client  `json:"redis_client,omitempty" mapstructure:"redis_client"`
	ExpireTime   time.Duration  `json:"expire_time,omitempty" mapstructure:"expire_time"`
}

func Init(cache *Cache) (err error) {
	c = cache

	if cache == nil || cache.RedisOptions == nil {
		return ErrCacheInfo
	}

	if !cache.Active {
		return nil
	}

	cache.RedisClient = redis.NewClient(cache.RedisOptions)

	return nil
}

func GetInstanse() *Cache {
	// if c == nil {
	// 	c = Init()
	// }

	return c
}

func (c *Cache) Get(ctx context.Context, key string, val interface{}) (err error) {
	if err = c.RedisClient.Get(ctx, key).Scan(val); err != nil {
		return
	}

	return
}

func (c *Cache) Set(ctx context.Context, key string, val interface{}) (err error) {
	if err = c.RedisClient.Set(ctx, key, val, c.ExpireTime).Err(); err != nil {
		return
	}

	return
}

func (c *Cache) SetResponse(ctx echo.Context, data, meta interface{}) (err error) {
	var (
		response   = new(templates.ResponseTemplate)
		requestKey = ctx.Request().RequestURI
	)

	response.Data = data
	response.Meta = meta

	if err = c.Set(ctx.Request().Context(), requestKey, response); err != nil {
		return
	}

	return nil
}

// func (c *Cache) Get(key string, val interface{}) (err error) {
// 	if err = c.RedisClient.Get(context.Background(), key).Scan(val); err != nil {
// 		return
// 	}

// 	return
// }

// func (c *Cache) Set(key string, val interface{}) (err error) {
// 	if err = c.RedisClient.Set(context.Background(), key, val, c.ExpireTime).Err(); err != nil {
// 		return
// 	}

// 	return
// }

// func (c *Cache) ReturnResponse(r *http.Request, data, meta interface{}) (err error) {
// 	var (
// 		response   = new(httpResponse.Response)
// 		requestKey = r.URL.String()
// 	)

// 	response.Data = data
// 	response.Meta = meta

// 	if err = c.Set(requestKey, response); err != nil {
// 		return
// 	}

// 	return nil
// }
