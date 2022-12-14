package api

import (
	"backend/helpers"
	"backend/interfaces"
	"backend/transactions"
	"backend/useraccounts"
	"backend/users"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type TransactionBody struct {
	UserID uint
	From   uint
	To     uint
	Amount int
}

func readBody(r *http.Request) []byte {
	body, err := ioutil.ReadAll(r.Body)
	helpers.HandleError(err)

	return body
}

func apiResponse(call map[string]interface{}, w http.ResponseWriter) {
	if call["message"] == "All is fine" {
		resp := call
		json.NewEncoder(w).Encode(resp)
	} else {
		// Handle Error
		resp := call
		json.NewEncoder(w).Encode(resp)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	var formattedBody interfaces.Login
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleError(err)
	login := users.Login(formattedBody.Username, formattedBody.Password)
	apiResponse(login, w)
}

func register(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	var formattedBody interfaces.Register
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleError(err)
	register := users.Register(formattedBody.Username, formattedBody.Email, formattedBody.Password)
	apiResponse(register, w)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
	jwt := r.Header.Get("Authorization")
	user := users.GetUser(userID, jwt)

	apiResponse(user, w)
}

func getMyTransactions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	jwt := r.Header.Get("Authorization")
	transactions := transactions.GetMyTransactions(userID, jwt)

	apiResponse(transactions, w)
}

func StartApp() {
	router := mux.NewRouter()
	router.Use(helpers.PanicHandler)
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/users/{id}", getUser).Methods("GET")
	router.HandleFunc("/transaction", transaction).Methods("POST")
	router.HandleFunc("/transaction/{userID}", getMyTransactions).Methods("GET")
	fmt.Println("App is listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func transaction(w http.ResponseWriter, r *http.Request) {
	body := readBody(r)
	auth := r.Header.Get("Authorization")

	var formattedBody TransactionBody
	err := json.Unmarshal(body, &formattedBody)
	helpers.HandleError(err)
	transaction := useraccounts.Transaction(formattedBody.UserID, formattedBody.From, formattedBody.To, formattedBody.Amount, auth)
	apiResponse(transaction, w)

}
