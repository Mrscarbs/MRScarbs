package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const Warning = "this is not financial advice do your own due diligence"

var file *os.File
var query_id = 0

type fundamental_response struct {
	Output     string
	Disclaimer string
}

func llm_economic(c *gin.Context) {
	log.SetOutput(file)
	query_id++
	var llm_type int = 1
	var llm_type_name string = "economic llm"
	var llm_entity_id = 0
	var llm_queryid = query_id

	query, _ := c.GetQuery("question")
	cmd := exec.Command("python", "economic_llm.py", fmt.Sprint(query))
	output, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}

	db, err_db := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")

	if err_db != nil {
		log.Println(err_db)
	}

	db.Exec("call stp_insert_tbl_llm_query_data(?,?,?,?,?,?)", llm_type, llm_entity_id, llm_queryid, llm_type_name, query, string(output))
	defer db.Close()
	response := fundamental_response{Output: string(output), Disclaimer: Warning}
	c.IndentedJSON(http.StatusOK, response)

}
func llm_news(c *gin.Context) {
	log.SetOutput(file)

	query_id++
	var llm_type int = 1
	var llm_type_name string = "newsllm"
	var llm_entity_id = 0
	var llm_queryid = query_id
	query, _ := c.GetQuery("question")

	query = string(query)

	cmd := exec.Command("python", "news_llm.py", fmt.Sprint(query))
	output, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	db, err3 := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")

	if err3 != nil {
		log.Println(err3)
	}

	db.Exec("call stp_insert_tbl_llm_query_data(?,?,?,?,?,?)", llm_type, llm_entity_id, llm_queryid, llm_type_name, query, string(output))
	defer db.Close()
	var response fundamental_response = fundamental_response{Output: string(output), Disclaimer: Warning}
	c.IndentedJSON(http.StatusOK, response)

}

func llm_fundamentals(c *gin.Context) {
	log.SetOutput(file)

	var llm_type int = 1
	var llm_type_name string = "fundamentals"
	var llm_entity_id = 0
	var llm_queryid = query_id

	query, _ := c.GetQuery("question")
	// log.SetOutput("llm_fundamentals.log")

	query = string(query)
	cmd := exec.Command("python", "vector_llm_fundamentals.py", fmt.Sprint(query))
	db, err_db := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")

	if err_db != nil {
		log.Println(err_db)
	}
	output, err := cmd.Output()
	if err != nil {
		log.Println(err)
	}
	db.Exec("call stp_insert_tbl_llm_query_data(?,?,?,?,?,?)", llm_type, llm_entity_id, llm_queryid, llm_type_name, query, string(output))

	var response fundamental_response = fundamental_response{Output: string(output), Disclaimer: Warning}
	c.IndentedJSON(http.StatusOK, response)

	defer db.Close()

}

func main() {
	log.Println("initializing llm api")
	// file, err := os.Create("llm_fundamentals.log")
	var err error
	file, err = os.Create("llm_api.log")
	log.SetOutput(file)
	if err != nil {
		log.Println(err)
	}
	// if err != nil {
	// 	log.Println(err)
	// }
	defer file.Close()
	router := gin.Default()
	router.GET("/fundamental_llm", llm_fundamentals)
	router.GET("/news-llm", llm_news)
	router.GET("/economic_llm", llm_economic)
	router.Run("localhost:8080")

}
