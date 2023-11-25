package main

import (
	"bytes"
	"io"
	"net"
	"regexp"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})
}

func TestMailServer(t *testing.T) {
	in := &bytes.Buffer{}
	out := &bytes.Buffer{}

	timeout, err := time.ParseDuration("10s")
	require.NoError(t, err)

	address := net.JoinHostPort("smtp.mail.com", "587")
	client := NewTelnetClient(address, timeout, io.NopCloser(in), out)
	require.NoError(t, client.Connect())
	defer func() { require.NoError(t, client.Close()) }()

	err = client.Receive()
	require.NoError(t, err)

	pattern := "^220 mail.com (.*) Nemesis ESMTP Service ready\r\n$"
	regExp, errRegExp := regexp.Compile(pattern)
	require.NoError(t, errRegExp)

	responseConnect := out.String()
	require.Regexp(t, regExp, responseConnect)

	regExp = regexp.MustCompile("220 smtp.mail.com ESMTP ")
	parts := regExp.Split(responseConnect, -1)
	require.Len(t, parts, 1)

	t.Run("HELLO server", func(t *testing.T) {
		out.Truncate(0)

		in.WriteString("HELO hi google\r\n")
		err = client.Send()
		require.NoError(t, err)

		err = client.Receive()
		require.NoError(t, err)
		regExp := regexp.MustCompile("250 mail.com Hello hi google \\[.*\r\n")
		require.Regexp(t, regExp, out.String())
	})

	t.Run("EHLO server", func(t *testing.T) {
		out.Truncate(0)

		in.WriteString("EHLO mail.com\r\n")
		err = client.Send()
		require.NoError(t, err)

		err = client.Receive()
		require.NoError(t, err)
		pattern := "250-mail.com Hello mail.com \\[.*\r\n250-8BITMIME\r\n250-SIZE 141557760\r\n250 STARTTLS\r\n$"
		regExp := regexp.MustCompile(pattern)
		require.Regexp(t, regExp, out.String())
	})
}
