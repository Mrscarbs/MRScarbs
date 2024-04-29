package main

import (
	"fmt"
	"os/exec"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	now := time.Now()
	fmt.Println("initializing multi qna")
	questions := []string{"give me some news about bitcoin", "give me some news in india", "give me some news about bitcoin", "give me some news for india", "give me some news from usa"}
	for _, val := range questions {
		wg.Add(1)
		go get_answer(val, &wg)
	}
	wg.Wait()
	dur := time.Since(now)

	fmt.Println(dur)

}

func get_answer(question string, wg *sync.WaitGroup) {
	cmd := exec.Command("python", "news_llm.py", fmt.Sprint(question))

	output, err := cmd.Output()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(output))
	wg.Done()

}
