package main

import (
	"fmt"
	"net/http"
)

func appruve(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Header)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/update/", appruve)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
