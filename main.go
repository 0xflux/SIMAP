package main

import (
	"bufio"
	"encoding/base64"
	"log"
	"net"
	"os"
	"strings"
)

/**
* Simple IMAP server using no third party dependancies.
* @author 0xflux
 */

func main() {
	listener, err := net.Listen("tcp", ":143")
	if err != nil {
		log.Fatal(err)
	}

	defer listener.Close()

	log.Println("SIMAP Server listening on port 143")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection, %v\n", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	authenticated := false

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from connection, %v\n", err)
			return
		}

		// log.Printf("Received command: %v", line)
		cmd, args := parseCommand(line)

		if cmd == "LOGIN" && !authenticated {
			authenticated = processLogin(args, writer)
		} else if cmd == "PROCESSDATA" && authenticated {
			processData(args, writer)
		} else {
			sendResponse(writer, "BAD Command unrecognised or not allowed.")
		}
	}
}

func sendResponse(writer *bufio.Writer, message string) {
	writer.WriteString(message + "\r\n")
	writer.Flush()
}

func parseCommand(line string) (string, string) {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return strings.ToUpper(strings.TrimSpace(parts[0])), strings.TrimSpace(parts[1])
}

/*
To process the login, set the following environment variables on your local machine (or use docker preconfigured)
in order to have credentials you may login with via plaintext.
*/
func processLogin(args string, writer *bufio.Writer) bool {
	// extract username and password
	credentials := strings.Fields(args)
	if len(credentials) != 2 {
		sendResponse(writer, "BAD Invalid arguments for LOGIN")
		return false
	}

	// log.Printf("Login attempt with args: %s", args)

	// set local variables from args into func
	username, password := credentials[0], credentials[1]

	// get environment variables
	validUsername := os.Getenv("simap_poc_username")
	validPassword := os.Getenv("simap_poc_password")

	if validUsername == "" || validPassword == "" {
		log.Fatal("Username or password env variable not set. If it is set, reload your shell / session.")
	}

	if username == validUsername && password == validPassword {
		sendResponse(writer, "OK LOGIN successful")
		return true
	} else {
		sendResponse(writer, "BAD LOGIN invalid username/password")
		return false
	}
}

func processData(args string, writer *bufio.Writer) {
	decoded, err := base64.StdEncoding.DecodeString(args)
	if err != nil {
		sendResponse(writer, "BAD Error in decoding base64 data")
		return
	}

	log.Printf("Processed data: %v\n", string(decoded))
	sendResponse(writer, "OK Data processed successfully")
}
