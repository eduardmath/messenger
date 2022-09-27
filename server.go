package main

import (
	"fmt"
	"net"
	"strings"
)

func process(pull map[int]net.Conn, number int) {
	var clientNo int
	buf := make([]byte, 1024*8)

	// получаем доступ к текущему соединению
	conn := pull[number]

	// определим, что перед выходом из функции, мы закроем соединение
	fmt.Println("Accept cnn:", number)
	defer conn.Close()

	for {
		readed_len, err := pull[number].Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Close ", number)
				delete(pull, number)
				break
			}
			fmt.Println(err)
		}

		// Распечатываем полученое сообщение
		//fmt.Println("Received Message:", read_len, buf)
		var message = string(buf[:readed_len])

		// парсинг полученного сообщения
		_, err = fmt.Sscanf(message, "%d", &clientNo) // определи номер клиента
		if err != nil {
			// обработка ошибки формата
			conn.Write([]byte("error format message\n"))
			continue
		}
		pos := strings.Index(message, " ") // нашли позицию разделителя

		if pos > 0 {
			out_message := message[pos+1:] // отчистили сообщение от номера клиента
			// Распечатываем полученое сообщение

			// if buf[0] == 0 {
			conn = pull[clientNo]
			if conn == nil {
				pull[number].Write([]byte("client is close"))
				continue
			}

			// }
			out_buf := []byte(fmt.Sprintf("%d->>%s\n", clientNo, out_message))

			// Отправить новую строку обратно клиенту
			_, err2 := conn.Write(out_buf)

			// анализируем на ошибку
			if err2 != nil {
				fmt.Println("Error:", err2.Error())

				break // выходим из цикла
			}
		}

	}
}

func main() {
	// start server
	fmt.Println("Start server...")

	// create a pull connect
	pull := make(map[int]net.Conn, 1)

	// listen port
	ln, err := net.Listen("tcp", ":8081")

	// handler errors
	if err != nil {
		fmt.Println(err.Error())
		panic(err.Error())
	}
	defer ln.Close()

	fmt.Println("error before for")

	// Запускаем цикл обработки соединений
	for i := 0; ; i++ {
		// Принимаем входящее соединение
		conn, err := ln.Accept()

		if err != nil {
			panic(err.Error())
		}

		// сохраняем соединение в пул
		pull[i] = conn

		// запускаем функцию process(conn)   как горутину
		go process(pull, i)
	}
}
