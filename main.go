// This is a bank application

package main

import "log"

// Not creating much folder structures(i.e packages) like controller, handlers, services
// Too much unnecessary folder structures can create circular dependencies.

func main() {
	// fmt.Println("Namaste")
	store, err := NewPostgresStore()

	if err != nil {
		log.Fatal("error connecting Postgres ", err)
	}

	if err := store.Init(); err != nil {
		log.Fatal("error creating table ", err)
	}
	srv := NewApiServer(":5000", store)
	srv.Run()
}
