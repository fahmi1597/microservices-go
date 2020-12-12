package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		log.Print("Hello World!")
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(resp, "Oops", http.StatusBadRequest)
			// resp.WriteHeader(http.StatusBadRequest)
			// resp.Write([]byte("Bad Request"))
			return
		}
		fmt.Fprintf(resp, "Hello %s \n", body)
	})
	http.HandleFunc("/goodbye", func(resp http.ResponseWriter, req *http.Request) {
		log.Print("Goodbye World")
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
