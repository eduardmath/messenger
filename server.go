package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"net"
	"os"
	"strings"
)

func process(pull *map[string]net.Conn, c net.Conn) {
	var friend string
	buf := make([]byte, 1024*8)

	// получаем доступ к текущему соединению
	conn := c

	rLen, _ := conn.Read(buf)

	var name = string(buf[:rLen])

	(*pull)[name] = conn

	// определим, что перед выходом из функции, мы закроем соединение
	fmt.Println("Accept user:", name)
	defer conn.Close()

	for {
		readLen, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Close ", name)
				delete(*pull, name)
				break
			}
			fmt.Println(err)
		}

		// Распечатываем полученое сообщение
		var message = string(buf[:readLen])

		// handler commands
		switch message {
		case "/h":
			conn.Write([]byte(fmt.Sprintf("%s\n", "help messenger")))
			continue
		case "/chat", "/c":
			pr(pull)
			// fmt.Println((*pull))
			conn.Write([]byte(fmt.Sprintf("%s", pr(pull))))
			continue
		}

		// парсинг полученного сообщения
		_, err = fmt.Sscanf(message, "%s", &friend) // определи номер клиента

		if err != nil {
			// обработка ошибки формата
			conn.Write([]byte("error format message\n"))
			continue
		}
		pos := strings.Index(message, " ") // нашли позицию разделителя

		if pos > 0 {
			outMessage := message[pos+1:] // отчистили сообщение от номера клиента

			if (*pull)[friend] == nil {
				conn.Write([]byte("client is close"))
				continue
			}

			out_buf := []byte(fmt.Sprintf("%s: %s\n", name, outMessage))

			// Отправить новую строку обратно клиенту
			_, errWrite := (*pull)[friend].Write(out_buf)

			// анализируем на ошибку
			if errWrite != nil {
				fmt.Println("Error:", errWrite.Error())
				break
			}
		}

	}
}

func pr(pull *map[string]net.Conn) string {
	var max int
	var list, answer string

	for i := range *pull {
		if max < len(i) {
			max = len(i)
		}
	}

	for i := range *pull {
		line := "* " + i
		for j := len(i); j < max; j++ {
			line += " "
		}
		list += line + " *\n"
	}

	for j := 0; j < max+4; j++ {
		answer += "*"
	}
	answer += "\n"

	answer = answer + list + answer
	return answer
}

func main() {
	// database
	log.Println("starting program")
	databaseUrl := "postgres://postgres:postgres@localhost:55001"
	dbPool, err := pgxpool.New(context.Background(), databaseUrl)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbPool.Close()
	ExecuteInsert(dbPool)
	// end

	// start server
	fmt.Println("Start server...")

	// create a pull connect
	pull := make(map[string]net.Conn, 1)

	// listen port
	ln, err := net.Listen("tcp", ":8081")

	// handler errors
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	defer ln.Close()

	// Запускаем цикл обработки соединений
	for {
		// Ждём входящее соединение
		conn, err := ln.Accept()

		if err != nil {
			panic(err.Error())
		}

		// сохраняем соединение в пул
		// pull[i] = conn

		// запускаем функцию process(conn)   как горутину
		go process(&pull, conn)
	}
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
//func ExecuteSelectQuery(dbPool *pgxpool.Pool) {
//	log.Println("starting execution of select query")
//
//	// execute the query and get result rows
//	rows, err := dbPool.Query(context.Background(), "select * from users")
//	if err != nil {
//		log.Fatal("error while executing query")
//	}
//
//	log.Println("result:")
//
//	// iterate through the rows
//	for rows.Next() {
//		values, err := rows.Values()
//		if err != nil {
//			log.Fatal("error while iterating dataset")
//		}
//
//		// convert DB types to Go types
//		id := values[0].(int32)
//		firstName := values[1].(string)
//		lastName := values[2].(string)
//		dateOfBirth := values[3].(string) // (string) // values[3].(time.Time)
//		log.Println("[id:", id, ", first_name:", firstName, ", last_name:", lastName, ", date:", dateOfBirth, "]")
//	}
//
//}
//func ExecuteFunction(dbPool *pgxpool.Pool, id int) {
//	log.Println("starting execution of database function")
//
//	// execute the query and get result rows
//	rows, err := dbPool.Query(context.Background(), "select * from get_ttt($1)", id)
//	log.Println("input id: ", id)
//	if err != nil {
//		log.Fatal("error while executing query")
//	}
//
//	log.Println("result:")
//
//	// iterate through the rows
//	for rows.Next() {
//		values, err := rows.Values()
//		if err != nil {
//			log.Fatal("error while iterating dataset")
//		}
//
//		//convert DB types to Go types
//		firstName := values[0].(string)
//		lastName := values[1].(string)
//		dateOfBirth := values[2].(time.Time)
//
//		log.Println("[first_name:", firstName, ", last_name:", lastName, ", date:", dateOfBirth, "]")
//	}
//
//}
