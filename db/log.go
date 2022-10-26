package db

import (
	"database/sql"
	"fmt"
	"mtui/types"
	"time"

	"github.com/google/uuid"
)

type LogRepository struct {
	DB *sql.DB
}

func (r *LogRepository) Insert(l *types.Log) error {
	if l.ID == "" {
		l.ID = uuid.NewString()
	}

	if l.Timestamp == 0 {
		l.Timestamp = time.Now().UnixMilli()
	}

	return Insert(r.DB, l)
}

func (r *LogRepository) Update(l *types.Log) error {
	return Update(r.DB, l, map[string]any{"id": l.ID})
}

func (r *LogRepository) RemoveBefore(timestamp int64) error {
	_, err := r.DB.Exec("delete from log where timestamp < $1", timestamp)
	return err
}

func (r *LogRepository) buildWhereClause(s *types.LogSearch) (string, []interface{}) {
	q := " where true "
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
	q += fmt.Sprintf(" limit %d", limit)

	return q, args
}

func (r *LogRepository) Search(s *types.LogSearch) ([]*types.Log, error) {
	q, args := r.buildWhereClause(s)
	return SelectMulti(r.DB, func() *types.Log { return &types.Log{} }, q, args...)
}

func (r *LogRepository) Count(s *types.LogSearch) (int64, error) {
	q, args := r.buildWhereClause(s)

	row := r.DB.QueryRow("select count(*) from log"+q, args...)
	var count int64
	err := row.Scan(&count)
	return count, err
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
