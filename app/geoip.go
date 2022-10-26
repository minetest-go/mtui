package app

import (
	"fmt"
	"net"
	"path"

	"github.com/oschwald/geoip2-golang"
	"github.com/sirupsen/logrus"
)

type GeoipResolver struct {
	db *geoip2.Reader
}

func NewGeoipResolver(basedir string) *GeoipResolver {
	db, err := geoip2.Open(path.Join(basedir, "GeoLite2-City.mmdb"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"reason": err.Error(),
		}).Debug("Skipping geoip resolver setup")
		return &GeoipResolver{}
	}

	return &GeoipResolver{db: db}
}

type GeoipResult struct {
	City    string
	Country string
	ASN     string
}

func (r *GeoipResolver) Resolve(ipstr string) *GeoipResult {
	if r.db == nil {
		return nil
	}

	ip := net.ParseIP(ipstr)
	result := &GeoipResult{}

	city, err := r.db.City(ip)
	if err != nil {
		return nil
	}

	result.City = city.City.Names["en"]
	result.Country = city.Country.Names["en"]

	asn, err := r.db.ASN(ip)
	if err != nil {
		return nil
	}

	result.ASN = fmt.Sprintf("%s (id: %d)", asn.AutonomousSystemOrganization, asn.AutonomousSystemNumber)

	return result
}
