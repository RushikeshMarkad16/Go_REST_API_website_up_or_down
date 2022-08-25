package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Error struct {
	ErrorMessage string `json:"error"`
}

var link_map = map[string]string{}

func getLinks(w http.ResponseWriter, r *http.Request) {
	fmt.Println("get link using post endpoint hit")

	for key, value := range link_map {
		fmt.Printf("[%s] = %s\n", key, value)
	}
	// fmt.Fprint(w, link_map)
	jsonStr, err := json.Marshal(link_map)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Println(string(jsonStr))
		fmt.Fprint(w, string(jsonStr))
	}
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

	_, isPresent := link_map[a]
	if !isPresent {
		errMsg := Error{
			ErrorMessage: "Webiste not presnt in the database",
		}
		json.NewEncoder(w).Encode(errMsg)
		return
	}

	responseMap := make(map[string]string)
	responseMap[a] = link_map[a]
	json.NewEncoder(w).Encode(responseMap)

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
		time.Sleep(1 * time.Minute)
	}
}
