package db

import (
	"fmt"
	"mtui/types"
	"time"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type LogRepository struct {
	dbu *dbutil.DBUtil[*types.Log]
	db  dbutil.DBTx
}

func (r *LogRepository) Insert(l *types.Log) error {
	if l.ID == "" {
		l.ID = uuid.NewString()
	}

	if l.Timestamp == 0 {
		l.Timestamp = time.Now().UnixMilli()
	}

	return r.dbu.Insert(l)
}

func (r *LogRepository) Update(l *types.Log) error {
	return r.dbu.Update(l, "where id = %s", l.ID)
}

func (r *LogRepository) DeleteBefore(timestamp int64) error {
	_, err := r.db.Exec("delete from log where timestamp < ?1", timestamp)
	return err
}

func (r *LogRepository) buildWhereClause(s *types.LogSearch) (string, []interface{}) {
	q := "where true "
	args := make([]interface{}, 0)

	if s.ID != nil {
		q += " and id = %s"
		args = append(args, *s.ID)
	}

	if s.Category != nil {
		q += " and category = %s"
		args = append(args, *s.Category)
	}

	if s.Event != nil {
		q += " and event = %s"
		args = append(args, *s.Event)
	}

	if s.Username != nil {
		q += " and username = %s"
		args = append(args, *s.Username)
	}

	if s.Nodename != nil {
		q += " and nodename = %s"
		args = append(args, *s.Nodename)
	}

	if s.PosX != nil {
		q += " and posx = %s"
		args = append(args, *s.PosX)
	}

	if s.PosY != nil {
		q += " and posy = %s"
		args = append(args, *s.PosY)
	}

	if s.PosZ != nil {
		q += " and posz = %s"
		args = append(args, *s.PosZ)
	}

	if s.IPAddress != nil {
		q += " and ip_address = %s"
		args = append(args, *s.IPAddress)
	}

	if s.GeoCountry != nil {
		q += " and geo_country = %s"
		args = append(args, *s.GeoCountry)
	}

	if s.GeoCity != nil {
		q += " and geo_city = %s"
		args = append(args, *s.GeoCity)
	}

	if s.GeoASN != nil {
		q += " and geo_asn = %s"
		args = append(args, *s.GeoASN)
	}

	if s.FromTimestamp != nil {
		q += " and timestamp > %s"
		args = append(args, *s.FromTimestamp)
	}

	if s.ToTimestamp != nil {
		q += " and timestamp < %s"
		args = append(args, *s.ToTimestamp)
	}

	// limit result length to 1000 per default
	limit := 1000
	if s.Limit != nil && *s.Limit < limit {
		limit = *s.Limit
	}

	sort := "desc"
	if s.SortAscending {
		sort = "asc"
	}
	q += fmt.Sprintf(" order by timestamp %s limit %d", sort, limit)

	return q, args
}

func (r *LogRepository) Search(s *types.LogSearch) ([]*types.Log, error) {
	q, args := r.buildWhereClause(s)
	return r.dbu.SelectMulti(q, args...)
}

func (r *LogRepository) Count(s *types.LogSearch) (int, error) {
	q, args := r.buildWhereClause(s)
	return r.dbu.Count(q, args...)
}

func (r *LogRepository) GetEvents(c types.LogCategory) ([]string, error) {
	rows, err := r.db.Query("select event from log where category = ?1 group by event", c)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]string, 0)
	for rows.Next() {
		var e string
		err = rows.Scan(&e)
		if err != nil {
			return nil, err
		}

		result = append(result, e)
	}

	return result, nil
}
