package cdb

import (
	"fmt"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
)

type CachedCDBClient struct {
	ttl              time.Duration
	client           *CDBClient
	search_cache     *cache.Cache[string, []*Package]
	dependency_cache *cache.Cache[string, PackageDependency]
	detail_cache     *cache.Cache[string, *PackageDetails]
	updates          PackageUpdates
	updates_time     time.Time
}

func NewCachedClient(client *CDBClient, ttl time.Duration) *CachedCDBClient {
	return &CachedCDBClient{
		ttl:              ttl,
		client:           client,
		search_cache:     cache.New[string, []*Package](),
		dependency_cache: cache.New[string, PackageDependency](),
		detail_cache:     cache.New[string, *PackageDetails](),
	}
}

func (c *CachedCDBClient) SearchPackages(q *PackageQuery) ([]*Package, error) {
	key := q.Params().Encode()
	res, ok := c.search_cache.Get(key)

	var err error
	if !ok {
		res, err = c.client.SearchPackages(q)
		if err != nil {
			return nil, err
		}
		c.search_cache.Set(key, res, cache.WithExpiration(c.ttl))
	}
	return res, nil
}

func (c *CachedCDBClient) GetDependencies(author, name string) (PackageDependency, error) {
	key := fmt.Sprintf("%s/%s", author, name)
	res, ok := c.dependency_cache.Get(key)

	var err error
	if !ok {
		res, err = c.client.GetDependencies(author, name)
		if err != nil {
			return nil, err
		}
		c.dependency_cache.Set(key, res, cache.WithExpiration(c.ttl))
	}
	return res, nil
}

func (c *CachedCDBClient) GetDetails(author, name string) (*PackageDetails, error) {
	key := fmt.Sprintf("%s/%s", author, name)
	res, ok := c.detail_cache.Get(key)

	var err error
	if !ok {
		res, err = c.client.GetDetails(author, name)
		if err != nil {
			return nil, err
		}
		c.detail_cache.Set(key, res, cache.WithExpiration(c.ttl))
	}
	return res, nil
}

func (c *CachedCDBClient) GetUpdates() (PackageUpdates, error) {
	last_update := time.Since(c.updates_time)
	if last_update > c.ttl || c.updates == nil {
		var err error
		c.updates, err = c.client.GetUpdates()
		if err != nil {
			return nil, err
		}
		c.updates_time = time.Now()
	}
	return c.updates, nil
}
