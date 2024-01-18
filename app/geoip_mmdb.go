package app

import (
	"net"
	"os"
	"path"

	"github.com/oschwald/geoip2-golang"
	"github.com/sirupsen/logrus"
)

type GeoIPMMDBResolver struct {
	citydb *geoip2.Reader
	asndb  *geoip2.Reader
}

const CITY_MMDB_NAME = "GeoLite2-City.mmdb"
const ASN_MMDB_NAME = "GeoLite2-ASN.mmdb"

func NewGeoIPMMDBResolver(basedir string) *GeoIPMMDBResolver {
	resolver := &GeoIPMMDBResolver{}
	var err error

	citydb_name := path.Join(basedir, CITY_MMDB_NAME)
	fs, _ := os.Stat(citydb_name)
	if fs == nil {
		return nil
	}

	resolver.citydb, err = geoip2.Open(citydb_name)
	if err != nil {
		panic(err)
	}

	asndb_name := path.Join(basedir, ASN_MMDB_NAME)
	fs, _ = os.Stat(asndb_name)
	if fs == nil {
		return nil
	}

	resolver.asndb, err = geoip2.Open(asndb_name)
	if err != nil {
		panic(err)
	}

	logrus.Info("geoip database setup successfully")
	return resolver
}

func (r *GeoIPMMDBResolver) Resolve(ipstr string) *GeoipResult {
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
