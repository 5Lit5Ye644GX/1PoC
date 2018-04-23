package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	//"github.com/gorilla/mux"
)

var transactions []Transaction
var balance []Balance

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TodoIndex(w http.ResponseWriter, r *http.Request) {

}

func TodoShow(w http.ResponseWriter, r *http.Request) {

}

func PostTodo(w http.ResponseWriter, r *http.Request) {

}

func getAllTransaction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(transactions); err != nil {
		panic(err)
	}
}

////////////
func getBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(balance); err != nil {
		panic(err)
	}
}

func postBalance(w http.ResponseWriter, r *http.Request) {
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Print(readErr)
	}
	data := Balance{}
	// json.NewDecoder(string(body)).Decode(transaction)
	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Print(jsonErr)
	}
	fmt.Println(data.Addresses)
	if data.Addresses != "" {
		fmt.Println("Receive Balances")
		balance = append(balance, data)
	} else {
		fmt.Println("Balances Fails")
	}

}

//////
func postTransaction(w http.ResponseWriter, r *http.Request) {
	body, readErr := ioutil.ReadAll(r.Body)
	if readErr != nil {
		log.Print(readErr)
	}
	data := Transaction{}
	// json.NewDecoder(string(body)).Decode(transaction)
	jsonErr := json.Unmarshal(body, &data)
	if jsonErr != nil {
		log.Print(jsonErr)
	}
	fmt.Println(data.Result)
	if data.Result != nil {
		fmt.Println("Receive transaction")
		transactions = append(transactions, data)
	} else {
		fmt.Println("Transaction Fails")
	}

}
