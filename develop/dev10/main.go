package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type App struct {
	timeout time.Duration
	host    string
	port    string
	err     error
}

type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

func main() {
	app := App{
		host: "127.0.0.1",
		port: "8080",
	}
	app.parseFlags()

	conn, err := app.connectTCP()
	if err != nil {
		fmt.Printf("Ошибка при создании сокета: %v", err)
		return
	}

	for {
		if err = WriteToTCP(conn); err != nil {
			break
		}
		if err = ReadFromTCP(conn); err != nil {
			break
		}
	}

	fmt.Printf("\nПрограмма остановлена: %v\n", err)
}

func (a *App) parseFlags() error {
	var timeout stringFlag
	var err error

	flag.Var(&timeout, "timeout", "таймаут для соединения")
	flag.Parse()

	dur, err := time.ParseDuration(timeout.String())
	if err != nil {
		return fmt.Errorf("Error on reading timeout: %v", err)
	}

	sliceStart := 0

	if timeout.set {
		sliceStart = 1
	}

	fmt.Print(sliceStart)
	args := flag.Args()[sliceStart:]

	if len(args) > 0 {
		a.host = args[0]
		fmt.Printf("%v", args)
	}

	if len(args) > 1 {
		a.port = args[1]
	}

	a.timeout = dur
	return nil
}

func (a *App) connectTCP() (net.Conn, error) {
	fmt.Printf("\n%v:%v", a.host, a.port)
	conn, err := net.DialTimeout("tcp", a.host+":"+a.port, a.timeout)
	if err != nil {
		return nil, fmt.Errorf("Ошибка: %v", err)
	}

	fmt.Printf("Соединение с сервером %v:%v установлено", a.host, a.port)
	return conn, nil
}

func WriteToTCP(c net.Conn) error {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\nТекст для отправки: ")
	text, err := reader.ReadString('\n')

	if err != nil {
		if err == io.EOF {
			c.Close()
			return fmt.Errorf("Выход из программы")
		}
		return fmt.Errorf("Ошибка при чтении из консоли: %v", err)
	}

	fmt.Fprintf(c, text+"\n")
	return nil
}

func ReadFromTCP(c net.Conn) error {
	line, _, err := bufio.NewReader(c).ReadLine()
	if err != nil {
		return fmt.Errorf("Ошибка при чтении из сокета: %v", err)
	}

	fmt.Print("Ответ от сервера: " + string(line))
	return nil
}
