package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var developerToken string

func main() {
	// Set your developer token here
	err := godotenv.Load("../.env")
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	developerToken = os.Getenv("AppleMusicDevToken")

	// Define a file server to serve static files (HTML, CSS, JS)
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// Endpoint to serve the developer token
	http.HandleFunc("/developer-token", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request for developer token")
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(developerToken))
		fmt.Println("Sent developer token", developerToken)

	})
	http.HandleFunc("/user-token", func(w http.ResponseWriter, r *http.Request) {
		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		// Print the received user token
		fmt.Println("Received user token:", string(body))

		// Respond with success
		w.WriteHeader(http.StatusOK)
	})

	// Start the HTTP server on port 8080
	println("Server is listening on port 8080...")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
