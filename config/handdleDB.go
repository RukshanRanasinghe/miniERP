package config

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"examples.com/grocery/utility"
	_ "github.com/go-sql-driver/mysql"
)

var (
	DB  *sql.DB
	err error
)

func Connection() {
	DB, err = sql.Open("mysql", "root:Patch12-34@tcp(127.0.0.1:3306)/grocery")
	Report(err, "Could not open the DB")
}

func GetCustomerData() []utility.Customer {
	var (
		customer  utility.Customer
		customers []utility.Customer
	)
	fetchCustomersFromDB, err := DB.Query("SELECT * FROM `grocery`.`customer`")
	Report(err, "Could not get data from customer DB")

	for fetchCustomersFromDB.Next() {
		err = fetchCustomersFromDB.Scan(
			&customer.Id,
			&customer.Category,
			&customer.PriceLevel,
			&customer.Status,
		)
		// record := utility.Customer{
		// 	Id:         customer.Id,
		// 	Category:   customer.Category,
		// 	PriceLevel: customer.PriceLevel,
		// 	Status:     customer.Status,
		// }
		customers = append(customers, customer)
		Report(err, "Could not fetch data from the Customer DB")
	}
	return customers
}

func UpdateCustomerData(w http.ResponseWriter, r *http.Request) {
	dbQuarry := fmt.Sprintf("UPDATE `grocery`.`customer` SET `category`=?,`priceLevel`=?, `status`=? WHERE (`id` = %s);", r.FormValue("CstGroup"))

	insertDB, err := DB.Query(dbQuarry,
		r.FormValue("customerGroup"),
		r.FormValue("PriceLevel"),
		r.FormValue("Status"),
	)
	Report(err, "Could not write to db, insert connection")
	defer insertDB.Close()
	fmt.Fprintf(w, "Form values -> %s, %s, %s, %s\n", r.FormValue("CstGroup"), r.FormValue("customerGroup"), r.FormValue("PriceLevel"), r.FormValue("Status"))
}

func GetSupplierData() []utility.Supplier {
	var (
		supplier     utility.Supplier
		supplierData []utility.Supplier
	)
	results, err := DB.Query("SELECT * FROM `grocery`.`suppliers`")
	Report(err, "Could not open db to read data from table")

	for results.Next() {
		err = results.Scan(&supplier.Id,
			&supplier.Date,
			&supplier.TelNum,
			&supplier.Salutations,
			&supplier.FirstName,
			&supplier.LastName,
			&supplier.Organization,
			&supplier.Address,
			&supplier.Status)
		// data := utility.Supplier{
		// 	Id:           supplier.Id,
		// 	Date:         supplier.Date,
		// 	TelNum:       supplier.TelNum,
		// 	Salutations:  supplier.Salutations,
		// 	FirstName:    supplier.FirstName,
		// 	LastName:     supplier.LastName,
		// 	Organization: supplier.Organization,
		// 	Address:      supplier.Address,
		// 	Status:       supplier.Status,
		// }
		supplierData = append(supplierData, supplier)
		Report(err, "Could not get data from DB")
	}
	return supplierData
}

func SQLFetchItem() []byte {
	var (
		item     utility.Item
		items    []utility.Item
		supplier utility.Supplier
	)
	suppOrganization := make(map[string]string)

	fetchSupplierFromDB, _ := DB.Query("SELECT `ID`,`organization` FROM `grocery`.`suppliers`;")

	for fetchSupplierFromDB.Next() {
		fetchSupplierFromDB.Scan(&supplier.Id, &supplier.Organization)
		suppOrganization[supplier.Id] = supplier.Organization
	}
	fetchSupplierFromDB.Close()

	fetchItemFrmDB, err := DB.Query("SELECT * FROM `grocery`.`items`;")
	Report(err, "Could not open item db")

	for fetchItemFrmDB.Next() {
		err = fetchItemFrmDB.Scan(&item.Id,
			&item.CreatedDate,
			&item.ItemName,
			&item.ItemDescription,
			&item.ItemCode,
			&item.PurchasePrice,
			&item.SellingPrice,
			&item.Units,
			&item.SupID,
			&item.VATRegistered,
			&item.Status)
		Report(err, "Could not read item from db")
		record := utility.Item{
			Id:              item.Id,
			CreatedDate:     item.CreatedDate,
			ItemName:        item.ItemName,
			ItemDescription: item.ItemDescription,
			ItemCode:        item.ItemCode,
			PurchasePrice:   item.PurchasePrice,
			SellingPrice:    item.SellingPrice,
			Units:           item.Units,
			SupID:           suppOrganization[item.SupID],
			VATRegistered:   item.VATRegistered,
			Status:          item.Status,
		}
		items = append(items, record)
	}
	itemJSON, err := json.Marshal(items)
	Report(err, "Could not convert item rec to JSON")
	return itemJSON
}

func SQLHandleFunc(r *http.Request, method string) {
	var dbQuarry, status string

	switch method {
	case "UPDATE":
		dbQuarry = fmt.Sprintf("UPDATE `grocery`.`suppliers` SET `date`=?,`tel_num`=?,`salutation`=?,`first_name`=?,`last_name`=?,`organization`=?,`address`=?,`status`=? WHERE (`ID` = %s);", r.FormValue("id"))
		status = r.FormValue("Status")
	case "INSERT":
		dbQuarry = "INSERT INTO `grocery`.`suppliers` (`date`,`tel_num`,`salutation`,`first_name`,`last_name`,`organization`,`address`,`status`)VALUES (?,?,?,?,?,?,?,?);"
		status = "active"
	}
	insertDB, err := DB.Query(dbQuarry,
		r.FormValue("supplierCrtDate"),
		r.FormValue("telNumber"),
		r.FormValue("selection"),
		r.FormValue("firstName"),
		r.FormValue("lastName"),
		r.FormValue("organization"),
		r.FormValue("longAddress"),
		status)

	Report(err, "Could not write to db, insert connection")
	defer insertDB.Close()
}

func SQLAddNewItem(r *http.Request) error {
	toDate := time.Now()
	dateNow := toDate.Format("2006-01-02")

	insertDB, err := DB.Query("INSERT INTO `grocery`.`items` (`date`, `item_name`,`item_description`,`item_code`,`purchase_price`,`selling_price`,`units`,`id_supplier`,`VATRegistered`,`status`)VALUES (?,?,?,?,?,?,?,?,?,?);",
		dateNow,
		r.FormValue("ItemName"),
		r.FormValue("ItemDescription"),
		r.FormValue("ItemCode"),
		r.FormValue("PurchasePrice"),
		r.FormValue("SellingPrice"),
		r.FormValue("Units"),
		r.FormValue("supID"),
		r.FormValue("VATRegistered"),
		r.FormValue("Status"))

	Report(err, "Could nor write to DB")
	defer insertDB.Close()
	return err
}

func SQLUpdateItem(r *http.Request) {
	var (
		itemEntry  utility.Item
		supplierID int
	)

	toDate := time.Now()
	dateNow := toDate.Format("2006-01-02")

	supplier := GetSupplierData()
	for _, item := range supplier {
		if r.FormValue("organization") == item.Organization {
			supplierID, _ = strconv.Atoi(item.Id)
		}
	}
	dbQuerySelect := fmt.Sprintf("SELECT * FROM `grocery`.`items` WHERE (`id` = '%s');", r.FormValue("id"))

	itemValues, _ := DB.Query(dbQuerySelect)
	for itemValues.Next() {
		err = itemValues.Scan(&itemEntry.Id,
			&itemEntry.CreatedDate,
			&itemEntry.ItemName,
			&itemEntry.ItemDescription,
			&itemEntry.ItemCode,
			&itemEntry.PurchasePrice,
			&itemEntry.SellingPrice,
			&itemEntry.Units,
			&itemEntry.SupID,
			&itemEntry.VATRegistered,
			&itemEntry.Status)
	}

	dbQueryChangeStatus := fmt.Sprintf("UPDATE `grocery`.`items` SET `status`='inactive' WHERE (`id` = %s);", r.FormValue("id"))
	_, err = DB.Query(dbQueryChangeStatus)
	Report(err, "Could not update status in itemDB")

	if r.FormValue("inactivate") != "true" {
		dbQueryInsertAmended := ("INSERT INTO `grocery`.`items` (`date`, `item_name`,`item_description`,`item_code`,`purchase_price`,`selling_price`,`units`,`id_supplier`,`VATRegistered`,`status`)VALUES (?,?,?,?,?,?,?,?,?,?);")
		_, err = DB.Query(dbQueryInsertAmended,
			dateNow,
			itemEntry.ItemName,
			itemEntry.ItemDescription,
			r.FormValue("ItemCode"),
			r.FormValue("PurchasePrice"),
			r.FormValue("SellingPrice"),
			r.FormValue("Units"),
			supplierID,
			r.FormValue("VATRegistered"),
			"active")

		Report(err, "Could not update item")

		selPrice, _ := strconv.ParseFloat(r.FormValue("SellingPrice"), 64)
		purPrice, _ := strconv.ParseFloat(r.FormValue("PurchasePrice"), 64)
		updateStocks := fmt.Sprintf("UPDATE `grocery`.`stocks` SET `purchase_price`=%f, `selling_price`=%f WHERE (`item_code` = '%s');", purPrice, selPrice, r.FormValue("ItemCode"))
		_, err = DB.Query(updateStocks)
		Report(err, "Could not update stocks")
	}
}

func SQLUpdateInventory(r *http.Request) {
	var (
		stock         utility.Stocks
		itemCodeFound = false
		query         string
		err           error
	)
	queryInventoryDB, err := DB.Query("INSERT INTO `grocery`.`purchases` (`purchase_date`,`item_code`,`item_name`,`item_description`,`purchase_price`,`quantity`,`units`,`selling_price`)VALUES (?,?,?,?,?,?,?,?);",
		r.FormValue("purchaseDate"),
		r.FormValue("itemCode"),
		r.FormValue("itemName"),
		r.FormValue("description"),
		r.FormValue("value"),
		r.FormValue("purchaseQuantity"),
		r.FormValue("units"),
		r.FormValue("sellingPrice"))
	Report(err, "Could not update the inventory")

	quantity, _ := strconv.Atoi(r.FormValue("purchaseQuantity"))
	purchaseCost, _ := strconv.ParseFloat(r.FormValue("value"), 64)
	sellingPrize, _ := strconv.ParseFloat(r.FormValue("sellingPrice"), 64)
	queryStocks, err := DB.Query("SELECT `item_code`, `quantity` FROM `grocery`.`stocks`;")
	Report(err, "Could not update the inventory")
	for queryStocks.Next() {
		queryStocks.Scan(&stock.ItemCode, &stock.Quantity)
		if stock.ItemCode == r.FormValue("itemCode") {
			stock.Quantity = stock.Quantity + quantity
			query = fmt.Sprintf("UPDATE `grocery`.`stocks` SET `quantity`=%d, `purchase_price`=%f, `selling_price`=%f WHERE (`item_code`='%s');", stock.Quantity, purchaseCost, sellingPrize, stock.ItemCode)
			itemCodeFound = true
		}
	}
	if itemCodeFound {
		queryStocks, err = DB.Query(query)
	} else {
		queryStocks, err = DB.Query("INSERT INTO `grocery`.`stocks` (`item_code`,`item_name`,`quantity`,`units`, `purchase_price`, `selling_price`)VALUES (?,?,?,?,?,?);",
			r.FormValue("itemCode"),
			r.FormValue("itemName"),
			quantity,
			r.FormValue("units"),
			purchaseCost,
			sellingPrize)
	}
	Report(err, "Could not update the inventory")
	defer queryInventoryDB.Close()
	defer queryStocks.Close()
}

func SQLFetchStocks() ([]byte, []byte) {
	var (
		stockItem       utility.Stocks
		stocks          []utility.Stocks
		customerData    utility.Customer
		customerDetails []utility.Customer
	)

	stockRecords, err := DB.Query("SELECT * FROM `grocery`.`stocks`;")
	Report(err, "Could not open stocks DB")

	for stockRecords.Next() {
		err = stockRecords.Scan(&stockItem.Id,
			&stockItem.ItemCode,
			&stockItem.ItemName,
			&stockItem.Quantity,
			&stockItem.Units,
			&stockItem.PurchasePrice,
			&stockItem.SellingPrice)
		// record := utility.Stocks{
		// 	Id:            stockItem.Id,
		// 	ItemCode:      stockItem.ItemCode,
		// 	ItemName:      stockItem.ItemName,
		// 	Quantity:      stockItem.Quantity,
		// 	Units:         stockItem.Units,
		// 	PurchasePrice: stockItem.PurchasePrice,
		// 	SellingPrice:  stockItem.SellingPrice}
		stocks = append(stocks, stockItem)
	}

	Report(err, "Could not open stock DB")
	customerRecords, err := DB.Query("SELECT * FROM `grocery`.`customer`;")
	for customerRecords.Next() {
		customerRecords.Scan(&customerData.Id,
			&customerData.Category,
			&customerData.PriceLevel,
			&customerData.Status)
		// record := utility.Customer{
		// 	Id:         customerData.Id,
		// 	Category:   customerData.Category,
		// 	PriceLevel: customerData.PriceLevel,
		// 	Status:     customerData.Status,
		// }
		customerDetails = append(customerDetails, customerData)
	}

	Report(err, "Could not open customer DB inside invoice")
	defer customerRecords.Close()

	customerJSON, _ := json.Marshal(customerDetails)
	stocksJSON, _ := json.Marshal(stocks)
	return stocksJSON, customerJSON
}

func RecordInvoice(r *http.Request) string {
	var (
		stocksBalance []utility.Stocks
		invoice       utility.Invoice
		invoiceRec    *sql.Rows
	)

	stocksForm := r.FormValue("balanceStock")
	json.Unmarshal([]byte(stocksForm), &stocksBalance)

	invoiceForm := r.FormValue("invoice")
	json.Unmarshal([]byte(invoiceForm), &invoice)

	for _, item := range stocksBalance {
		query := fmt.Sprintf("UPDATE `grocery`.`stocks` SET `quantity`=%d WHERE (`item_code`='%s');", item.Quantity, item.ItemCode)
		_, err = DB.Query(query)

		Report(err, "Could not update stocks")
	}

	invoiceQuery := fmt.Sprintln("INSERT INTO `grocery`.`invoices` (`invoice_date`, `invoice_no`, `item_code`, `items`, `quantity`, `units`, `sub_total`, `total_bill_value`) VALUES (?,?,?,?,?,?,?,?);")
	invoiceRec, err = DB.Query(invoiceQuery,
		invoice.InvoiceDate,
		invoice.InvoiceNumber,
		StringsJoin(invoice.ItemCodes),
		StringsJoin(invoice.ItemNames),
		StringsJoin(invoice.Quantities),
		StringsJoin(invoice.ItemUnits),
		StringsJoin(invoice.ItemSubTotal),
		invoice.TotalBillValue)

	Report(err, "Could not update the invoice DB")
	fmt.Printf("%v\n%v\n%v\n%v\n%v\n", StringsJoin(invoice.ItemCodes),
		StringsJoin(invoice.ItemNames),
		StringsJoin(invoice.Quantities),
		StringsJoin(invoice.ItemUnits),
		StringsJoin(invoice.ItemSubTotal))
	defer invoiceRec.Close()
	return invoiceForm
}

func GetInvoiceDataFrmDB(r *http.Request) ([]byte, []byte, []byte) {
	type groupedDetails struct {
		Quantity int     `json:"quantity"`
		Units    string  `json:"units"`
		Value    float64 `json:"purchased_price"`
		ItemName string  `json:"item_name"`
	}

	var (
		queryStocks, queryPurchases, queryInvoices *sql.Rows
		stocks                                     utility.Stocks
		inventory                                  []utility.Stocks
		purchase                                   utility.Inventory
		invoice                                    utility.Bill
		separateInvoices                           utility.Invoice
		item                                       groupedDetails
		items                                      = make(map[string]groupedDetails)
		invoiceData                                = make(map[string]groupedDetails)
		ok                                         bool
	)
	dateLayout := "2006-01-02"

	startDate, _ := time.Parse(dateLayout, r.FormValue("stateDate"))
	endDate, _ := time.Parse(dateLayout, r.FormValue("endDate"))

	queryStocks, err = DB.Query("SELECT * FROM `grocery`.`stocks`;")
	queryPurchases, err = DB.Query("SELECT * FROM `grocery`.`purchases`;")
	queryInvoices, err = DB.Query("SELECT * FROM `grocery`.`invoices`;")

	defer queryStocks.Close()
	defer queryPurchases.Close()
	defer queryInvoices.Close()

	for queryStocks.Next() {
		queryStocks.Scan(&stocks.Id, &stocks.ItemCode, &stocks.ItemName, &stocks.Quantity, &stocks.Units, &stocks.PurchasePrice, &stocks.SellingPrice)
		inventory = append(inventory, stocks)
	}

	for queryPurchases.Next() {
		queryPurchases.Scan(&purchase.Id, &purchase.CreatedDate, &purchase.ItemCode, &purchase.ItemName, &purchase.ItemDescription, &purchase.PurchasePrice, &purchase.Quantity, &purchase.Units, &purchase.SellingPrice)
		dateOfRecord, _ := time.Parse(dateLayout, purchase.CreatedDate)

		if dateOfRecord.After(startDate) && dateOfRecord.Before(endDate) {
			if item, ok = items[purchase.ItemCode]; ok {
				item.Quantity += purchase.Quantity
				items[purchase.ItemCode] = item

			} else {
				item.Quantity = purchase.Quantity
				item.Units = purchase.Units
				item.ItemName = purchase.ItemName
				item.Value = purchase.PurchasePrice
				items[purchase.ItemCode] = item
			}
		}
	}

	for queryInvoices.Next() {
		queryInvoices.Scan(&invoice.Id, &separateInvoices.InvoiceDate, &separateInvoices.InvoiceNumber, &invoice.ItemCodes, &invoice.ItemNames, &invoice.Quantities, &invoice.ItemUnits, &invoice.ItemSubTotal, &separateInvoices.TotalBillValue)
		separateInvoices.ItemCodes = strings.Split(invoice.ItemCodes, ",")
		separateInvoices.ItemNames = strings.Split(invoice.ItemNames, ",")
		separateInvoices.ItemUnits = strings.Split(invoice.ItemUnits, ",")
		separateInvoices.Quantities = nil
		separateInvoices.ItemSubTotal = nil

		for _, qty := range strings.Split(invoice.Quantities, ",") {
			value, _ := strconv.Atoi(qty)
			separateInvoices.Quantities = append(separateInvoices.Quantities, value)
		}
		for _, STol := range strings.Split(invoice.ItemSubTotal, ",") {
			value, _ := strconv.ParseFloat(STol, 64)
			separateInvoices.ItemSubTotal = append(separateInvoices.ItemSubTotal, value)
		}

		for index, code := range separateInvoices.ItemCodes {
			item = invoiceData[code]
			if _, ok := invoiceData[code]; ok {
				item.Quantity += separateInvoices.Quantities[index]
				item.Value += separateInvoices.ItemSubTotal[index]
				invoiceData[code] = item
			} else {
				item.ItemName = separateInvoices.ItemNames[index]
				item.Quantity = separateInvoices.Quantities[index]
				item.Units = separateInvoices.ItemUnits[index]
				item.Value = separateInvoices.ItemSubTotal[index]
				invoiceData[code] = item
			}
		}

		delete(invoiceData, "")
	}

	Report(err, "Could not open DB for reporting")
	inventoryJSON, _ := json.Marshal(inventory)
	purchasesJSON, _ := json.Marshal(items)
	invoicesJSON, _ := json.Marshal(invoiceData)
	return inventoryJSON, purchasesJSON, invoicesJSON
}
