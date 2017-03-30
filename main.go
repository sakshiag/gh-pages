package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.ibm.com/terraform-devops-tools/e2erunner/utils"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 9080, "Port on which this server listens")
	flag.Parse()
	mux := http.NewServeMux()
	mux.HandleFunc("/e2e", utils.E2EHandler)
	fmt.Println("Server will listen at port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		fmt.Printf("Couldn't start the server %v", err)
	}
}
