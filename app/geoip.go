package app

import (
	"mtui/types"
	"net"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/oschwald/geoip2-golang"
	"github.com/sirupsen/logrus"
)

type GeoipResolver struct {
	citydb *geoip2.Reader
	asndb  *geoip2.Reader
}

const CITY_MMDB_NAME = "GeoLite2-City.mmdb"
const ASN_MMDB_NAME = "GeoLite2-ASN.mmdb"

func NewGeoipResolver(basedir string) *GeoipResolver {
	resolver := &GeoipResolver{}
	var err error

	citydb_name := path.Join(basedir, CITY_MMDB_NAME)
	fs, _ := os.Stat(citydb_name)
	if fs == nil {
		return &GeoipResolver{}
	}

	resolver.citydb, err = geoip2.Open(citydb_name)
	if err != nil {
		panic(err)
	}

	asndb_name := path.Join(basedir, ASN_MMDB_NAME)
	fs, _ = os.Stat(asndb_name)
	if fs == nil {
		return &GeoipResolver{}
	}

	resolver.asndb, err = geoip2.Open(asndb_name)
	if err != nil {
		panic(err)
	}

	logrus.Info("geoip database setup successfully")
	return resolver
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
