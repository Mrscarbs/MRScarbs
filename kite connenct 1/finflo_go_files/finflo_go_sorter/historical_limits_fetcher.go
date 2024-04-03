package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Time_diff struct {
	Api_id   int
	Interval string
	Limit    int
}

func (t *Time_diff) Database_fetcher() int {
	api_id := t.Api_id
	interval := t.Interval

	db, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")
	if err != nil {
		fmt.Println(err)
	}
	db.QueryRow("call stp_get_limits(?,?)", api_id, interval).Scan(&t.Limit)
	fmt.Println(t.Limit)
	db.Close()
	return t.Limit

}

func (t Time_diff) Last_time() (time.Time, int64, int64) {
	current_time := time.Now()
	current_timestamp := current_time.Unix()
	last_time := current_timestamp - int64(t.Limit)*60
	convertabletime := time.Unix(int64(last_time), 0)
	return convertabletime, last_time, current_timestamp

}
