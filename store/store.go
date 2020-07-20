package store

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
	"time"
)

const (
	host     = "raja.db.elephantsql.com"
	port     = 5432
	user     = "stzrurfj"
	password = "mHLqxoPKfj2P0R5XD2AImPSr8Ozu7rWr"
	dbname   = "stzrurfj"
)

var db *sqlx.DB

func init() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	dbx := sqlx.NewDb(conn, "postgres")
	dbx.SetMaxIdleConns(2)
	dbx.SetConnMaxLifetime(10 * time.Minute)
	dbx.SetMaxOpenConns(5)
	log.Println("database connection created")
	db = dbx
}

type User struct {
	Id       uint64 `json:"id" db:"id"`
	Uid      string `json:"uid" db:"uid"`
	Name     string `json:"name" db:"name"`
	Age      int    `json:"age" db:"age"`
	Address  string `json:"address" db:"address"`
	MobileNo string `json:"mobile_no" db:"mobile"`
}

func GetUserCountByName(name string, mobNo string) (count uint64, err error) {
	if db == nil {
		err = errors.New("unable to connect to database")
		return
	}

	q := `SELECT count(id) FROM user_details WHERE name=$1 and mobile=$2;`
	err = db.Get(&count, q, name, mobNo)

	return
}

func AddNewUser(u User) (resp User, err error) {
	if db == nil {
		err = errors.New("unable to connect to database")
		return
	}
	q := `INSERT INTO user_details (uid, name, age, address, mobile)
	VALUES ($1, $2, $3, $4, $5) returning id, uid, name, age, address, mobile`

	err = db.Get(&resp, q, u.Uid, u.Name, u.Age, u.Address, u.MobileNo)

	return
}

func GetAllUsers() (r []User, err error) {
	if db == nil {
		err = errors.New("unable to connect to database")
		return
	}
	q := `SELECT id, uid, name, age, address, mobile FROM user_details;`

	err = db.Select(&r, q)

	return
}

func GetAllUsersCount() (r int64, err error) {
	if db == nil {
		err = errors.New("unable to connect to database")
		return
	}
	q := `SELECT count(id) FROM user_details;`

	err = db.Select(&r, q)

	return
}

func GetUsersWithPagination(limit int, offset int) (r []User, err error) {
	if db == nil {
		err = errors.New("unable to connect to database")
		return
	}
	q := `SELECT id, uid, name, age, address, mobile FROM user_details LIMIT $1 OFFSET $2;`

	err = db.Select(&r, q, limit, offset)

	return
}
