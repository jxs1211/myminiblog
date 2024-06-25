package coding

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func runHttpServer() {
	err := http.ListenAndServe(":8000", handler())
	if err != nil {
		log.Fatal(err)
	}
}

func handler() http.Handler {
	h := http.NewServeMux()
	h.Handle("/double", http.HandlerFunc(doubleHandler))
	return h
}

func doubleHandler(w http.ResponseWriter, r *http.Request) {
	v := r.FormValue("v")
	if len(v) == 0 {
		http.Error(w, "no value", http.StatusBadRequest)
		return
	}
	val, err := strconv.Atoi(v)
	if err != nil {
		http.Error(w, "invalid value", http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, val*2)
}
