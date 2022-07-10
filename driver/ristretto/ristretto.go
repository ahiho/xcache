package ristretto

import (
	"time"

	"github.com/dgraph-io/ristretto"
)

type RistrettoDriver struct {
	cache *ristretto.Cache
}

func NewRistrettoDriver(cache *ristretto.Cache) *RistrettoDriver {
	return &RistrettoDriver{
		cache: cache,
	}
}

func (r *RistrettoDriver) Set(k string, v string, d time.Duration) error {
	r.cache.SetWithTTL(k, v, int64(len(v)), d)
	return nil
}

func (r *RistrettoDriver) Get(k string) (*string, error) {
	v, b := r.cache.Get(k)
	if !b {
		return nil, nil
	}
	s := v.(string)
	return &s, nil
}

func (r *RistrettoDriver) Del(k string) error {
	r.cache.Del(k)
	return nil
}

func (r *RistrettoDriver) MultiSet(m map[string]string, d time.Duration) error {
	for k, v := range m {
		r.cache.SetWithTTL(k, v, int64(len(v)), d)
	}
	return nil
}

func (r *RistrettoDriver) MultiGet(keys []string) (map[string]string, error) {
	m := map[string]string{}
	for _, k := range keys {
		v, b := r.cache.Get(k)
		if b {
			m[k] = v.(string)
		}
	}
	return m, nil
}

func (r *RistrettoDriver) MultiDel(keys []string) error {
	for _, k := range keys {
		r.cache.Del(k)
	}
	return nil
}
