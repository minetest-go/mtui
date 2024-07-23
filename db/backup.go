package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func BackupSqlite3Database(ctx context.Context, srcfile, dstfile string) error {
	srcdb, err := sql.Open("sqlite3", "file:"+srcfile)
	if err != nil {
		return fmt.Errorf("source db open error: %v", err)
	}

	dstdb, err := sql.Open("sqlite3", "file:"+dstfile)
	if err != nil {
		return fmt.Errorf("destination db open error: %v", err)
	}

	srcconn, err := srcdb.Conn(context.Background())
	if err != nil {
		return fmt.Errorf("source conn open error: %v", err)
	}
	defer srcconn.Close()

	dstconn, err := dstdb.Conn(ctx)
	if err != nil {
		return fmt.Errorf("destination conn open error: %v", err)
	}
	defer dstconn.Close()

	var srcsq3, dstsq3 *sqlite3.SQLiteConn

	err = srcconn.Raw(func(driverConn any) error {
		var ok bool
		srcsq3, ok = driverConn.(*sqlite3.SQLiteConn)
		if !ok {
			return errors.New("failed to get sqlite3 connection")
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("source raw open error: %v", err)
	}

	err = dstconn.Raw(func(driverConn any) error {
		var ok bool
		dstsq3, ok = driverConn.(*sqlite3.SQLiteConn)
		if !ok {
			return errors.New("failed to get sqlite3 connection")
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("destination raw open error: %v", err)
	}

	b, err := dstsq3.Backup("main", srcsq3, "main")
	if err != nil {
		return fmt.Errorf("backup error: %v", err)
	}
	defer b.Finish()

	done := false
	for !done {
		done, err = b.Step(-1)
		if err != nil {
			return fmt.Errorf("step error: %v", err)
		}
	}

	return nil
}
