package app

import (
	"net"
	"path"

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
		}).Debug("Skipping geoip resolver setup")
		return &GeoipResolver{}
	}

	asndb, err := geoip2.Open(path.Join(basedir, "GeoLite2-ASN.mmdb"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"reason": err.Error(),
		}).Debug("Skipping geoip resolver setup")
		return &GeoipResolver{}
	}

	return &GeoipResolver{citydb: citydb, asndb: asndb}
}

type GeoipResult struct {
	City       string
	Country    string
	ISOCountry string
	ASN        uint
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

	result.ASN = asn.AutonomousSystemNumber

	return result
}
