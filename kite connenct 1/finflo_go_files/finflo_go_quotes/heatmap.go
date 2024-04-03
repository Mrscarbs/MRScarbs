package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	kitemodels "github.com/zerodha/gokiteconnect/v4/models"
	kiteticker "github.com/zerodha/gokiteconnect/v4/ticker"
)

var (
	ticker *kiteticker.Ticker
)

func onError(err error) {
	fmt.Println(" Websocket Error: ", err)
}

func onClose(code int, reason string) {
	fmt.Println("Close: ", code, reason)
}

func onConnect() {

	db1, err4 := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")
	if err4 != nil {
		fmt.Println(err4)
	}
	var (
		instToken = []uint32{}
	)

	rows, err5 := db1.Query("call stp_get_tbl_current_ltp()")
	if err5 != nil {
		fmt.Println(err5)
	}
	fmt.Println(rows)
	defer rows.Close()

	for rows.Next() {
		var token uint32
		// if err := rows.Scan(&token); err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		err := rows.Scan(&token)
		if err != nil {
			fmt.Println(err)
		}
		instToken = append(instToken, uint32(token))
	}

	fmt.Println("Connected")
	err := ticker.Subscribe(instToken)
	if err != nil {
		fmt.Println("err: ", err)
	}

	err = ticker.SetMode(kiteticker.ModeLTP, instToken)
	if err != nil {
		fmt.Println("err: ", err)
	}

}

func onTick(tick kitemodels.Tick) {

	time := time.Now()
	unix_time := time.Unix()

	ltp := tick.LastPrice
	time_fetched := unix_time
	instrument_token := tick.InstrumentToken
	fmt.Println("Tick: ", ltp)
	fmt.Println("time:", time_fetched)
	fmt.Println("instrument_token :", instrument_token)

	db, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")
	if err != nil {
		panic(err)
	}

	db.Exec("CALL stp_Update_Instrument(?, ?, ?)", int(instrument_token), int(unix_time), int(ltp))
	if err != nil {
		fmt.Println(err)
	}
	db.Close()

}

func onReconnect(attempt int, delay time.Duration) {
	fmt.Printf("Reconnect attempt %d in %fs\n", attempt, delay.Seconds())
}

func onNoReconnect(attempt int) {
	fmt.Printf("Maximum no of reconnect attempt reached: %d", attempt)
}

func main() {

	db2, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")
	if err != nil {
		fmt.Println(err)
	}
	var api_key string
	var secret_key string
	var api_provider string
	var access_Token string
	var last_purchase_time int
	var first_purchase_time int
	var hstorical int
	var instrument_type string
	var api_id int

	err2 := db2.QueryRow("call stp_get_api_config(?)", 1).Scan(&api_key, &secret_key, &api_provider, &access_Token, &last_purchase_time, &first_purchase_time, &hstorical, &instrument_type, &api_id)
	if err2 != nil {
		fmt.Println(err2)
	}
	db2.Close()
	fmt.Println(api_key)
	fmt.Println(access_Token)

	ticker = kiteticker.New(api_key, access_Token)

	// Assign callbacks
	ticker.OnError(onError)
	ticker.OnClose(onClose)
	ticker.OnConnect(onConnect)
	ticker.OnReconnect(onReconnect)
	ticker.OnNoReconnect(onNoReconnect)
	ticker.OnTick(onTick)

	// Start the connection
	ticker.Serve()

}
