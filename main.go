package main

import (
	"fmt"
	"net/http"
	"os"
)

var hn string

func init() {
	hn, _ = os.Hostname()
}

func main() {
	h := &H{}
	http.HandleFunc("/", h.ServeHTTP)
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

type H struct {
}

func (h *H) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, e := fmt.Fprintf(w, "hello %s!", "envlibs")
	if e != nil {
		println(e.Error())
	}
}
