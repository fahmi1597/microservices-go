package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	h.l.Println("handle Hello requests")

	// Read request body
	d, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.l.Println("Error reading body", err.Error())

		http.Error(resp, "Unable to read request body", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(resp, "Hello %s \n", d)
}
