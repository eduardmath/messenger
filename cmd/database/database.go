package database

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"time"
)

const databaseUrl string = "postgres://postgres:postgres@localhost:55000"

type Database struct {
	s    string
	pool *pgxpool.Pool
}

func help() {
	// md5
	fmt.Printf("%x\n", md5.Sum([]byte("hello")))
	fmt.Printf("%x\n", md5.Sum([]byte("hella")))
	// return

	// time
	t := time.Now()
	fmt.Println(t)
	year, month, day := t.Date()
	fmt.Println(year, month, day)
	// return

}

func (db Database) Print(s string) string {
	var ss string
	var len = len(s)
	for i := 0; i < len; i++ {
		ss += string(s[len-i-1])
	}

	return ss
}

func (db Database) CheckUser(name string) bool {
	return true
}

func (db Database) Connect() bool {
	var errDB error
	db.pool, errDB = pgxpool.New(context.Background(), databaseUrl)

	if errDB != nil {
		return false
	}
	defer db.pool.Close()

	return true
}

func ExecuteInsert(pool *pgxpool.Pool) {
	var f, l, d string
	fmt.Scan(&f, &l, &d)
	var id int
	err := pool.QueryRow(context.Background(), "INSERT INTO users (first_name, last_name, date) VALUES ($1, $2, $3) RETURNING id", f, l, d).Scan(&id)
	if err != nil {
		panic(err)
	}
	fmt.Println("New record ID is:", id)

}

func ExecuteSelectQuery(dbPool *pgxpool.Pool) {
	log.Println("starting execution of select query")

	// execute the query and get result rows
	rows, err := dbPool.Query(context.Background(), "select * from users")
	if err != nil {
		log.Fatal("error while executing query")
	}

	log.Println("result:")

	// iterate through the rows
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal("error while iterating dataset")
		}

		// convert DB types to Go types
		id := values[0].(int32)
		firstName := values[1].(string)
		lastName := values[2].(string)
		dateOfBirth := values[3].(string) // (string) // values[3].(time.Time)
		log.Println("[id:", id, ", first_name:", firstName, ", last_name:", lastName, ", date:", dateOfBirth, "]")
	}
}

func ExecuteFunction(dbPool *pgxpool.Pool, id int) {
	log.Println("starting execution of database function")

	// execute the query and get result rows
	rows, err := dbPool.Query(context.Background(), "select * from get_ttt($1)", id)
	log.Println("input id: ", id)
	if err != nil {
		log.Fatal("error while executing query")
	}

	log.Println("result:")

	// iterate through the rows
	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			log.Fatal("error while iterating dataset")
		}

		//convert DB types to Go types
		firstName := values[0].(string)
		lastName := values[1].(string)
		dateOfBirth := values[2].(time.Time)

		log.Println("[first_name:", firstName, ", last_name:", lastName, ", date:", dateOfBirth, "]")
	}

}
