package config

import (
	"fmt"
	"log"
	"net/http"
)

type person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Sex  string `json:"sex"`
}

var worker = person{
	Name: "Juliet Jackson",
	Age:  35,
	Sex:  "Female",
}

func Report(err error, message string) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func PostError(w http.ResponseWriter, r *http.Request) string {
	if r.Method != "POST" {
		http.Redirect(w, r, "./grocery", http.StatusSeeOther)
		return "Error posting the record"
	}
	return "nil"
}

func StringsJoin[t interface{}](items []t) string {
	var itemList string = ""
	for _, item := range items {
		itemList = fmt.Sprintf("%v,%v", item, itemList)
	}
	return itemList
}
func ConfirmUpdate(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/confirmUpdateToDB.html")
	http.ServeFile(w, r, "./html/html/transactions/printInvoice.html")
}
