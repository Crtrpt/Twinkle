package main

import (
	"encoding/json"
	"net/http"
)

func Handle(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]any, 0)
	data["RequestURI"] = r.RequestURI
	data["Method"] = r.Method
	data["Header"] = r.Header.Clone()

	dataStr, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Write([]byte(dataStr))
}

func main() {
	http.HandleFunc("/", Handle)
	http.ListenAndServe(":8088", nil)
}
