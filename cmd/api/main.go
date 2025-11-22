package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type user struct {
	Name string `json:"name"`
	Age int `json:"age"`
	City string `json:"city"`
}

func main() {
	port := ":3000"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Write([]byte("Hello from teachers!"))
		case http.MethodPost:
			// parse form data (necessary for x-www-form-urlencoded)
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error parsing form data!", http.StatusBadRequest)
				return
			}

			fmt.Println(r.Form)

			// Process the raw body
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}
			defer r.Body.Close()

			// Unmarshal raw json data
			var userInstance user
			err = json.Unmarshal(body, &userInstance)
			if err != nil {
				return
			}
			fmt.Println(userInstance)
		}
		
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from students"))
	})

	fmt.Println("Server running on port:", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Error starting new server: ", err)
	}

}