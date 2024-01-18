package app

import (
	"mtui/types"
	"net/http"
	"strings"
)

type GeoipResult struct {
	City       string `json:"city"`
	Country    string `json:"country"`
	ISOCountry string `json:"country_iso"`
	ASN        int    `json:"asn"`
}

type GeoIPResolver interface {
	Resolve(ipstr string) *GeoipResult
}

type NoOPGeoIPResolver struct{}

func (r *NoOPGeoIPResolver) Resolve(ipstr string) *GeoipResult {
	return nil
}

func NewGeoIPResolver(basedir, api_url string) GeoIPResolver {
	// api-client -> mmdb -> no-op
	if api_url != "" {
		return NewApiGeoIPResolver(api_url)
	}

	r := NewGeoIPMMDBResolver(basedir)
	if r == nil {
		return &NoOPGeoIPResolver{}
	}
	return r
}

func (r *App) ResolveLogGeoIP(l *types.Log, req *http.Request) {
	if req != nil {
		// web request
		fwdfor := req.Header.Get("X-Forwarded-For")
		if fwdfor != "" {
			// behind reverse proxy
			parts := strings.Split(fwdfor, ",")
			l.IPAddress = &parts[0]
		} else {
			// direct access
			parts := strings.Split(req.RemoteAddr, ":")
			l.IPAddress = &parts[0]
		}
	}

	if l.IPAddress != nil {
		geoip := r.GeoipResolver.Resolve(*l.IPAddress)
		if geoip != nil {
			l.GeoCity = &geoip.City
			l.GeoCountry = &geoip.ISOCountry
			l.GeoASN = &geoip.ASN
		}
	}
}
