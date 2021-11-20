package client

import (
	"errors"
	"net"
	"net/textproto"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8964"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
				server := textproto.NewConn(conn)
				defer server.Close()
				server.Writer.PrintfLine("220 Service ready for new user.")
			}
		}()
	}

	client, _ := NewFtpClient("localhost:8964")
	if client == nil {
		t.Fatal("client is nil")
	}
}

func TestLogin(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8965"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
				server := textproto.NewConn(conn)
				defer server.Close()
				server.Writer.PrintfLine("220 Service ready for new user.")

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
					server.Writer.PrintfLine("331 User name okay, need password.")
				}

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASS") {
					server.Writer.PrintfLine("230 User logged in, proceed.")
				}
			}
		}()
	}

	client, _ := NewFtpClient("localhost:8965")
	if client == nil {
		t.Fatal("client is nil")
	}

	if err := client.Login("user", "pass"); err != nil {
		t.Fatal(err)
	}
}

func TestUsernameNotExist(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8966"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
				server := textproto.NewConn(conn)
				defer server.Close()
				server.Writer.PrintfLine("220 Service ready for new user.")

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
					server.Writer.PrintfLine("332 Need account for login.")
				}
			}
		}()
	}

	client, _ := NewFtpClient("localhost:8966")
	if client == nil {
		t.Fatal("client is nil")
	}

	if err := client.Login("user", "pass"); err == nil || !errors.Is(err, ErrUsernameNotExist) {
		t.Fatal("should not login")
	}
}

func TestPasswordError(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8967"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
				server := textproto.NewConn(conn)
				defer server.Close()
				server.Writer.PrintfLine("220 Service ready for new user.")

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
					server.Writer.PrintfLine("331 User name okay, need password.")
				}

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASS") {
					server.Writer.PrintfLine("530 Not logged in.")
				}
			}
		}()
	}

	client, _ := NewFtpClient("localhost:8967")
	if client == nil {
		t.Fatal("client is nil")
	}

	if err := client.Login("user", "pass"); err == nil || !errors.Is(err, ErrPasswordNotMatch) {
		t.Fatal("should not login")
	}
}

func TestLogout(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8968"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
				server := textproto.NewConn(conn)
				defer server.Close()
				server.Writer.PrintfLine("220 Service ready for new user.")

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "USER") {
					server.Writer.PrintfLine("331 User name okay, need password.")
				}

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "PASS") {
					server.Writer.PrintfLine("230 User logged in, proceed.")
				}

				if line, _ := server.ReadLine(); strings.HasPrefix(line, "QUIT") {
					server.Writer.PrintfLine("221 Service closing control connection.")
				}
			}
		}()
	}

	client, err := NewFtpClient("localhost:8968")
	if err != nil {
		t.Fatal(err)
	}
	if err = client.Login("user", "pass"); err != nil {
		t.Fatal(err)
	}
	if err := client.Logout(); err != nil {
		t.Fatal(err)
	}
}

func TestMode(t *testing.T) {
	if listener, err := net.Listen("tcp", ":8969"); err != nil {
		t.Fatal(err)
	} else {
		go func() {
			defer listener.Close()
			if conn, err := listener.Accept(); err != nil {
				return
			} else {
				server := textproto.NewConn(conn)
				defer server.Close()
				server.Writer.PrintfLine("220 Service ready for new user.")

				for {
					if line, _ := server.ReadLine(); strings.HasPrefix(line, "MODE S") ||
						strings.HasPrefix(line, "MODE C") {
						server.Writer.PrintfLine("200 Command okay.")
					} else if strings.HasPrefix(line, "MODE B") {
						server.Writer.PrintfLine("504 Command not implemented for that parameter.")
					}
				}
			}
		}()
	}

	client, err := NewFtpClient("localhost:8969")
	if err != nil {
		t.Fatal(err)
	}
	if err = client.Mode(ModeStream); err != nil {
		t.Fatal(err)
	}
	if err = client.Mode(ModeBlock); err == nil || !errors.Is(err, ErrModeNotSupported) {
		t.Fatal("should not change mode")
	}
	if err = client.Mode(ModeCompressed); err != nil {
		t.Fatal(err)
	}
}