package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Goodbye struct {
	l *log.Logger
}

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}

func (g *Goodbye) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	g.l.Println("handle Goodbye requests")

	// Read request body data
	d, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(resp, "Oops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(resp, "Goodbye %s \n", d)
}
