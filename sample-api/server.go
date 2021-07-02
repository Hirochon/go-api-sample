package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"work/middleware"
)

var addr = ":8080"

type helloJSON struct {
	UserName string `json:"user_name"`
	Content  string `json:"content"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "GET hello!\n")
	case "POST":
		body := r.Body
		defer body.Close()

		buf := new(bytes.Buffer)
		io.Copy(buf, body)

		var hello helloJSON
		json.Unmarshal(buf.Bytes(), &hello)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "POST hello! %s\n", hello)

	default:
		fmt.Fprint(w, "Method not allowed.!\n")
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/hello", Handler)
	fmt.Printf("[START] server. port: %s\n", addr)
	if err := http.ListenAndServe(addr, middleware.Log(router)); err != nil {
		panic(fmt.Errorf("[FAILED] start sever. err: %v", err))
	}
}
