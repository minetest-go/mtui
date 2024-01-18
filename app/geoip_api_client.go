package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type ApiGeoIPResolver struct {
	url        string
	cache      map[string]*GeoipResult
	cache_lock *sync.RWMutex
	hc         *http.Client
}

func NewApiGeoIPResolver(url string) GeoIPResolver {
	return &ApiGeoIPResolver{
		url:        url,
		cache:      map[string]*GeoipResult{},
		cache_lock: &sync.RWMutex{},
		hc:         &http.Client{},
	}
}

func (r *ApiGeoIPResolver) Resolve(ipstr string) *GeoipResult {

	r.cache_lock.RLock()
	if r.cache[ipstr] != nil {
		r.cache_lock.RUnlock()
		return r.cache[ipstr]
	}
	r.cache_lock.RUnlock()

	url := fmt.Sprintf("%s/%s", r.url, ipstr)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil
	}

	resp, err := r.hc.Do(req)
	if err != nil {
		return nil
	}
	if resp.StatusCode != 200 {
		return nil
	}

	result := &GeoipResult{}
	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	r.cache_lock.Lock()
	r.cache[ipstr] = result
	r.cache_lock.Unlock()

	return result
}
