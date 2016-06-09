package main

import(
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"os"
	"io"
	"bufio"
	"github.com/google/uuid"
	"github.com/unixpickle/wav"
)

func healthcheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Ok")
}

func updateConversation(latest string) {
	convFileName := "./audio/converstion.wav"

	if _, err := os.Stat(convFileName); os.IsNotExist(err) {
		latestSound, _ := wav.ReadSoundFile(latest)
		wav.WriteFile(latestSound, convFileName)

	} else {
		convSound, _ := wav.ReadSoundFile(convFileName)
		latestSound, _ := wav.ReadSoundFile(latest)
		wav.Append(convSound, latestSound)
		wav.WriteFile(convSound, convFileName)
	}

}

func saveLatest(filename string, contents io.Reader) {
	f, _ := os.Create(filename)
	defer f.Close()

	w := bufio.NewWriter(f)
	io.Copy(w, contents)
	w.Flush()
	f.Sync()

}

func uploadWav(wr http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	fileName := "./audio/" + uuid.New().String() + ".wav"
	saveLatest(fileName, r.Body)
	updateConversation(fileName)

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
