package main

import (
	"fmt"
	"net/http"
	"time"
)

type handler struct {
}

func (s *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello world for exmaple at %s", time.Now().Format("2006-01-02 15:04:05.000"))
}
