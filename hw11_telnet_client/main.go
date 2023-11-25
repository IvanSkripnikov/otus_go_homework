package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second*10, "timeout of closing")
}

func main() {
	flag.Parse()

	if !validateRequest(os.Args) {
		log.Fatalln("Required arguments has missing (host or port)")
	}

	ctxNotify, stopFunc := signal.NotifyContext(context.Background(), syscall.SIGINT)
	defer stopFunc()

	ctxTimeout, cancelFunc := context.WithTimeout(ctxNotify, timeout)
	defer cancelFunc()

	// получаем хост и порт.
	countArgs := len(os.Args)
	address := net.JoinHostPort(os.Args[countArgs-2], os.Args[countArgs-1])
	telnet := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	defer func() {
		telnet.Close()
	}()

	err := telnet.Connect()
	if err != nil {
		log.Println(err)
		return
	}

	// запускаем горутину для отправки данных в сокет.
	go func() {
		for {
			select {
			case <-ctxTimeout.Done():
				return
			default:
				err = telnet.Send()
				if checkError(err) {
					log.Printf("Error while sending data to socket, %v", err)
				}

				if errors.Is(err, io.EOF) {
					stopFunc()
				}
			}
		}
	}()

	// запускаем горутину для чтения данных из сокета.
	go func() {
		for {
			select {
			case <-ctxTimeout.Done():
				return
			default:
				err = telnet.Receive()
				if checkError(err) {
					log.Printf("Error while get data from socket, %v", err)
				}

				if errors.Is(err, io.EOF) {
					stopFunc()
				}
			}
		}
	}()

	<-ctxTimeout.Done()
	err = ctxTimeout.Err()
	if checkError(err) {
		log.Println(err)
	}
}

// проверяем наличие флага.
func checkExistsFlags(args []string) (int, bool) {
	n := 0
	found := false
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			n++
			found = true
		}
	}

	return n, found
}

// проверяем входные данные.
func validateRequest(args []string) bool {
	countArgs := len(args)

	if n, ok := checkExistsFlags(args); ok {
		countArgs -= n
	}

	return countArgs >= 3
}

// перенести данные из одного источника в другой.
func transferData(in io.ReadCloser, out io.Writer) error {
	buffer := make([]byte, 1024)
	n, err := in.Read(buffer)
	if checkError(err) {
		return fmt.Errorf("read error: %w", err)
	}

	_, err = out.Write(buffer[:n])
	if err != nil {
		return fmt.Errorf("write error: %w", err)
	}

	return err
}

// проверяем, действительно ли это ошибка.
func checkError(err error) bool {
	return err != nil && !errors.Is(err, io.EOF)
}
