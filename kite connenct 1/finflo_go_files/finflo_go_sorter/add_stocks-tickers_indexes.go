package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/go-gota/gota/dataframe"
)

func Get_index_tickers(index_no int, interval string, wg *sync.WaitGroup) {
	// var stp_nifty = "call stp_insert_nifty50_sorted_tickers(?)"
	// var stp_bank_nifty = "call stp_insert_banknifty_sorted_tickers(?)"
	var stp string
	var url string
	if index_no == 1 {
		stp = "call stp_insert_nifty50_sorted_tickers(?,?)"
		url = "https://archives.nseindia.com/content/indices/ind_nifty50list.csv"

	}
	if index_no == 2 {
		stp = "call stp_insert_banknifty_sorted_tickers(?,?)"
		url = "https://archives.nseindia.com/content/indices/ind_niftybanklist.csv"
	}

	response, err := http.Get(url)

	if err != nil {
		fmt.Println(err)
	}

	content, err2 := io.ReadAll(response.Body)
	if err2 != nil {
		fmt.Println(err2)
	}

	str_content := string(content)

	// fmt.Println(str_content)

	df := dataframe.ReadCSV(strings.NewReader(str_content))

	// fmt.Println(df)
	var symbol_l = []string{}
	symbols := df.Col("Symbol")
	for i := 0; i < symbols.Len(); i++ {

		symbol := symbols.Elem(i).String()
		symbol_l = append(symbol_l, symbol)

	}
	db, err3 := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")
	if err != nil {
		fmt.Println(err3)
	}
	for _, val := range symbol_l {

		db.Exec(stp, val, interval)
		fmt.Println(val)

	}
	db.Close()
	wg.Done()

}
