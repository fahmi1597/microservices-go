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
	h.l.Println("Handled - Hello World")
	d, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, "Oops", http.StatusBadRequest)
		// resp.WriteHeader(http.StatusBadRequest)
		// resp.Write([]byte("Bad Request"))
		return
	}

	fmt.Fprintf(resp, "Hello %s \n", d)
}
