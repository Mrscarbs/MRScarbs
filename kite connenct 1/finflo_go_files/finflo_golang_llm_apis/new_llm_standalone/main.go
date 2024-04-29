package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const warning string = "This is not financiall advice do your ownn due diligence before investing or trading"

type chat_details struct {
	Llm_tye       []int
	Entity_id     []int
	Queryid       []int
	Llm_type_name []string
	Query         []string
	Answer        []string
}

type response struct {
	question   string
	answer     string
	disclaimer string
}

var mut sync.Mutex

func main() {

	var wg sync.WaitGroup

	log_file, err := os.Create("standalone_llm.log")
	log.SetOutput(log_file)
	if err != nil {
		panic(err)
	}
	fmt.Println("ask your question")
	// var question_list = []string{"give me some news for bitcoin and its url"}

	var llm_type string
	var llm_type2 string
	var llm_type3 string
	var llm_type4 string
	var input_question string
	var input_question2 string
	var input_question3 string
	var input_question4 string

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter your llm type 1 is news, 2 is fundamental, 3 is economic : ")
	llm_type, _ = reader.ReadString('\n')
	llm_type = strings.TrimSpace(llm_type)
	fmt.Println("Enter your question: ")
	input_question, _ = reader.ReadString('\n')

	fmt.Println("Enter your llm type 1 is news, 2 is fundamental, 3 is economic : ")
	llm_type2, _ = reader.ReadString('\n')
	llm_type2 = strings.TrimSpace(llm_type2)
	if llm_type2 == "1" || llm_type2 == "2" {
		fmt.Println("working")
	} else {
		fmt.Println("not working")
	}
	fmt.Println("Enter your question: ")
	input_question2, _ = reader.ReadString('\n')

	fmt.Println("Enter your llm type 1 is news, 2 is fundamental, 3 is economic : ")
	llm_type3, _ = reader.ReadString('\n')
	llm_type3 = strings.TrimSpace(llm_type3)

	fmt.Println("Enter your question: ")
	input_question3, _ = reader.ReadString('\n')
	fmt.Println("Enter your llm type 1 is news, 2 is fundamental, 3 is economic : ")
	llm_type4, _ = reader.ReadString('\n')
	llm_type4 = strings.TrimSpace(llm_type4)
	fmt.Println("Enter your question: ")
	input_question4, _ = reader.ReadString('\n')

	var question = make(map[string][]string)
	question["question"] = []string{"give me 5 news for bitcoin and its url", "give me the cashflow details for google"}

	question["llm_id"] = []string{"1", "2"}

	question["question"] = append(question["question"], input_question, input_question2, input_question3, input_question4)

	question["llm_id"] = append(question["llm_id"], llm_type, llm_type2, llm_type3, llm_type4)

	current_time := time.Now()
	fmt.Println(question)
	for i := 0; i < len(question["question"]); i++ {
		question_user := question["question"][i]
		llm_id := question["llm_id"][i]
		if llm_id == "1" {
			wg.Add(1)
			go news_llm(question_user, log_file, &wg)
		}
		if llm_id == "2" {
			wg.Add(1)
			go fundamendamentaal_llm(question_user, log_file, &wg)
		}
		if llm_id == "3" {
			wg.Add(1)
			go economic_llm(question_user, log_file, &wg)
		}

	}
	wg.Wait()
	scince2 := time.Since(current_time)
	fmt.Println(scince2)
	chat_reader := bufio.NewReader(os.Stdin)
	fmt.Println("if you want to get your previous chats type y if not type n")
	input_chat, _ := chat_reader.ReadString('\n')
	input_chat = strings.TrimSpace(input_chat)
	eid := "1"
	if input_chat == "y" {
		chat_dett := get_chat_history(eid, log_file)

		chat_json, err := json.MarshalIndent(chat_dett, "", "\t")

		if err != nil {
			log.Println(err)
		}

		fmt.Println("these are your previous chats")

		fmt.Println(string(chat_json))
	} else {
		fmt.Println("thank you for using traders pilot llm chat bot")
	}

	since := time.Since(current_time)
	fmt.Println(since)

}

func news_llm(question string, log_file *os.File, wg *sync.WaitGroup) {

	var llm_type = 1
	var llm_entity_id = 1
	var llm_queryid = 1
	var llm_type_name = "news"
	var query = question

	log.SetOutput(log_file)
	cmd := exec.Command("python", "news_llm.py", fmt.Sprint(question))

	output, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	mut.Lock()
	res := response{question: question, answer: string(output), disclaimer: warning}
	send_to_db(log_file, llm_type, llm_entity_id, llm_queryid, llm_type_name, query, output)
	fmt.Println(res)
	mut.Unlock()
	wg.Done()
}

func fundamendamentaal_llm(question string, log_file *os.File, wg *sync.WaitGroup) {

	var llm_type = 2
	var llm_entity_id = 1
	var llm_queryid = 1
	var llm_type_name = "fundamental"
	var query = question

	cmd := exec.Command("python", "vector_llm_fundamentals.py", fmt.Sprint(question))
	log.SetOutput(log_file)

	output, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	mut.Lock()
	res := response{question: question, answer: string(output), disclaimer: warning}
	send_to_db(log_file, llm_type, llm_entity_id, llm_queryid, llm_type_name, query, output)
	fmt.Println(res)
	mut.Unlock()
	wg.Done()
}

func economic_llm(question string, log_file *os.File, wg *sync.WaitGroup) {

	var llm_type = 3
	var llm_entity_id = 1
	var llm_queryid = 1
	var llm_type_name = "economic"
	var query = question

	log.SetOutput(log_file)

	cmd := exec.Command("python", "economic_llm.py", fmt.Sprint(question))

	output, err := cmd.Output()

	if err != nil {
		log.Println(err)
	}
	mut.Lock()
	res := response{question: question, answer: string(output), disclaimer: warning}
	send_to_db(log_file, llm_type, llm_entity_id, llm_queryid, llm_type_name, query, output)
	fmt.Println(res)
	mut.Unlock()
	wg.Done()
}

func send_to_db(log_file *os.File, llm_type int, llm_entity_id int, llm_queryid int, llm_type_name string, query string, output []byte) {

	log.SetOutput(log_file)

	db, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")
	if err != nil {
		log.Println(err)
	}

	db.Exec("call stp_insert_tbl_llm_query_data(?,?,?,?,?,?)", llm_type, llm_entity_id, llm_queryid, llm_type_name, query, string(output))

	defer db.Close()

}
func get_chat_history(eid string, log_file *os.File) chat_details {
	log.SetOutput(log_file)
	url := fmt.Sprintf("http://localhost:8080/get_prev_chats?eid=%s", eid)
	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}

	databytes, err := io.ReadAll(res.Body)

	if err != nil {
		log.Println(err)
	}
	var chat_json_unfolded chat_details
	err_json := json.Unmarshal(databytes, &chat_json_unfolded)
	if err != nil {
		log.Println(err_json)
	}
	return chat_json_unfolded
}
