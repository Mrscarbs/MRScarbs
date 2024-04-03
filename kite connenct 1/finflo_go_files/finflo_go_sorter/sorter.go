package main

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	kiteconnect "github.com/zerodha/gokiteconnect/v4"
)

var mut sync.Mutex

var token_list = []int64{}
var change = []float64{}
var token_list_5m = []int64{}
var change_5m = []float64{}

var token_list_m = []int64{}
var change_m = []float64{}

func main() {
	time_checker := time.Now()
	var wg sync.WaitGroup
	wg.Add(6)
	go Get_index_tickers(1, "5minute", &wg)
	go Get_index_tickers(2, "5minute", &wg)
	go Get_index_tickers(2, "60minute", &wg)
	go Get_index_tickers(1, "60minute", &wg)
	go Get_index_tickers(2, "minute", &wg)
	go Get_index_tickers(1, "minute", &wg)
	wg.Wait()

	var dict_token_change = make(map[string][]float64)
	dict_token_change["token"] = []float64{}
	dict_token_change["change"] = []float64{}
	dict_token_change["token_5m"] = []float64{}
	dict_token_change["change_5m"] = []float64{}
	dict_token_change["token_m"] = []float64{}
	dict_token_change["change_m"] = []float64{}

	token_arr := get_tokens(1)
	token_arr_banknifty := get_tokens(2)
	interval_arr := get_db_times()
	fmt.Println(interval_arr)
	var dict_times = make(map[string]string)

	dict_times["5minute"] = "minute"
	dict_times["60minute"] = "hour"
	dict_times["minute"] = "oneminute"

	for i := 0; i < len(token_arr); i++ {
		wg.Add(1)
		go sorter(token_arr[i], &wg, dict_token_change, interval_arr, dict_times)
		if i < 36 {

			wg.Add(1)
			go sorter_banknifty(token_arr_banknifty[i], &wg, interval_arr, dict_times)

		}

		time.Sleep(time.Millisecond * 400)

	}

	wg.Wait()
	// fmt.Println("before sort", token_list)
	fmt.Println(dict_token_change)
	fmt.Println(token_list)
	fmt.Println(change)
	fmt.Println(token_list_5m)
	fmt.Println(change_5m)
	fmt.Println("done")
	Since := time.Since(time_checker)
	fmt.Println(Since)

}

func sorter(tokens int64, wg1 *sync.WaitGroup, dict map[string][]float64, intervals []string, times map[string]string) {

	api_key, access_token := get_db_details()
	kc := kiteconnect.New(api_key)

	kc.SetAccessToken(access_token)
	for z := 0; z < len(intervals); z++ {
		var inter string
		if intervals[z] == times["60minute"] {

			inter = "60minute"

		} else if intervals[z] == times["5minute"] {
			inter = "5minute"
		} else if intervals[z] == times["minute"] {
			inter = "minute"

		} else {

			continue
		}
		fmt.Println(inter)
		token := tokens
		fmt.Println("initializing sorter")

		fetch_time1 := Time_diff{Api_id: 1, Interval: intervals[z]}
		fetch_time1.Database_fetcher()
		formated_time_last, un_time_last, current_timestamp := fetch_time1.Last_time()

		fmt.Println(formated_time_last)
		fmt.Println(un_time_last)
		fmt.Println(current_timestamp)

		from := time.Now()

		continuous := false
		oi := true
		historical, err_major := kc.GetHistoricalData(int(token), inter, formated_time_last, from, continuous, oi)
		if err_major != nil {
			fmt.Println(err_major)
		}
		fmt.Println("indexer")
		lates_close := historical[len(historical)-1].Close
		second_lates := historical[len(historical)-2].Close

		// fmt.Println(historical)
		abs_change := lates_close - second_lates
		mut.Lock()
		data_base, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")
		if err != nil {
			fmt.Println(err)
		}
		if inter == "60minute" {
			fmt.Println("start1")
			db_interval := inter
			db_token := token
			db_abs_change := abs_change

			data_base.Exec("call update_nifty50_sorted_points_change(?,?,?)", db_interval, db_token, db_abs_change)

			token_list = append(token_list, int64(token))
			change = append(change, float64(abs_change))
			dict["token"] = append(dict["token"], float64(token))
			dict["change"] = append(dict["change"], float64(abs_change))

		} else if inter == "5minute" {
			fmt.Println("start2")
			db_interval := inter
			db_token := token
			db_abs_change := abs_change

			data_base.Exec("call update_nifty50_sorted_points_change(?,?,?)", db_interval, db_token, db_abs_change)

			token_list_5m = append(token_list_5m, int64(token))
			change_5m = append(change_5m, float64(abs_change))
			dict["token_5m"] = append(dict["token_5m"], float64(token))
			dict["change_5m"] = append(dict["change_5m"], float64(abs_change))

		} else if inter == "minute" {
			fmt.Println("start2")
			db_interval := inter
			db_token := token
			db_abs_change := abs_change

			data_base.Exec("call update_nifty50_sorted_points_change(?,?,?)", db_interval, db_token, db_abs_change)

			token_list_m = append(token_list_m, int64(token))
			change_m = append(change_m, float64(abs_change))
			dict["token_m"] = append(dict["token_m"], float64(token))
			dict["change_m"] = append(dict["change_m"], float64(abs_change))

		}
		data_base.Close()

		mut.Unlock()

		// percentage_change := (abs_change/lates_close)*100
		fmt.Println(abs_change)

	}
	wg1.Done()

}
func sorter_banknifty(tokens int64, wg2 *sync.WaitGroup, intervals []string, times map[string]string) {

	api_key, acces_token := get_db_details()
	kc := kiteconnect.New(api_key)
	kc.SetAccessToken(acces_token)

	for i := 0; i < len(intervals); i++ {
		var inter string
		if intervals[i] == times["60minute"] {
			inter = "60minute"
		} else if intervals[i] == times["5minute"] {
			inter = "5minute"
		} else if intervals[i] == times["minute"] {
			inter = "minute"
		} else {
			continue
		}
		token := tokens

		fetch_time2 := Time_diff{Api_id: 1, Interval: intervals[i]}
		fetch_time2.Database_fetcher()
		formated_time_last, un_time_last, current_timestamp := fetch_time2.Last_time()
		fmt.Println(formated_time_last)
		fmt.Println(un_time_last)
		fmt.Println(current_timestamp)
		from := time.Now()

		continuous := false
		oi := true
		historical, err_major := kc.GetHistoricalData(int(token), inter, formated_time_last, from, continuous, oi)
		if err_major != nil {
			fmt.Println(err_major)
		}
		fmt.Println("indexer")
		lates_close := historical[len(historical)-1].Close
		second_lates := historical[len(historical)-2].Close
		abs_change := lates_close - second_lates
		mut.Lock()
		data_base, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")
		if err != nil {
			fmt.Println(err)
		}
		if inter == "60minute" {
			fmt.Println("start1")
			db_interval := inter
			db_token := token
			db_abs_change := abs_change

			data_base.Exec("call update_banknifty_sorted_points_change(?,?,?)", db_interval, db_token, db_abs_change)

		} else if inter == "5minute" {
			fmt.Println("start2")
			db_interval := inter
			db_token := token
			db_abs_change := abs_change

			data_base.Exec("call update_banknifty_sorted_points_change(?,?,?)", db_interval, db_token, db_abs_change)

		} else if inter == "minute" {
			fmt.Println("start2")
			db_interval := inter
			db_token := token
			db_abs_change := abs_change

			data_base.Exec("call update_banknifty_sorted_points_change(?,?,?)", db_interval, db_token, db_abs_change)

		}
		data_base.Close()

		mut.Unlock()

	}
	wg2.Done()
}

func get_db_details() (string, string) {

	var api_key string
	var secret_key string
	var api_provider string
	var access_Token string
	var last_purchase_time int
	var first_purchase_time int
	var hstorical int
	var instrument_type string
	var api_id int

	db, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")
	if err != nil {
		fmt.Println(err)
	}
	err2 := db.QueryRow("call stp_get_api_config(?)", 1).Scan(&api_key, &secret_key, &api_provider, &access_Token, &last_purchase_time, &first_purchase_time, &hstorical, &instrument_type, &api_id)
	if err2 != nil {
		fmt.Println(err2)
	}
	db.Close()
	return api_key, access_Token
}

func get_tokens(index_no int64) []int64 {

	var token_list = []int64{}
	var rows *sql.Rows
	var err2 error
	db, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")

	if err != nil {
		fmt.Println(err)
	}
	if index_no == 1 {
		rows, err2 = db.Query("call stp_GetNifty50InstrumentTokens()")
	}
	if index_no == 2 {
		rows, err2 = db.Query("call stp_GetBankNiftyInstrumentTokens()")
	}
	// rows, err2 := db.Query("call stp_GetNifty50InstrumentTokens()")

	if err2 != nil {
		fmt.Println(err2)
	}
	for rows.Next() {
		var new_token int
		err := rows.Scan(&new_token)
		if err != nil {
			fmt.Println(err)
		}
		token_list = append(token_list, int64(new_token))

	}
	rows.Close()
	db.Close()

	return token_list

}

func get_db_times() []string {
	var list_intervals = []string{}
	var Interval string

	db, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")
	if err != nil {
		fmt.Println(err)
	}

	row, err1 := db.Query("call stp_GetSInterval(?)", 1)
	if err != nil {
		fmt.Println(err1)
	}

	for row.Next() {

		err2 := row.Scan(&Interval)
		if err2 != nil {
			fmt.Println(err2)
		}
		list_intervals = append(list_intervals, string(Interval))
	}

	return list_intervals

}
