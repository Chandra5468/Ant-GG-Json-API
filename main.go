// This is a bank application

package main

// Not creating much folder structures(i.e packages) like controller, handlers, services
// Too much unnecessary folder structures can create circular dependencies.

func main() {
	// fmt.Println("Namaste")
	srv := NewApiServer(":5000")
	srv.Run()
}
