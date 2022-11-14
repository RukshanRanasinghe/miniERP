package utility

type Supplier struct {
	Id           string `json:"id"`
	Date         string `json:"date"`
	TelNum       string `json:"telnum"`
	Salutations  string `json:"salutations"`
	FirstName    string `json:"firstname"`
	LastName     string `json:"lastname"`
	Organization string `json:"organization"`
	Address      string `json:"address"`
	Status       string `json:"status"`
}

type Customer struct {
	Id         string  `json:"id"`
	Category   string  `json:"category"`
	PriceLevel float32 `json:"priceLevel"`
	Status     string  `json:"status"`
}

type Item struct {
	Id              string  `json:"id"`
	CreatedDate     string  `json:"createdate"`
	ItemName        string  `json:"itemname"`
	ItemDescription string  `json:"itemdescription"`
	ItemCode        string  `json:"itemcode"`
	PurchasePrice   float64 `json:"purchasedprice"`
	SellingPrice    float64 `json:"sellingprice"`
	Units           string  `json:"units"`
	SupID           string  `json:"supid"`
	VATRegistered   string  `json:"vatreg"`
	Status          string  `json:"status"`
}

type Inventory struct {
	Id              string  `json:"id"`
	CreatedDate     string  `json:"createdate"`
	ItemCode        string  `json:"itemcode"`
	ItemName        string  `json:"itemname"`
	ItemDescription string  `json:"itemdescription"`
	PurchasePrice   float64 `json:"purchasedprice"`
	Quantity        int     `json:"quantity"`
	Units           string  `json:"units"`
	SellingPrice    float64 `json:"sellingprice"`
}

type Stocks struct {
	Id            string  `json:"id"`
	ItemCode      string  `json:"itemcode"`
	ItemName      string  `json:"itemname"`
	Quantity      int     `json:"quantity"`
	Units         string  `json:"units"`
	PurchasePrice float64 `json:"purchasedprice"`
	SellingPrice  float64 `json:"sellingprice"`
}

type Invoice struct {
	Id             string    `json:"id"`
	InvoiceDate    string    `json:"InvoiceDate"`
	InvoiceNumber  string    `json:"InvoiceNumber"`
	ItemCodes      []string  `json:"ItemCodes"`
	ItemNames      []string  `json:"ItemNames"`
	Quantities     []int     `json:"Quantities"`
	ItemUnits      []string  `json:"ItemUnits"`
	ItemSubTotal   []float64 `json:"ItemSubTotal"`
	TotalBillValue float64   `json:"TotalBillValue"`
}

type Bill struct {
	Id             string `json:"id"`
	InvoiceDate    string `json:"InvoiceDate"`
	InvoiceNumber  string `json:"InvoiceNumber"`
	ItemCodes      string `json:"ItemCodes"`
	ItemNames      string `json:"ItemNames"`
	Quantities     string `json:"Quantities"`
	ItemUnits      string `json:"ItemUnits"`
	ItemSubTotal   string `json:"ItemSubTotal"`
	TotalBillValue string `json:"TotalBillValue"`
}
