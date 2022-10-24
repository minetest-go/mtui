package db

import (
	"database/sql"
	"fmt"
	"mtui/types"
	"strings"
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

func (r *LogRepository) RemoveBefore(timestamp int64) error {
	_, err := r.DB.Exec("delete from log where timestamp < $1", timestamp)
	return err
}

func (r *LogRepository) buildWhereClause(fields string, s *types.LogSearch) (string, []interface{}) {
	l := &types.Log{}
	q := fmt.Sprintf("select %s from %s where true ", fields, l.Table())
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
	l := &types.Log{}
	q, args := r.buildWhereClause(strings.Join(l.Columns(SelectAction), ","), s)
	rows, err := r.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	list := make([]*types.Log, 0)
	for rows.Next() {
		entry := &types.Log{}
		err := entry.Scan(SelectAction, rows.Scan)
		if err != nil {
			return nil, err
		}
		list = append(list, entry)
	}
	return list, nil
}

func (r *LogRepository) Count(s *types.LogSearch) (int64, error) {
	q, args := r.buildWhereClause("count(*)", s)
	row := r.DB.QueryRow(q, args...)
	var count int64
	err := row.Scan(&count)
	return count, err
}
