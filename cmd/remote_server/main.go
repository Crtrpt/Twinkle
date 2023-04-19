package main

import (
	"encoding/json"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	header, _ := json.Marshal(r.Header)
	w.Write([]byte(header))
}

func main() {
	http.HandleFunc("/", Handle)
	http.ListenAndServe(":8088", nil)
}
