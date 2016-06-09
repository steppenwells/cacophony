package main

import(
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"os"
	"io"
	"bufio"
	"github.com/google/uuid"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}

func uploadWav(wr http.ResponseWriter, r *http.Request) {
	fileName := "./audio/" + uuid.New().String() + ".wav"
	f, _ := os.Create(fileName)
	defer f.Close()
	defer r.Body.Close()

	w := bufio.NewWriter(f)
	io.Copy(w, r.Body)
	w.Flush()
	f.Sync()

	fmt.Fprintf(wr, fileName)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	r.HandleFunc("/uploadWav", uploadWav).Methods("POST")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)

	http.ListenAndServe(":5000", nil)
}
