package model

import (
	"context"
	"database/sql"
)

const (
	sqlAddBook        = "insert or ignore into book_list (name, md5) values (?, ?)"
	sqlQueryBooks     = "select id, name, md5 from book_list"
	sqlUpdateBookName = "update book_list set name = ? where md5 = ?"
)

type Book struct {
	ID   int64
	Name string
	Md5  string
}

type BookListDatabase struct{}

func (BookListDatabase) AddNewBook(ctx context.Context, book *Book, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, sqlAddBook, book.Name, book.Md5)
	return err
}

func (BookListDatabase) UpdateBookName(ctx context.Context, book *Book, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, sqlUpdateBookName, book.Name, book.Md5)
	return err
}

func (BookListDatabase) QueryBooks(ctx context.Context) ([]*Book, error) {
	var books []*Book
	rows, err := DbConn.QueryContext(ctx, sqlQueryBooks)
	if err != nil {
		return books, err
	}

	defer rows.Close()
	for rows.Next() {
		var id int64
		var name string
		var md5 string
		err := rows.Scan(&id, &name, &md5)
		if err != nil {
			return books, err
		}

		books = append(books, &Book{
			ID:   id,
			Name: name,
			Md5:  md5,
		})
	}

	return books, nil
}
