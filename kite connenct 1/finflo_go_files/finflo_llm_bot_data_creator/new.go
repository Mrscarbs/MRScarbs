package main

import (
	"fmt"
	"io"
	"net/http"
)

func new() {

	url := "https://alpha-vantage.p.rapidapi.com/query?interval=5min&function=TIME_SERIES_INTRADAY&symbol=MSFT&datatype=json&output_size=compact"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "cf0a4a5f7amshc1e476b658e05e3p140de0jsnc4c8d2dcf811")
	req.Header.Add("X-RapidAPI-Host", "alpha-vantage.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))

}
