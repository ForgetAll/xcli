package cmd

import (
	"bufio"
	"context"
	"crypto/md5"
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/ForgetAll/xcli/model"
	hook "github.com/robotn/gohook"
	"github.com/spf13/cobra"
)

var (
	path string

	code       chan int
	saveSingle chan bool
	CodeExit   = 1
	CodeNext   = 2
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read",
	Short: "read book",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var outputInfo string
		var bddb model.BookDirDatabase
		ctx := context.TODO()

		if len(path) > 0 {
			result, err := bddb.AddOrUpdateBookDirPath(ctx, path)
			if err != nil {
				fmt.Printf("set path error: %v", err)
				return
			}

			if !result {
				fmt.Println("set path failed")
				return
			}

			return
		}

		p, err := bddb.QueryPath(ctx)
		if err != nil {
			fmt.Printf("query path error: %v", err)
			return
		}

		if !strings.HasSuffix(p, "/") {
			p += "/"
		}

		outputInfo += "book dir path: " + p + "\n"
		path = p

		err = updateBooks(ctx, path)
		if err != nil {
			fmt.Printf("read dir error: %v, error: $%v\n", path, err)
			return
		}

		fmt.Println(outputInfo)
		var bldb model.BookListDatabase
		books, err := bldb.QueryBooks(ctx)
		if err != nil {
			fmt.Printf("query books error: %v", err)
		}

		fmt.Println("input number to read: ")
		bookMap := make(map[int64]*model.Book)
		for _, b := range books {
			bookMap[b.ID] = b
			fmt.Printf("%v  《%v》\n", b.ID, strings.ReplaceAll(b.Name, path, ""))
		}

		var bookPath string
		var currentBook *model.Book
		for {
			var num int64
			n, err := fmt.Scan(&num)
			if err != nil {
				fmt.Printf("%v", err)
				return
			}

			if n == 0 {
				return
			}

			if b, ok := bookMap[num]; !ok {
				fmt.Printf("number %v book not exist, please re-enter: \n", num)
			} else {
				bookPath = path + b.Name
				currentBook = b
				break
			}
		}

		f, err := os.Open(bookPath)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		var rldb model.ReadLogDatabase
		readLog, err := rldb.QueryReadLog(ctx, currentBook.ID)
		if err != nil {
			fmt.Printf("%v", err)
			return
		}

		lineCount := readLog.LineCount
		if lineCount == 0 {
			lineCount = 1
		}

		code = make(chan int)
		saveSingle = make(chan bool)
		go keyEventListen()
		var next bool
		exit := false
		fileScanner := bufio.NewScanner(f)
		var count int64 = 1
		go autoSaveReadLog(&lineCount, currentBook.ID)
		for fileScanner.Scan() {
			if lineCount >= count {
				count++
				continue
			} else {
				next = true
			}

			if exit {
				break
			}

			if next {
				fmt.Println(fileScanner.Text())
				saveSingle <- true
				// nolint
				next = false
			}

			input := <-code
			switch input {
			case CodeExit:
				exit = true
			case CodeNext:
				// nolint
				next = true
				lineCount++
				count++
			default:
				continue
			}
		}
		model.Release()
	},
}

func autoSaveReadLog(lineCount *int64, bookID int64) {
	var rldb model.ReadLogDatabase
	ctx := context.TODO()
	for {
		<-saveSingle
		_, _ = rldb.AddOrUpdateReadLog(ctx, *lineCount, bookID)
		time.Sleep(1000)
	}
}

func keyEventListen() {
	hook.Register(hook.KeyDown, []string{"q"}, func(e hook.Event) {
		code <- CodeExit
		hook.End()
	})

	hook.Register(hook.KeyDown, []string{"j"}, func(e hook.Event) {
		code <- CodeNext
	})

	s := hook.Start()
	<-hook.Process(s)
}

func init() {
	rootCmd.AddCommand(readCmd)
	model.InitDb()
	readCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "books dir")
}

func updateBooks(ctx context.Context, path string) error {
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	txtFile := make([]fs.FileInfo, 0)
	// var dirBooks []*model.Book
	dirBooks := make([]*model.Book, 0)
	for _, f := range dir {
		if !strings.HasSuffix(f.Name(), "txt") {
			continue
		}

		txtFile = append(txtFile, f)

		file, err := os.Open(path + f.Name())
		if err != nil {
			return err
		}

		md5Hash := md5.New()
		_, err = io.Copy(md5Hash, file)
		if err != nil {
			return err
		}

		md5Str := fmt.Sprintf("%x", md5Hash.Sum(nil))
		dirBooks = append(dirBooks, &model.Book{
			Name: f.Name(),
			Md5:  md5Str,
		})
		_ = file.Close()
	}

	if len(txtFile) == 0 {
		return fmt.Errorf("no books in %v", path)
	}

	var bldb model.BookListDatabase
	books, err := bldb.QueryBooks(ctx)
	if err != nil {
		return err
	}

	var needUpdateBooks []*model.Book
	var needAddBooks []*model.Book
	dbBooksMap := make(map[string]*model.Book)
	for _, b := range books {
		dbBooksMap[b.Md5] = b
	}

	for _, b := range dirBooks {
		book, ok := dbBooksMap[b.Md5]
		if !ok {
			needAddBooks = append(needAddBooks, b)
			continue
		}

		if book.Name != b.Name {
			needUpdateBooks = append(needUpdateBooks, b)
		}
	}

	fmt.Printf("%v book need to be added, %v book need to be update\n", len(needAddBooks), len(needUpdateBooks))
	err = addNewBooks(ctx, needAddBooks)
	if err != nil {
		return err
	}
	fmt.Printf("%v book add success\n", len(needAddBooks))

	err = updateBookInfo(ctx, needUpdateBooks)
	if err != nil {
		return err
	}
	fmt.Printf("%v book update success\n", len(needUpdateBooks))
	return nil
}

func addNewBooks(ctx context.Context, books []*model.Book) error {
	var bldb model.BookListDatabase
	tx, err := model.DbConn.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: false,
	})
	if err != nil {
		return err
	}

	for _, b := range books {
		err := bldb.AddNewBook(ctx, b, tx)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	fmt.Println("commit transaction")
	return tx.Commit()
}

func updateBookInfo(ctx context.Context, books []*model.Book) error {
	var bldb model.BookListDatabase
	tx, err := model.DbConn.BeginTx(ctx, &sql.TxOptions{
		ReadOnly: false,
	})
	if err != nil {
		return err
	}

	for _, b := range books {
		err := bldb.UpdateBookName(ctx, b, tx)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
