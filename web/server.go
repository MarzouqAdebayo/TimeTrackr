package web

import (
	boltdb "TimeTrackr/boltDB"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

var (
	PORT   = 8080
	o, err = os.Getwd()
	tmpl   = template.Must(template.ParseFiles(filepath.Join(o, "web/index.html")))
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, nil)
}

func startTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		taskName := r.FormValue("taskName")
		category := r.FormValue("category")
		fmt.Fprintf(w, "Started task: %s in category: %s", taskName, category)
		boltdb.StartTask(taskName, category)
	}
}

func displayTasksHandler(w http.ResponseWriter, r *http.Request) {
	res, err := boltdb.FilterTasks(boltdb.FilterObject{})
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Printf("An erro occurred: %s", err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
	w.WriteHeader(http.StatusOK)
}

func RunServer(port *int) {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/start", startTaskHandler)
	http.HandleFunc("/tasks", displayTasksHandler)

	if port != nil {
		PORT = *port
	}

	fmt.Printf("Server running on localhost:%d", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
