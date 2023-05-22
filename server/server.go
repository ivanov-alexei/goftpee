package server

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Address string  `json:"address"`
	Port    int32   `json:"port"`
	Users   []*User `json:"users"`
}

type Session struct {
	User        *User
	DataAddress string
}

type FtpServer struct {
	Config *Config
}

func NewFtpServer(config *Config) *FtpServer {
	return &FtpServer{
		Config: config,
	}
}

func (server *FtpServer) Start() error {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Config.Address, server.Config.Port))
	if err != nil {
		return err
	}

	log.Println("FTP Server was running on port", server.Config.Port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go handleConnection(conn, server.Config.Users)
	}
}

// func handleConnection(conn net.Conn, dataListener net.Listener, users []*User) {
func handleConnection(conn net.Conn, users []*User) {

	log.Println("The connection is established with ", conn.RemoteAddr())

	// sent "ready" message
	sendMsg(conn, fmt.Sprintf("%d %s", StatusServiceReady, "Service ready for new user."))
	user := User{}
	session := &Session{
		User: &user,
	}

	// user auth
	for {
		msg, err := getMsg(conn)
		if err != nil {
			panic(err)
		}
		response := handleLogin(msg, &user, users)
		sendMsg(conn, response)
		if user.valid {
			log.Printf("User %s is logged in.\n", user.Username)
			break
		}
	}

	// if user is logged, server is ready to get message from user
	for {
		input, err := getMsg(conn)
		if err != nil {
			log.Println("Something went wrong", err)
			break
		}

		handleCommand(session, input, conn)
	}
}

func handleCommand(session *Session, input string, conn net.Conn) {
	cmd, args, err := parseCommand(input)

	switch {
	case cmd == "QUIT":
		sendMsg(conn, fmt.Sprintf(ResponseTemplate, StatusCloseControlConn, "Service closing control connection."))
		if err = conn.Close(); err != nil {
			panic(err)
		}
		log.Println("The connection was closed by", conn.RemoteAddr())
	case cmd == "TYPE" && (args == "A" || args == "I"):
		sendMsg(conn, fmt.Sprintf("%d %s", StatusCommandOk, "Type is OK."))
	case cmd == "PORT":
		connectData := strings.Split(args, ",")
		n1, _ := strconv.Atoi(connectData[4])
		n2, _ := strconv.Atoi(connectData[5])
		port := n1*256 + n2
		session.DataAddress = fmt.Sprintf("%s.%s.%s.%s:%d", connectData[0], connectData[1], connectData[2], connectData[3], port)
		sendMsg(conn, fmt.Sprintf("%d %s.\n", StatusCommandOk, "Command Ok."))
	//not min command, but I'd like to add it into my code
	case cmd == "LIST":
		dataConn, closeDataConn := handleDataConnection(session, conn)
		files, err := os.ReadDir(".")
		if err != nil {
			panic(err)
		}

		var msg string
		for _, file := range files {
			msg += file.Name() + "\n"
		}
		sendMsg(dataConn, fmt.Sprintf("%s", msg))

		closeDataConn()

	case cmd == "MODE":
		sendMsg(conn, fmt.Sprintf(ResponseTemplate, StatusCommandOk, ""))
	case cmd == "STRU":
		sendMsg(conn, fmt.Sprintf(ResponseTemplate, StatusCommandOk, ""))
	case cmd == "RETR":
		dataConn, closeDataConn := handleDataConnection(session, conn)

		data, err := os.ReadFile(args)
		if err != nil {
			dataConn.Close()
			conn.Close()
			panic(err)
		}
		sendMsg(dataConn, string(data))

		closeDataConn()

	case cmd == "STOR":
		dataConn, closeDataConn := handleDataConnection(session, conn)

		f, err := os.Create(args)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		data, err := getMsg(dataConn)
		_, err = f.WriteString(data)

		closeDataConn()
	case cmd == "NOOP":
		sendMsg(conn, fmt.Sprintf(ResponseTemplate, StatusCommandOk, ""))

		/*	****** not min impl, commands for further implementation ******

			case cmd == "CWD":
				log.Println(args)
				// todo здесь полное гавно,надо как-то сделать отслеживание текущей и родительской директории
				// todo плюс надо выдавать права на определенный корень, чтоб нельзя было
				if strings.HasPrefix(args, "/") {
					sendMsg(conn, "530 Not logged in for that operation")
					} else {
						err := os.Chdir(args)
						if err != nil {
							sendMsg(conn, "501 Syntax error.")
							} else {
								sendMsg(conn, "250 Requested file action okay, completed.")
							}
						}
			case cmd == "PASV":
				response := fmt.Sprintf("%d 127,0,0,1,0,20", StatusPassiveMode)
				sendMsg(conn, response)
			case cmd == "PWD":
				path, err := os.Getwd()
				if err != nil {
					panic(err)
				}
				sendMsg(conn, fmt.Sprintf("257 %s", path))
		*/
	default:
		sendMsg(conn, fmt.Sprintf(ResponseTemplate, StatusCommandNotImplemented, "Command not implemented."))
	}
}

func handleDataConnection(session *Session, controlConn net.Conn) (net.Conn, func ()) {
	if session.DataAddress == "" {
		sendMsg(controlConn, fmt.Sprintf(ResponseTemplate, StatusSyntaxError, "Need to get data address."))
	}
	conn1, err := net.Dial("tcp", session.DataAddress)
	if err != nil {
		conn1.Close()
		panic(err)
	}
	sendMsg(controlConn, "150 File status okay; about to open data connection.")

	return conn1, func () {
		sendMsg(controlConn, fmt.Sprintf(ResponseTemplate, StatusFileActionOk, "Data was transferred.\n"))
		if err != nil {
			controlConn.Close()
			panic(err)
		}
		log.Println(conn1.Close())
	}
}