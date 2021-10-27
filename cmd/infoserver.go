package cmd

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opxyc/wdc/alert"
)

// runHTTPServer creates a simple http server that listens on localhost:8080
// with single route /{id} to server details of an alert with given id
func runHTTPServer() {
	r := mux.NewRouter()
	r.HandleFunc("/{id}", alertInfoHandler)
	fmt.Println("[+] Info server started")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func alertInfoHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logID := vars["id"]
	alert, err := alert.ReadFromLog(logDir, logID)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(rw, "%v", err)
		return
	}

	t, err := template.ParseFiles("cmd/templates/alertinfo.html")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(rw, "%v", err)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/html")
	t.Execute(rw, &alert)
}
