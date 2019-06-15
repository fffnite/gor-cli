package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var (
	createTable = `CREATE TABLE IF NOT EXISTS creds (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		client_id TEXT,
		client_secret_hash BLOB,
		tag TEXT
		)`
	flagNew    = flag.Bool("new", false, "Create new creds")
	flagTag    = flag.String("tag", "notSet", "human readable tag for new cred")
	flagList   = flag.Bool("list", false, "List existing credentials")
	flagRemove = flag.String("remove", "", "Creds to remove based off ID")
)

type Result struct {
	Client_Id string
	Tag       string
}

func main() {

	db := connectDatabase()
	defer db.Close()
	db.Exec(createTable)
	flag.Parse()
	if *flagNew {
		create(db, *flagTag)
	}
	if *flagList {
		list(db)
	}
	if *flagRemove != "" {
		remove(db, *flagRemove)
	}
}

func create(db *sql.DB, tag string) {
	u, err := randomHex(25)
	if err != nil {
		fmt.Println(err)
		return
	}
	p, err := randomHex(32)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("IMPORTANT: store somewhere safe as this is not recoverable")
	fmt.Printf("client_id: %v \n", u)
	fmt.Printf("client_secret: %v \n", p)
	en := encrypt(p)
	_, err = squirrel.
		Insert("creds").
		Columns("client_id", "client_secret_hash", "tag").
		Values(u, en, tag).
		RunWith(db).
		Exec()
	if err != nil {
		fmt.Println(err)
		return
	}

	return
}

func encrypt(s string) []byte {
	b := []byte(s)
	en, err := bcrypt.GenerateFromPassword(b, 10)
	if err != nil {
		fmt.Println(err)
		n := make([]byte, 0)
		return n
	}
	return en
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func list(db *sql.DB) {
	rows, err := squirrel.
		Select("client_id", "tag").
		From("creds").
		OrderBy("id").
		RunWith(db).
		Query()
	if err != nil {
		fmt.Println(err)
		return
	}
	var results []Result
	for rows.Next() {
		var r Result
		err := rows.Scan(&r.Client_Id, &r.Tag)
		if err != nil {
			fmt.Println(err)
		}
		results = append(results, r)
	}
	fmt.Printf("list: %v", results)
}

func connectDatabase() *sql.DB {
	conn, err := sql.Open("sqlite3", "credstore")
	if err != nil {
		fmt.Println(err)
	}
	return conn
}

func remove(db *sql.DB, s string) {
	_, err := squirrel.
		Delete("").
		From("creds").
		Where("client_id = ?", s).
		RunWith(db).
		Exec()
	if err != nil {
		fmt.Println(err)
	}
}
