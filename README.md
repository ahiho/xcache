# XCache - Just a simple cache wrapper

Inspired from `IDistributedCache` of `Asp.net`

Caching is very general things backend usage. In development or small scale we can just use in-memory cache libraries such as `go-cache`, `ristretto`. But when application scale up and need run in multiple replicas, it's neccessary to using centralize cache such as `Redis` or `Memcached`. But ofc no one wanna refactor all of cache usage code with new libraries. 

We made this as a simple interfaces function for `xcache` and can use `redis` or `ristretto` as a backed end for store data. When change from `ristretto` to `redis` we only need change code for create driver. 

You can also develop your own driver which just need implement below functions 
```
type Driver interface {
	Set(k string, v string, d time.Duration) error
	Get(k string) (*string, error)
	Del(k string) error
	MultiSet(m map[string]string, d time.Duration) error
	MultiGet(k []string) (map[string]string, error)
	MultiDel(k []string) error
}
```