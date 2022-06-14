package model

import (
	"context"
	"database/sql"
)

const (
	sqlAddReadLog    = "insert or ignore into read_log (book_id, line_count) values(?, ?)"
	sqlUpdateReadLog = "update read_log set line_count = ? where book_id = ?"
	sqlQueryReadLog  = "select line_count from read_log where book_id = ?"
)

type ReadLogDatabase struct{}

type ReadLog struct {
	LineCount int64
}

func (ReadLogDatabase) QueryReadLog(ctx context.Context, bookID int64) (*ReadLog, error) {
	var readLog ReadLog
	rows, err := DbConn.QueryContext(ctx, sqlQueryReadLog, bookID)
	if err == sql.ErrNoRows {
		return &readLog, nil
	}
	if err != nil {
		return &readLog, err
	}

	defer rows.Close()
	for rows.Next() {
		var lineCount int64
		err := rows.Scan(&lineCount)
		if err != nil {
			return &readLog, nil
		}

		readLog.LineCount = lineCount
	}

	return &readLog, nil
}

func (ReadLogDatabase) AddOrUpdateReadLog(ctx context.Context, lineCount, bookID int64) (bool, error) {
	res, err := DbConn.ExecContext(ctx, sqlUpdateReadLog, lineCount, bookID)
	if err != nil {
		return false, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	if rows > 0 {
		return true, nil
	}

	res, err = DbConn.ExecContext(ctx, sqlAddReadLog, bookID, lineCount)
	if err != nil {
		return false, err
	}

	rows, err = res.RowsAffected()
	return rows > 0, err
}
