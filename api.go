package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Calling jwt auth middleware")
		tokenString := r.Header.Get("Authorization")
		token, err := validateJWT(tokenString)

		if err != nil {
			WriteJson(w, http.StatusForbidden, &ApiError{Error: "invalid token permission denied"})
			return
		}
		// Get user Id from user
		idstr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idstr)
		if err != nil {
			WriteJson(w, http.StatusForbidden, &ApiError{Error: "invalid id permission denied"})
			return
		}
		account, err := s.GetAccountByID(id)
		if err != nil {
			WriteJson(w, http.StatusForbidden, &ApiError{Error: "Account id not found in db, permission denied"})
			return
		}
		claims := token.Claims.(jwt.MapClaims)

		if account.Number != claims["accountNumber"] {
			WriteJson(w, http.StatusForbidden, &ApiError{Error: "Permission denied"})
			return
		}

		handlerFunc(w, r)
	}
}

const secret = "hunter9999" // store it in envs

func createJWT(account *Account) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"accountNumber": account.Number,
		"exp":           time.Now().Add(time.Hour * 24).Unix(), // expires after 24 hrs
	})
	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	// secret := os.Getenv("jwtsecret")
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token, nil
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
	store      Storage
}

func NewApiServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	// router := http.NewServeMux()

	// http.ListenAndServe("127.0.0.1:5000", router)

	// Using the most minimalistic routing package. Gorilla Mux

	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.handleAccount)) // we are transforming handleAccount a normal function to http Handler.

	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandlerFunc(s.handleGETAccountByID), s.store)) // Q : Why are all the handlers in js(nodejs) and in go used inside a router without () this function braces.
	router.HandleFunc("/transfer/", makeHTTPHandlerFunc(s.handleTransfer))
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
	return fmt.Errorf("methods not allowed %s", r.Method)
}

// get /account
func (s *APIServer) handleGETAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, &accounts)
}

func (s *APIServer) handleGETAccountByID(w http.ResponseWriter, r *http.Request) error {

	if r.Method == http.MethodGet {

		idstr := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idstr)
		if err != nil {
			return fmt.Errorf("invalid id given %s", idstr)
		}
		// Make database call here
		// log.Println("This is id which got called ", id)

		account, err := s.store.GetAccountByID(id)

		if err != nil {
			return err
		}
		return WriteJson(w, http.StatusOK, account)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method out of scope for this request")
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	// createAccountReq := &CreateAccountRequest{

	// }
	/*
		p := new(chan int)   // p has type: *chan int
		c := make(chan int)  // c has type: chan int


		or you can also use

		createAccountReq := &CreateAccountRequest{}
	*/

	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}
	tokenString, err := createJWT(account)
	if err != nil {
		return err
	}

	log.Println("JWT TOKEN ", tokenString)
	return WriteJson(w, http.StatusCreated, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	// _, err := s.store.DeleteAccount()
	idstr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idstr)
	if err != nil {
		return fmt.Errorf("invalid id given %s", idstr)
	}
	err = s.store.DeleteAccount(id)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	defer r.Body.Close() // This should be the approach of closing the body after reading from it
	transferReq := &TransferRequest{}

	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, transferReq)
}
