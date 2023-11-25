package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Telnet{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

type Telnet struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (telnet *Telnet) Connect() error {
	conn, err := net.DialTimeout("tcp", telnet.address, telnet.timeout)
	if err != nil {
		return fmt.Errorf("cannot connect: %w", err)
	}

	telnet.conn = conn
	return nil
}

func (telnet *Telnet) Close() error {
	err := telnet.in.Close()
	if err != nil {
		return fmt.Errorf("cannot closing input: %w", err)
	}

	if telnet.conn != nil {
		err = telnet.conn.Close()
		if err != nil {
			return fmt.Errorf("cannot closing connect: %w", err)
		}
	}

	return nil
}

func (telnet *Telnet) Send() error {
	return transferData(telnet.in, telnet.conn)
}

func (telnet *Telnet) Receive() error {
	return transferData(telnet.conn, telnet.out)
}
