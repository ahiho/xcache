package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisDriver struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisDriver(client *redis.Client) *RedisDriver {
	return &RedisDriver{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisDriver) Set(k string, v string, d time.Duration) error {
	return r.client.Set(r.ctx, k, v, d).Err()
}

func (r *RedisDriver) Get(k string) (*string, error) {
	s, e := r.client.Get(r.ctx, k).Result()
	if e == redis.Nil {
		return nil, nil
	}
	if e != nil {
		return nil, e
	}
	return &s, nil
}

func (r *RedisDriver) Del(k string) error {
	return r.client.Del(r.ctx, k).Err()
}

func (r *RedisDriver) MultiSet(m map[string]string, d time.Duration) error {
	p := r.client.Pipeline()
	for k, v := range m {
		p.Set(r.ctx, k, v, d)
	}
	_, e := p.Exec(r.ctx)
	return e
}

func (r *RedisDriver) MultiGet(keys []string) (map[string]string, error) {
	c := r.client.MGet(r.ctx, keys...)
	rs, e := c.Result()
	if e != nil {
		return nil, e
	}
	ret := map[string]string{}
	for i, k := range keys {
		val := rs[i]
		if val != nil {
			s := val.(string)
			ret[k] = s
		}

	}
	return ret, nil
}

func (r *RedisDriver) MultiDel(keys []string) error {
	return r.client.Del(r.ctx, keys...).Err()
}
