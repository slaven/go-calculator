package calcserver

import (
	"fmt"
	"time"

	"github.com/patrickmn/go-cache"
)

type CalcCache struct {
	items *cache.Cache
}

type Caching interface {
	Get(string) (float64, bool)
	SetOrUpdate(string, float64) (bool, error)
}

// Create cache instance
func NewCache() *CalcCache {
	return &CalcCache{
		items: cache.New(getCacheExpiration(), getCacheExpiration()),
	}
}

// Get from cache
func (cc *CalcCache) Get(key string) (float64, bool) {
	val, found := cc.items.Get(key)
	if !found {
		return 0, false
	}

	value, ok := val.(float64)
	if !ok {
		cc.items.Delete(key)
		return 0, false
	}
	return value, found
}

// Set or update cache
func (cc *CalcCache) SetOrUpdate(key string, val float64) (bool, error) {
	_, found := cc.items.Get(key)
	if found {
		replaceErr := cc.items.Replace(key, val, getCacheExpiration())
		if replaceErr != nil {
			return false, replaceErr
		}
		return true, nil
	}

	addErr := cc.items.Add(key, val, getCacheExpiration())
	if addErr != nil {
		return false, addErr
	}
	return true, nil
}

// Build cache key from params
func buildCacheKey(calculation CalcOperation, x, y float64) string {
	// Order of params is irrelevant for add and multiply calculations
	if calculation == addCalc || calculation == multiplyCalc {
		if x > y {
			return fmt.Sprintf("%s|%v|%v", calculation, y, x)
		}
		return fmt.Sprintf("%s|%v|%v", calculation, x, y)
	}
	return fmt.Sprintf("%s|%v|%v", calculation, x, y)
}

// Get cache expiration duration
func getCacheExpiration() time.Duration {
	// TODO: get from config
	return time.Minute
}
