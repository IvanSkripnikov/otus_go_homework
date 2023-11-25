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

	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "Hello server",
			message:  "HELO hi google\r\n",
			expected: "250 mail.com Hello hi google \\[.*\\]\r\n",
		},
		{
			name:    "Ehlo server",
			message: "EHLO mail.com\r\n",
			//nolint:lll
			expected: "250-mail.com Hello mail.com \\[.*\\]\r\n250-8BITMIME\r\n250-SIZE 141557760\r\n250 STARTTLS\r\n$",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			out.Truncate(0)

			in.WriteString(tc.message)
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			regExp, errRegExp := regexp.Compile(tc.expected)
			require.NoError(t, errRegExp)
			require.Regexp(t, regExp, out.String())
		})
	}
}
