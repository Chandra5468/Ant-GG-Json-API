package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJson(w http.ResponseWriter, status int, v any) error { // in v we will always send pointer
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	// w.Header().Set("Content-Type", "application/json") // Its not set method. It should be add method
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string
}

type apiFunc func(http.ResponseWriter, *http.Request) error // We are creating our own function type. This is used for passing a regular function as handler and using below httpHandlerFunc we will return a handler func

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// Handle error here
			WriteJson(w, http.StatusBadRequest, &ApiError{Error: err.Error()})
		}
	}
}

type APIServer struct {
	listenAddr string
}

func NewApiServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	// router := http.NewServeMux()

	// http.ListenAndServe("127.0.0.1:5000", router)

	// Using the most minimalistic routing package. Gorilla Mux

	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount)) // we are transforming handleAccount a normal function to http Handler.

	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.handleGETAccount)) // Q : Why are all the handlers in js(nodejs) and in go used inside a router without () this function braces.

	log.Println("JSON API Server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error { // In golang it is a standart to write handle as prefix to all controllers

	if r.Method == "GET" {
		return s.handleGETAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("methods not allowed %s", r.Method)
}

func (s *APIServer) handleGETAccount(w http.ResponseWriter, r *http.Request) error {

	id := mux.Vars(r)["id"]

	// Make database call here
	log.Println("This is id which got called ", id)
	return WriteJson(w, http.StatusOK, &Account{})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
