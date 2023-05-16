package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"unicode/utf8"
)

func getMsg(conn net.Conn) (string, error) {
	buffer := make([]byte, 2048)
	l, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}
	msg := strings.TrimSpace(string(buffer[:l]))
	log.Println(fmt.Sprintf("Received from %s: %s", conn.RemoteAddr(), msg))
	return msg, nil
}

func sendMsg(conn net.Conn, response string) {
	log.Println(fmt.Sprintf("Send to %s: %s", conn.RemoteAddr(), response))
	_, err := conn.Write([]byte(fmt.Sprintf("%s\n", response)))
	if err != nil {
		panic(err)
	}
}

func parseCommand(input string) (string, string, error) {
	var cmd, args string

	if utf8.RuneCountInString(input) < 3 {
		return cmd, args, errors.New(fmt.Sprintf("%d %s", StatusSyntaxError, "Syntax error, command unrecognized."))
	}

	inputs := strings.SplitAfterN(input, " ", 2)
	if len(inputs) == 2 {
		cmd = strings.TrimSpace(inputs[0])
		args = strings.TrimSpace(inputs[1])
	} else if len(inputs) == 1 {
		cmd = strings.TrimSpace(inputs[0])
	}

	return cmd, args, nil
}
