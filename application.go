package main

import(
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)

	http.ListenAndServe(":5000", nil)
}
