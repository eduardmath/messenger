package main

import (
	"fmt"
	"net"
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
	fmt.Println("Accept user:", name, len(name))
	defer conn.Close()

	for {
		readed_len, err := conn.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				fmt.Println("Close ", name)
				delete(*pull, name)
				break
			}
			fmt.Println(err)
		}

		// Распечатываем полученое сообщение
		//fmt.Println("Received Message:", read_len, buf)
		var message = string(buf[:readed_len])

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
		// friend = message[:pos]
		// _f, err = fmt.Sscan(message[:pos], &friend) // определи номер клиента

		if pos > 0 {
			out_message := message[pos+1:] // отчистили сообщение от номера клиента
			// Распечатываем полученое сообщение

			// if buf[0] == 0 {

			// !!!fatal error!!!
			// conn = (*pull)[friend]
			if (*pull)[friend] == nil {
				conn.Write([]byte("client is close"))
				continue
			}

			// }
			out_buf := []byte(fmt.Sprintf("%s: %s\n", name, out_message))

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
		li := "* " + i

		fmt.Println((*pull)[i])

		for j := len(i); j < max; j++ {
			li += " "
		}
		list += li + " *\n"
	}

	for j := 0; j < max+4; j++ {
		answer += "*"
	}
	answer += "\n"

	answer = answer + list + answer
	return answer
}

func main() {
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
