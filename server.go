package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ehames/dhtstore"
	"github.com/gorilla/mux"
)

var store = map[string]string{}

func Start() {
	router := mux.NewRouter()
	router.HandleFunc("/map/{key}", put).Methods("PUT")
	router.HandleFunc("/map/{key}", get).Methods("GET")
	http.ListenAndServe(":9000", router)
}

func put(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["key"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Couldn't read body: %s\n", err)
		panic(err)
	}
	entry := new(dhtstore.Entry)
	json.Unmarshal(body, &entry)
	store[entry.Key] = entry.Value
	dhtstore.Put(key, entry.Value)
}

func get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	key := vars["key"]

	value := store[key]

	rv := dhtstore.Entry{Key: key, Value: value}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(rv); err != nil {
		panic(err)
	}
}
