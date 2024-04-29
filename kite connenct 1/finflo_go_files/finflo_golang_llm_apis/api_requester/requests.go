package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	now := time.Now()
	file, _ := os.Create("requester.log")
	log.SetOutput(file)
	fmt.Println("requester initializing")
	url := "http://localhost:8080/fundamental_llm?question=give%20me%20the%20cashflow%20details%20for%20google"
	wg.Add(2)
	content := make_request(url, file, &wg)
	content2 := make_request(url, file, &wg)
	wg.Wait()
	fmt.Println(content)
	fmt.Println(content2)

	timetocomp := time.Since(now)
	fmt.Println(timetocomp)

}
func make_request(url string, file *os.File, wg *sync.WaitGroup) string {
	log.SetOutput(file)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	response, _ := http.DefaultClient.Do(req)
	content, _ := io.ReadAll(response.Body)
	wg.Done()
	return string(content)

}
