package main

import "math/rand"

type Account struct { // we mention `json:"something"` because, this is some kind of struct annotation. If this struct is getting serialized and coded to json, how that variable name should be
	// What is serialization and deserialization ?
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int32  `json:"number"`
	Balance   int16  `json:"balance"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.Intn(1000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    rand.Int31n(10000000), // Use this uuid as bank account number
	}
}
