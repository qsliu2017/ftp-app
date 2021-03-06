package server

import (
	"io"
	"net"
	"net/textproto"
	"strings"
)

type clientHandler struct {
	ctrl *textproto.Conn
	conn net.Conn

	login    bool
	username string

	mode  byte
	type_ byte
	stru  byte

	rootDir string
}

func handleClient(conn net.Conn, rootDir string) {
	defer conn.Close()
	handler := &clientHandler{
		ctrl: textproto.NewConn(conn),

		mode:  ModeStream,
		type_: TypeAscii,
		stru:  StruFile,

		rootDir: rootDir,
	}

	handler.reply(StatusReady)

	for {
		cmd, err := handler.ctrl.ReadLine()
		if err != nil {
			if err == io.EOF {
				logger.Printf("%s:%s disconnected", conn.RemoteAddr(), handler.username)
				break
			} else {
				logger.Printf("read command error: %v", err)
				continue
			}
		}

		var param string
		part := strings.SplitN(cmd, " ", 2)
		if len(part) != 2 {
			param = ""
		} else {
			param = part[1]
		}

		// Sometime it read a empty line. Skip it.
		if len(part[0]) < 4 {
			continue
		}

		if cmdHandler, has := commandHandlers[part[0]]; has {
			logger.Printf("%s:%s %s", conn.RemoteAddr(), handler.username, cmd)
			if err := cmdHandler(handler, param); err != nil {
				logger.Printf("%s:%s %s error: %v", conn.RemoteAddr(), handler.username, cmd, err)
				if err == ErrCloseConn {
					break
				}
			}
		} else {
			handler.reply(StatusSyntaxError)
		}
	}
}

type commandHandler func(c *clientHandler, param string) error

var commandHandlers = map[string]commandHandler{
	//auth commands
	"USER": (*clientHandler).handleUSER,
	"PASS": (*clientHandler).handlePASS,
	"QUIT": (*clientHandler).handleQUIT,

	//dial commands
	"PORT": (*clientHandler).handlePORT,
	"PASV": (*clientHandler).handlePASV,

	//file commands
	"RETR": (*clientHandler).handleRETR,
	"STOR": (*clientHandler).handleSTOR,

	//param commands
	"MODE": (*clientHandler).handleMODE,
	"TYPE": (*clientHandler).handleTYPE,
	"STRU": (*clientHandler).handleSTRU,
}
