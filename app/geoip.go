package app

import (
	"mtui/types"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/oschwald/geoip2-golang"
	"github.com/sirupsen/logrus"
)

type GeoipResolver struct {
	citydb *geoip2.Reader
	asndb  *geoip2.Reader
}

func NewGeoipResolver(basedir string) *GeoipResolver {
	citydb, err := geoip2.Open(path.Join(basedir, "GeoLite2-City.mmdb"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"reason": err.Error(),
		}).Info("Skipping geoip resolver setup (missing city db)")
		return &GeoipResolver{}
	}

	asndb, err := geoip2.Open(path.Join(basedir, "GeoLite2-ASN.mmdb"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"reason": err.Error(),
		}).Info("Skipping geoip resolver setup (missing asn db)")
		return &GeoipResolver{}
	}

	return &GeoipResolver{citydb: citydb, asndb: asndb}
}

type GeoipResult struct {
	City       string `json:"city"`
	Country    string `json:"country"`
	ISOCountry string `json:"country_iso"`
	ASN        int    `json:"asn"`
}

func (r *GeoipResolver) Resolve(ipstr string) *GeoipResult {
	if r.citydb == nil || r.asndb == nil {
		return nil
	}

	ip := net.ParseIP(ipstr)
	result := &GeoipResult{}

	city, err := r.citydb.City(ip)
	if err != nil {
		return nil
	}

	result.City = city.City.Names["en"]
	result.Country = city.Country.Names["en"]
	result.ISOCountry = city.Country.IsoCode

	asn, err := r.asndb.ASN(ip)
	if err != nil {
		return nil
	}

	result.ASN = int(asn.AutonomousSystemNumber)

	return result
}

func (a *App) ResolveLogGeoIP(l *types.Log, r *http.Request) {
	if r == nil {
		// nothing to work with
		return
	}

	fwdfor := r.Header.Get("X-Forwarded-For")
	if fwdfor != "" {
		// behind reverse proxy
		parts := strings.Split(fwdfor, ",")
		l.IPAddress = &parts[0]
	} else {
		// direct access
		parts := strings.Split(r.RemoteAddr, ":")
		l.IPAddress = &parts[0]
	}

	if a.GeoipResolver != nil && l.IPAddress != nil {
		geoip := a.GeoipResolver.Resolve(*l.IPAddress)
		if geoip != nil {
			l.GeoCity = &geoip.City
			l.GeoCountry = &geoip.ISOCountry
			l.GeoASN = &geoip.ASN
		}
	}
}
