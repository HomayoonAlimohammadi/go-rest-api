package api

import (
	"backend/helpers"
	"backend/users"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Login struct {
	Username string
	Password string
}

type ErrResponse struct {
	Message string
}

func login(w http.ResponseWriter, r *http.Request) {

	// Ready body
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleError(err)

	var formattedBody Login
	err = json.Unmarshal(body, &formattedBody)
	helpers.HandleError(err)
	login := users.Login(formattedBody.Username, formattedBody.Password)

	if login["message"] == "All is fine" {
		resp := login
		json.NewEncoder(w).Encode(resp)
	} else {
		// Handle Error
		resp := ErrResponse{Message: "Wrong Username or Password"}
		json.NewEncoder(w).Encode(resp)
	}
}

func StartApp() {
	router := mux.NewRouter()
	router.HandleFunc("/login", login).Methods("POST")
	fmt.Println("App is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
