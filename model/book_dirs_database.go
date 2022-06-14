package model

import "context"

const (
	// bookDirTable = "book_dir"

	sqlAddBookDir    = "insert or ignore into book_dir (id, path) values (1, ?)"
	sqlUpdateBookDir = "update book_dir set path = ? where id = 1"
	sqlQueryPath     = "select path from book_dir where id = 1"
)

type BookDirDatabase struct{}

func (*BookDirDatabase) AddOrUpdateBookDirPath(ctx context.Context, path string) (bool, error) {
	res, err := DbConn.ExecContext(ctx, sqlUpdateBookDir, path)
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

	res, err = DbConn.ExecContext(ctx, sqlAddBookDir, path)
	if err != nil {
		return false, err
	}

	rows, err = res.RowsAffected()
	return rows > 0, err
}

func (*BookDirDatabase) QueryPath(ctx context.Context) (string, error) {
	rows, err := DbConn.QueryContext(ctx, sqlQueryPath)
	if err != nil {
		return "", err
	}

	var path string
	rows.Next()
	err = rows.Scan(&path)
	_ = rows.Close()

	return path, err
}
