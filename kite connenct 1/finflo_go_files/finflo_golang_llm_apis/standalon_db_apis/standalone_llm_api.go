package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var log_file *os.File

type chat_details struct {
	Llm_tye       []int
	Entity_id     []int
	Queryid       []int
	Llm_type_name []string
	Query         []string
	Answer        []string
}

func get_chats_db(c *gin.Context) {

	log.SetOutput(log_file)
	eid, _ := c.GetQuery("eid")
	num_id, err := strconv.Atoi(eid)

	if err != nil {
		log.Println(err)
	}

	db, err := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")

	if err != nil {
		log.Println(err)
	}

	rows, err := db.Query("call stp_GetLLMQueryDataByEntityID(?)", num_id)
	if err != nil {
		log.Println(err)
	}

	Chatdet := chat_details{Llm_tye: []int{}, Queryid: []int{}, Llm_type_name: []string{}, Query: []string{}, Answer: []string{}}

	for rows.Next() {
		var nllm_type int
		var nentityid int
		var nqueeryid int
		var llm_type_name string
		var query string
		var answer string

		err := rows.Scan(&nllm_type, &nentityid, &nqueeryid, &llm_type_name, &query, &answer)
		if err != nil {
			log.Println(err)
		}
		Chatdet.Llm_tye = append(Chatdet.Llm_tye, nllm_type)
		Chatdet.Queryid = append(Chatdet.Queryid, nqueeryid)
		Chatdet.Llm_type_name = append(Chatdet.Llm_type_name, llm_type_name)
		Chatdet.Query = append(Chatdet.Query, query)
		Chatdet.Answer = append(Chatdet.Answer, answer)
	}

	c.IndentedJSON(http.StatusOK, Chatdet)

}
func main() {

	var err error
	log_file, err = os.Create("standalone_api.log")
	log.SetOutput(log_file)

	if err != nil {
		log.Println(err)
	}

	fmt.Println("initializing llm_api")
	router := gin.Default()
	router.GET("/get_prev_chats", get_chats_db)

	router.Run("localhost:8080")
}
