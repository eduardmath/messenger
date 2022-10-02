package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

const buffer int = 1024

var name string

// прием данных из сокета и вывод на печать
func readSock(conn *net.Conn) {

	if conn == nil {
		panic("Connection is nil")
	}

	buf := make([]byte, buffer)
	eof_count := 0

	for {
		// чистим буфер
		for i := 0; i < buffer; i++ {
			buf[i] = 0
		}
		readedLen, err := (*conn).Read(buf)
		fmt.Println(*conn, name)
		if err != nil {
			if err.Error() == "EOF" {
				eof_count++
				time.Sleep(time.Second * 2)
				checkConn(conn)
				fmt.Println(*conn, name)
				(*conn).Write([]byte(name)[:len(name)])
				//if eof_count > 7 {
				//	fmt.Println("Timeout connection")
				//	break
				//}
				continue
			}
			if strings.Index(err.Error(), "use of closed network connection") > 0 {

				fmt.Println("connection not exist or closed")
				continue
			}
			panic(err.Error())
		}
		if readedLen > 0 {
			fmt.Println(string(buf))
		}

	}
}

// ввод данных с консоли и вывод их в канал
func readConsole(ch chan string) {
	for {
		line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		if len(line) > buffer {
			fmt.Println("Error: message is very large")
			continue
		}
		if len(name) == 0 {
			name = line
		}
		// fmt.Print(">")
		out := line[:len(line)-1] // убираем символ возврата каретки
		ch <- out                 // отправляем данные в канал
	}
}

func main() {
	ch := make(chan string)

	defer close(ch) // закрываем канал при выходе из программы

	// readConsole(ch)

	var conn net.Conn
	checkConn(&conn)
	fmt.Println(conn)

	fmt.Print("Enter your name: ")

	go readConsole(ch)
	go readSock(&conn)

	for {
		val, ok := <-ch
		if ok { // если есть данные, то их пишем в сокет
			// val_len := len(val)
			out := []byte(val)
			_, err := conn.Write(out)
			if err != nil {
				fmt.Println("Write error:", err.Error())
				break
			}

		} else {
			// данных в канале нет, задержка на 2 секунды

			time.Sleep(time.Second * 2)
		}

	}
	fmt.Println("Finished...")

	conn.Close()
}

func checkConn(conn *net.Conn) {
	var err error
	for i := 1; ; i *= 2 {
		fmt.Print("Connecting")
		for j := 1; j <= 3; j++ {
			time.Sleep(time.Second / 2)
			fmt.Print(".")
		}
		time.Sleep(time.Second)
		fmt.Println()

		*conn, err = net.Dial("tcp", "127.0.0.1:8081")
		if err == nil {
			fmt.Println("Connection is established")
			return
		}

		fmt.Print("The connection is not established, please wait ")
		fmt.Println(i, "s.")
		time.Sleep(time.Second * time.Duration(i))
	}
}
