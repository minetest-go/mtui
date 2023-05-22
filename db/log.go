package db

import (
	"fmt"
	"mtui/types"
	"time"

	"github.com/google/uuid"
	"github.com/minetest-go/dbutil"
)

type LogRepository struct {
	DB dbutil.DBTx
}

func (r *LogRepository) Insert(l *types.Log) error {
	if l.ID == "" {
		l.ID = uuid.NewString()
	}

	if l.Timestamp == 0 {
		l.Timestamp = time.Now().UnixMilli()
	}

	return dbutil.Insert(r.DB, l)
}

func (r *LogRepository) Update(l *types.Log) error {
	return dbutil.Update(r.DB, l, "where id = $1", l.ID)
}

func (r *LogRepository) DeleteBefore(timestamp int64) error {
	_, err := r.DB.Exec("delete from log where timestamp < $1", timestamp)
	return err
}

func (r *LogRepository) buildWhereClause(s *types.LogSearch) (string, []interface{}) {
	q := "where true "
	args := make([]interface{}, 0)
	i := 1

	if s.ID != nil {
		q += fmt.Sprintf(" and id = $%d", i)
		args = append(args, *s.ID)
		i++
	}

	if s.Category != nil {
		q += fmt.Sprintf(" and category = $%d", i)
		args = append(args, *s.Category)
		i++
	}

	if s.Event != nil {
		q += fmt.Sprintf(" and event = $%d", i)
		args = append(args, *s.Event)
		i++
	}

	if s.Username != nil {
		q += fmt.Sprintf(" and username = $%d", i)
		args = append(args, *s.Username)
		i++
	}

	if s.Nodename != nil {
		q += fmt.Sprintf(" and nodename = $%d", i)
		args = append(args, *s.Nodename)
		i++
	}

	if s.PosX != nil {
		q += fmt.Sprintf(" and posx = $%d", i)
		args = append(args, *s.PosX)
		i++
	}

	if s.PosY != nil {
		q += fmt.Sprintf(" and posy = $%d", i)
		args = append(args, *s.PosY)
		i++
	}

	if s.PosZ != nil {
		q += fmt.Sprintf(" and posz = $%d", i)
		args = append(args, *s.PosZ)
		i++
	}

	if s.IPAddress != nil {
		q += fmt.Sprintf(" and ip_address = $%d", i)
		args = append(args, *s.IPAddress)
		i++
	}

	if s.GeoCountry != nil {
		q += fmt.Sprintf(" and geo_country = $%d", i)
		args = append(args, *s.GeoCountry)
		i++
	}

	if s.GeoCity != nil {
		q += fmt.Sprintf(" and geo_city = $%d", i)
		args = append(args, *s.GeoCity)
		i++
	}

	if s.GeoASN != nil {
		q += fmt.Sprintf(" and geo_asn = $%d", i)
		args = append(args, *s.GeoASN)
		i++
	}

	if s.FromTimestamp != nil {
		q += fmt.Sprintf(" and timestamp > $%d", i)
		args = append(args, *s.FromTimestamp)
		i++
	}

	if s.ToTimestamp != nil {
		q += fmt.Sprintf(" and timestamp < $%d", i)
		args = append(args, *s.ToTimestamp)
		i++
	}

	// limit result length to 1000 per default
	limit := 1000
	if s.Limit != nil && *s.Limit < limit {
		limit = *s.Limit
	}
	q += fmt.Sprintf(" order by timestamp desc limit %d", limit)

	return q, args
}

func (r *LogRepository) Search(s *types.LogSearch) ([]*types.Log, error) {
	q, args := r.buildWhereClause(s)
	return dbutil.SelectMulti(r.DB, func() *types.Log { return &types.Log{} }, q, args...)
}

func (r *LogRepository) Count(s *types.LogSearch) (int, error) {
	q, args := r.buildWhereClause(s)
	return dbutil.Count(r.DB, &types.Log{}, q, args...)
}

func (r *LogRepository) GetEvents(c types.LogCategory) ([]string, error) {
	rows, err := r.DB.Query("select event from log where category = $1 group by event", c)
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
