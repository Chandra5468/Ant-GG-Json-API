package main

import (
	"math/rand"
	"time"
)

type Account struct { // we mention `json:"something"` because, this is some kind of struct annotation. If this struct is getting serialized and coded to json, how that variable name should be
	// What is serialization and deserialization ?
	ID        int       `json:"id,omitempty"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Number    int32     `json:"number"`
	Balance   int16     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		// ID:        rand.Intn(10000), we will not use this random id, because in postgresql we are using serial type which increments on each creation. These random and serial might get conflict or give unorders id
		FirstName: firstName,
		LastName:  lastName,
		Number:    rand.Int31n(10000000), // Use this uuid as bank account number
		CreatedAt: time.Now().UTC(),
	}
}
