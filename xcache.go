package xcache

import (
	"errors"
	"strconv"
	"time"

	"github.com/vmihailenco/msgpack/v5"
)

type options struct {
	expiration time.Duration
}

type Option func(o *options)

func WithExpiration(expiration time.Duration) Option {
	return func(o *options) {
		o.expiration = expiration
	}
}

var (
	ErrInvalidDuration = errors.New("invalid duration")
)

type Driver interface {
	Set(k string, v string, d time.Duration) error
	Get(k string) (*string, error)
	Del(k string) error
	MultiSet(m map[string]string, d time.Duration) error
	MultiGet(k []string) (map[string]string, error)
	MultiDel(k []string) error
}

type XCache struct {
	driver     Driver
	expiration time.Duration
}

func NewCache(driver Driver, ops ...Option) (*XCache, error) {
	op := options{
		expiration: time.Hour * 24,
	}
	for _, o := range ops {
		o(&op)
	}
	if op.expiration <= 0 {
		return nil, ErrInvalidDuration
	}
	return &XCache{
		driver:     driver,
		expiration: op.expiration,
	}, nil
}

/* get cache functions */

func (c *XCache) GetString(k string) (*string, error) {
	return c.driver.Get(k)
}

func (c *XCache) GetMultipleString(keys ...string) (map[string]string, error) {
	return c.driver.MultiGet(keys)
}

func (c *XCache) GetBool(k string) (*bool, error) {
	v, e := c.driver.Get(k)
	if e != nil {
		return nil, e
	}
	if v == nil {
		return nil, nil
	}
	var b bool
	s := *v
	b, e = strconv.ParseBool(s)
	if e != nil {
		return nil, e
	}
	return &b, nil
}

func (c *XCache) GetInt(k string) (*int, error) {
	v, e := c.driver.Get(k)
	if e != nil {
		return nil, e
	}
	if v == nil {
		return nil, nil
	}
	s := *v
	i, e := strconv.Atoi(s)
	if e != nil {
		return nil, e
	}
	return &i, nil
}

func (c *XCache) GetMultipleInt(keys ...string) (map[string]int, error) {
	rs, e := c.driver.MultiGet(keys)
	if e != nil {
		return nil, e
	}
	m := map[string]int{}
	for k, v := range rs {
		i, e := strconv.Atoi(v)
		if e != nil {
			return nil, e
		}
		m[k] = i
	}
	return m, nil
}

func (c *XCache) GetInt64(k string) (*int64, error) {
	v, e := c.driver.Get(k)
	if e != nil {
		return nil, e
	}
	if v == nil {
		return nil, nil
	}
	s := *v
	i, e := strconv.ParseInt(s, 10, 64)
	if e != nil {
		return nil, e
	}
	return &i, nil
}

func (c *XCache) GetMultipleInt64(keys ...string) (map[string]int64, error) {
	rs, e := c.driver.MultiGet(keys)
	if e != nil {
		return nil, e
	}
	m := map[string]int64{}
	for k, v := range rs {
		i, e := strconv.ParseInt(v, 10, 64)
		if e != nil {
			return nil, e
		}
		m[k] = i
	}
	return m, nil
}

func GetObject[T any](c *XCache, k string) (*T, error) {
	v, e := c.driver.Get(k)
	if e != nil {
		return nil, e
	}
	if v == nil {
		return nil, nil
	}
	b := []byte(*v)
	var t T
	e = msgpack.Unmarshal(b, &t)
	if e != nil {
		return nil, e
	}
	return &t, nil
}

func GetMultipleObject[T any](c *XCache, keys ...string) (map[string]*T, error) {
	rs, e := c.GetMultipleString(keys...)
	if e != nil {
		return nil, e
	}

	m := map[string]*T{}
	var b []byte
	for k, v := range rs {
		b = []byte(v)
		var t T
		e = msgpack.Unmarshal(b, &t)
		if e != nil {
			return nil, e
		}
		m[k] = &t
	}
	return m, nil
}

/* set cache functions */

func (c *XCache) SetBool(k string, v bool, op ...Option) error {
	d, e := c.duration(op...)
	if e != nil {
		return e
	}
	return c.driver.Set(k, strconv.FormatBool(v), d)
}

func (c *XCache) SetString(k string, v string, op ...Option) error {
	d, e := c.duration(op...)
	if e != nil {
		return e
	}
	return c.driver.Set(k, v, d)
}

func (c *XCache) SetMultipleString(m map[string]string, op ...Option) error {
	d, e := c.duration(op...)
	if e != nil {
		return e
	}
	return c.driver.MultiSet(m, d)
}

func (c *XCache) SetInt(k string, v int, op ...Option) error {
	d, e := c.duration(op...)
	if e != nil {
		return e
	}
	b := strconv.Itoa(v)
	return c.driver.Set(k, b, d)
}

func (c *XCache) SetMultipleInt(m map[string]int, op ...Option) error {
	d, e := c.duration(op...)
	if e != nil {
		return e
	}
	mv := map[string]string{}
	for k, v := range m {
		mv[k] = strconv.Itoa(v)
	}
	return c.driver.MultiSet(mv, d)
}

func (c *XCache) SetInt64(k string, v int64, op ...Option) error {
	d, e := c.duration(op...)
	if e != nil {
		return e
	}
	b := strconv.FormatInt(v, 10)
	return c.driver.Set(k, b, d)
}

func (c *XCache) SetMultipleInt64(m map[string]int64, op ...Option) error {
	d, e := c.duration(op...)
	if e != nil {
		return e
	}
	mv := map[string]string{}
	for k, v := range m {
		mv[k] = strconv.FormatInt(v, 10)
	}
	return c.driver.MultiSet(mv, d)
}

func (c *XCache) SetObject(k string, v interface{}, op ...Option) error {
	d, e := c.duration(op...)
	if e != nil {
		return e
	}
	b, e := msgpack.Marshal(v)
	if e != nil {
		return e
	}
	return c.driver.Set(k, string(b), d)
}

func (c *XCache) SetMultipleObject(m map[string]interface{}, op ...Option) error {
	d, e := c.duration(op...)
	if e != nil {
		return e
	}
	mb := map[string]string{}
	for k, v := range m {
		b, e := msgpack.Marshal(v)
		if e != nil {
			return e
		}
		mb[k] = string(b)
	}
	return c.driver.MultiSet(mb, d)
}

func (c *XCache) duration(op ...Option) (time.Duration, error) {
	if len(op) == 0 {
		return c.expiration, nil
	}
	o := options{}
	op[0](&o)
	if o.expiration <= 0 {
		return 0, ErrInvalidDuration
	}
	return o.expiration, nil
}
