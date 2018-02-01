package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/terraform-ibm-provider-api/utils"
)

func main() {
	var port int
	flag.IntVar(&port, "p", 9080, "Port on which this server listens")
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/configuration", utils.ConfHandler)
	r.HandleFunc("/configuration/{repo_name}/plan", utils.PlanHandler)
	r.HandleFunc("/configuration/{repo_name}/apply", utils.ApplyHandler)
	r.HandleFunc("/configuration/{repo_name}/destroy", utils.DestroyHandler)
	r.HandleFunc("/configuration/{repo_name}/{action}/log/{logID}", utils.LogHandler)
	fmt.Println("Server will listen at port", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)
	if err != nil {
		fmt.Printf("Couldn't start the server %v", err)
	}
}
