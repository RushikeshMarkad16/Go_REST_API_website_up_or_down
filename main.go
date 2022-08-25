package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var link_map = map[string]string{}

func getLinks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get link using post endpoint hit")

	for key, value := range link_map {
		fmt.Printf("[%s] = %s\n", key, value)
	}

	fmt.Fprint(w, link_map)

}

func PostLinks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("post link endpoint hit")
	post_map := map[string][]string{}
	err := json.NewDecoder(r.Body).Decode(&post_map)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(post_map)

	for _, link := range post_map["websites"] {
		link_map[link] = "Down"
	}

	fmt.Fprint(w, "Successfully Added Website Links")
}

func getLinks_by_id(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	a := "https://" + params["link"]
	_, err := http.Get(a)
	if err != nil {
		fmt.Println(params["link"], " : Down")
		fmt.Fprint(w, params["link"], " : Down")
	} else {
		fmt.Println(params["link"], " : Up")
		fmt.Fprint(w, params["link"], " : Up")
	}
}

func main() {

	go checkStatus()
	r := mux.NewRouter()
	r.HandleFunc("/post-links", PostLinks)
	r.HandleFunc("/get-links", getLinks)
	r.HandleFunc("/get-links/{link}", getLinks_by_id)
	fmt.Println("Server running on... localhost:8081")
	http.ListenAndServe(":8081", r)
}

func checkStatus() {
	for {
		for key := range link_map {
			resp, err := http.Get(key)
			if err != nil {
				fmt.Println("Error Occured")
				link_map[key] = "Down"
				continue
			}

			if resp.StatusCode == 200 {
				fmt.Println("Successful")
				link_map[key] = "Up"
			}

		}
	}

}
