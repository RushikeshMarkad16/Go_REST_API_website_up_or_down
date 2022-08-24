package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

var link_map = map[string]string{}
var post_map map[string][]string

func getLinks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get link using post endpoint hit")

	for _, link := range post_map["websites"] {
		link_map[link] = "Down"
	}

	for key := range link_map {
		_, err := http.Get(key)
		if err != nil {
			link_map[key] = "Down"
		} else {
			link_map[key] = "Up"
		}
	}

	for key, value := range link_map {
		fmt.Printf("[%s] = %s\n", key, value)
	}

	fmt.Fprint(w, link_map)

}

func PostLinks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("post link endpoint hit")
	post_map = map[string][]string{}
	err := json.NewDecoder(r.Body).Decode(&post_map)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(post_map)
	err = json.NewEncoder(w).Encode(post_map)
	if err != nil {
		fmt.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
}

func getLinks_by_id(w http.ResponseWriter, r *http.Request) {

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/post-links", PostLinks)
	r.HandleFunc("/get-links", getLinks)
	r.HandleFunc("/get-links/{key}", getLinks_by_id)
	fmt.Println("Server running on... localhost:8081")
	http.ListenAndServe(":8081", r)
}
