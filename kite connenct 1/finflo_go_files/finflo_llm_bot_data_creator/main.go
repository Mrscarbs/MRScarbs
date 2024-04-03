package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const apiKey = "U6TZY67PIO8WLRWB"

func main() {

	symbol := "AAPL" // Symbol of the stock you want to fetch data for
	url := fmt.Sprintf("https://www.alphavantage.co/query?function=OVERVIEW&symbol=%s&apikey=%s", symbol, apiKey)

	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
	}

	databytes, err2 := io.ReadAll(response.Body)

	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(string(databytes))

	var stock_datamap = make(map[string]interface{})

	err3 := json.Unmarshal(databytes, &stock_datamap)

	if err3 != nil {
		fmt.Println(err3)
	}

	pb := stock_datamap["PriceToBookRatio"]
	yearHigh := stock_datamap["52WeekHigh"]
	yearlow := stock_datamap["52WeekLow"]
	EVToEBITDA := stock_datamap["EVToEBITDA"]
	EVToRevenue := stock_datamap["EVToRevenue"]
	EPS := stock_datamap["EPS"]
	PERatio := stock_datamap["PERatio"]
	DividendYield := stock_datamap["DividendYield"]
	Sector := stock_datamap["Sector"]
	Industry := stock_datamap["Industry"]
	MarketCapitalization := stock_datamap["MarketCapitalization"]
	BookValue := stock_datamap["BookValue"]
	ProfitMargin := stock_datamap["ProfitMargin"]

	db, err_db := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")

	if err_db != nil {
		fmt.Println(err_db)
	}
	time_st := time.Now()

	unix_time := time_st.Unix()

	db.Exec("call stp_InsertStockData(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", pb, yearHigh, yearlow, EVToEBITDA, EVToRevenue, EPS, PERatio, DividendYield, Sector, Industry, MarketCapitalization, BookValue, ProfitMargin, unix_time, symbol)
	defer db.Close()
	defer response.Body.Close()
	llm_content := fmt.Sprintf("The pb or price to bookvalue for %s is %s, the year high for %s is %s, the year low for %s is%s,the EVToEBITDA for %s is %s,the EVToRevenue for %s is %s,the EPS for %s is %s,the pe or pe ratio or price to earning for %s is %s,the DividendYield for %s is %s,the sector for %s is %s,the industry for %s is %s,the marketcap or market capitalization for %s is %s, the bookvalue for %s is %s, the profitmargin for %s is %s", symbol, pb, symbol, yearHigh, symbol, yearlow, symbol, EVToEBITDA, symbol, EVToRevenue, symbol, EPS, symbol, PERatio, symbol, DividendYield, symbol, Sector, symbol, Industry, symbol, MarketCapitalization, symbol, BookValue, symbol, ProfitMargin)
	fmt.Println(llm_content)
	file_creator(llm_content)

	data, err := os.ReadFile("llm_doc.txt")
	fmt.Println(string(data))
}
func file_creator(content string) {
	file, err := os.Create("llm_doc.txt")
	if err != nil {
		fmt.Println(err)
	}
	io.WriteString(file, content)

}
