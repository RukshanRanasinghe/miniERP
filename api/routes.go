package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"examples.com/grocery/config"
	_ "github.com/go-sql-driver/mysql"
)

var (
	TplSuppliers *template.Template
	// TplSupplierData *template.Template
	// TplAltSup       *template.Template
	TplEditCustomer *template.Template
	TplItems        *template.Template
	// TplAddNewItem   *template.Template
	// TplEditItem    *template.Template
	TplTransaction *template.Template
	// TplProcure         *template.Template
	// TplInvoice *template.Template
	// TplPrintInvoice    *template.Template
	// TplConfirmedUpdate *template.Template
	TplInvoiceReport *template.Template
	TplReports       *template.Template
)

func init() {
	TplSuppliers = template.Must(template.ParseGlob("html/suppliers/*"))
	// TplSupplierData = template.Must(template.ParseGlob("html/suppliers/editSuppliers.html"))
	// TplAltSup = template.Must(template.ParseGlob("html/suppliers/alterSupplierData.html"))
	TplEditCustomer = template.Must(template.ParseFiles("html/Customers/editCustomers.html"))
	TplItems = template.Must(template.ParseGlob("html/Items/*"))
	// TplAddNewItem = template.Must(template.ParseGlob("html/Items/addItems.html"))
	// TplEditItem = template.Must(template.ParseGlob("html/Items/*"))
	TplTransaction = template.Must(template.ParseGlob("html/transactions/*"))
	// TplProcure = template.Must(template.ParseGlob("html/transactions/procure.html"))
	// TplInvoice = template.Must(template.ParseGlob("html/transactions/invoicing.html"))
	// TplPrintInvoice = template.Must(template.ParseGlob("html/transactions/printInvoice.html"))
	// TplConfirmedUpdate = template.Must(template.ParseFiles("html/confirmUpdateToDB.html"))
	TplReports = template.Must(template.ParseFiles("html/reports/viewStatus.html"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/index.html")
}

func AddSupplier(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/suppliers/addSuppliers.html")
}

func ConfirmSuppler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "A supplier record was saved...................")
	if config.PostError(w, r) != "nil" {
		log.Panic("Err Post")
	}
	config.SQLHandleFunc(r, "INSERT")
}

func EditSupplier(w http.ResponseWriter) {
	supplierInfo := config.GetSupplierData()
	TplSuppliers.ExecuteTemplate(w, "editSuppliers.html", supplierInfo) // Send the data to the html template
}

func AlterSupplier(w http.ResponseWriter, _ *http.Request) {
	supplierInfo := config.GetSupplierData()
	supData, err := json.Marshal(supplierInfo)
	config.Report(err, "Could not convert to JSON in function alterSupplier")
	TplSuppliers.ExecuteTemplate(w, "alterSupplierData.html", string(supData))
}

func AlterConfirm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "The supplier data was updated...................")
	config.SQLHandleFunc(r, "UPDATE")
}

func AddNewCustomerGroup(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/Customers/addCustomers.html")
}

func ConfirmCustomerGroup(w http.ResponseWriter, r *http.Request) {
	if priceLevel, err := strconv.ParseFloat(r.FormValue("PriceLevel"), 32); err == nil && r.FormValue("Category") != "" {
		insertDB, err := config.DB.Query("INSERT INTO `grocery`.`customer` (`category`,`priceLevel`, `status`)VALUES (?,?,?);", r.FormValue("Category"), priceLevel, "active")
		config.Report(err, "Could not write to db, insert connection")
		defer insertDB.Close()
		fmt.Fprintln(w, "New customer group was created.................")
	} else {
		fmt.Fprintln(w, "The data was not saved. Please check the data and resubmit")
	}
}

func EditCustomerGroup(w http.ResponseWriter, _ *http.Request) {
	customerInf := config.GetCustomerData()
	customerRecs, err := json.Marshal(customerInf)
	config.Report(err, "Could not convert to JSON in function EditCustomerGroup")
	TplEditCustomer.ParseFiles("/")
	TplEditCustomer.Execute(w, string(customerRecs))
}

func ConfirmEditCustomerGroup(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("CstGroup") != "NA" {
		config.UpdateCustomerData(w, r)
	} else {
		fmt.Fprintln(w, "Could not save data")
	}
}

func AddNewItem(w http.ResponseWriter, _ *http.Request) {
	supplierInfo := config.GetSupplierData()
	supData, err := json.Marshal(supplierInfo)
	config.Report(err, "Could not convert to JSON in function alterSupplier")
	TplItems.ExecuteTemplate(w, "addItems.html", string(supData))
}

func ConfirmAddItem(w http.ResponseWriter, r *http.Request) {
	err := config.SQLAddNewItem(r)
	if err == nil {
		fmt.Fprintln(w, "Item was saved")
	}
}

func EditItem(w http.ResponseWriter, _ *http.Request) {
	items := config.SQLFetchItem()
	TplItems.ExecuteTemplate(w, "editItem.html", string(items))
}

func ConfirmEditItem(w http.ResponseWriter, r *http.Request) {
	config.SQLUpdateItem(r)
	fmt.Fprintln(w, "The item was updated")
}

func ProcureItems(w http.ResponseWriter, _ *http.Request) {
	items := config.SQLFetchItem()
	TplTransaction.ExecuteTemplate(w, "procure.html", string(items))
	// TplProcure.Execute(w, string(items))
}

func AddToInventory(w http.ResponseWriter, r *http.Request) {
	config.SQLUpdateInventory(r)
	fmt.Fprintln(w, "Successfully add item to inventory")
}

func Invoice(w http.ResponseWriter, _ *http.Request) {
	type invoiceData struct {
		Stocks    string `json:"stocks"`
		Customers string `json:"customer"`
	}

	invStocks, invCustomers := config.SQLFetchStocks()
	detailsForInvoice := invoiceData{string(invStocks), string(invCustomers)}
	TplTransaction.ExecuteTemplate(w, "invoicing.html", detailsForInvoice)
}

func UpdateInvoice(w http.ResponseWriter, r *http.Request) {
	invoiceForm := config.RecordInvoice(r)
	TplTransaction.ExecuteTemplate(w, "printInvoice.html", invoiceForm)
}

func ViewInvoice(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./html/reports/viewStatus.html")
}

func SubmitInvoiceData(w http.ResponseWriter, r *http.Request) {
	stocks, purchases, invoices := config.GetInvoiceDataFrmDB(r)
	report := struct {
		Stocks    string `json:"stocks"`
		Purchases string `json:"purchases"`
		Invoices  string `json:"invoice"`
	}{
		string(stocks),
		string(purchases),
		string(invoices)}

	TplReports.Execute(w, report)
}
