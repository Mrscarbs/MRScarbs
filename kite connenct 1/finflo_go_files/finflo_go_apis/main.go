package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

type Response struct {
	Index         string
	Tickes_sorted []string
	Interval_time string
}
type Response_quants struct {
	Sortino          float64
	Sharpe           float64
	Instrument_token int64
	Last_update_time int64
}

func get_top_nifty_coins(c *gin.Context) {

	interval, _ := c.GetQuery("interval")

	db, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")

	if err != nil {
		fmt.Println(err)
	}
	rows, err2 := db.Query("call stp_api_GetNifty50SortedByInterval(?)", string(interval))

	if err2 != nil {
		fmt.Println(err2)
	}
	var ticker_list = []string{}
	for rows.Next() {
		var ticker string
		err3 := rows.Scan(&ticker)
		if err3 != nil {
			fmt.Println(err3)
		}
		ticker_list = append(ticker_list, ticker)

	}
	response := Response{Index: "nifty50", Tickes_sorted: ticker_list, Interval_time: interval}

	c.IndentedJSON(http.StatusOK, response)

}

func get_top_banknifty_coins(c *gin.Context) {
	interval, _ := c.GetQuery("interval")

	db, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")

	if err != nil {
		fmt.Println(err)
	}

	rows, err2 := db.Query("call stp_api_GetBankNiftySortedByInterval(?)", string(interval))
	if err2 != nil {
		fmt.Println(err2)
	}
	var ticker_list = []string{}
	for rows.Next() {
		var ticker string
		err := rows.Scan(&ticker)
		if err != nil {
			fmt.Println(err)
		}
		ticker_list = append(ticker_list, ticker)

	}

	response := Response{Index: "banknifty", Tickes_sorted: ticker_list, Interval_time: interval}
	c.IndentedJSON(http.StatusOK, response)
}

func get_quant_statsbytoken(c *gin.Context) {

	var sortino float64
	var sharpe float64
	var instrumenttoken int64
	var last_update_time int64

	symbol, _ := c.GetQuery("symbol")
	var instrument_token int
	db, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")

	if err != nil {
		fmt.Println(err)
	}

	err2 := db.QueryRow("call stp_get_instrument_info_by_symbol(?)", symbol).Scan(&instrument_token)

	if err2 != nil {
		fmt.Println(err2)
	}

	err3 := db.QueryRow("call stp_GetQuantStatsByToken(?)", instrument_token).Scan(&sortino, &instrumenttoken, &sharpe, &last_update_time)
	if err3 != nil {
		fmt.Println(err3)
	}
	response := Response_quants{Sortino: sortino, Sharpe: sharpe, Instrument_token: instrumenttoken, Last_update_time: last_update_time}
	c.IndentedJSON(http.StatusOK, response)

}

func main() {

	fmt.Println("initializing_api")
	router := gin.Default()
	router.GET("/heatmap_nifty", get_top_nifty_coins)
	router.GET("/heatmap_banknifty", get_top_banknifty_coins)
	router.GET("/quant_stats", get_quant_statsbytoken)
	router.Run("localhost:8080")

}
