package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func InitiateServer() {
	fmt.Println("Initiating the server........")
	route := mux.NewRouter()

	fs := http.FileServer(http.Dir("html/assets/"))
	route.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	route.HandleFunc("/grocery", Index).Methods("GET")

	route.HandleFunc("/addsuppliers", AddSupplier).Methods("GET")
	route.HandleFunc("/confirmsuppler", ConfirmSuppler).Methods("POST")

	route.HandleFunc("/altersupplierdetails", AlterSupplier).Methods(("GET"))
	route.HandleFunc("/confirmAlter", AlterConfirm).Methods(("POST"))

	route.HandleFunc("/addNewCustomerGroup", AddNewCustomerGroup).Methods("GET")
	route.HandleFunc("/confirmCustomerGroup", ConfirmCustomerGroup).Methods("POST")

	route.HandleFunc("/editCustomerGroup", EditCustomerGroup).Methods("GET")
	route.HandleFunc("/confirmEditCustomerGroup", ConfirmEditCustomerGroup).Methods("POST")

	route.HandleFunc("/addNewItem", AddNewItem).Methods("GET")
	route.HandleFunc("/confirmAddItem", ConfirmAddItem).Methods("POST")

	route.HandleFunc("/editItem", EditItem).Methods("GET")
	route.HandleFunc("/confirmEditItem", ConfirmEditItem).Methods("POST")

	route.HandleFunc("/procureItems", ProcureItems).Methods("GET")
	route.HandleFunc("/addToInventory", AddToInventory).Methods("POST")

	route.HandleFunc("/invoice", Invoice).Methods("GET")
	route.HandleFunc("/updateInvoice", UpdateInvoice).Methods("POST")

	route.HandleFunc("/viewInvoice", ViewInvoice).Methods("GET")
	route.HandleFunc("/submitInvoiceData", SubmitInvoiceData).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", route))
}
