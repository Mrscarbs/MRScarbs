package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-gota/gota/dataframe"
	_ "github.com/go-sql-driver/mysql"
)

var mut sync.Mutex

func main() {

	time_now := time.Now()

	var wg sync.WaitGroup
	f_log, err10 := os.Create("cash_flowlog.log")

	if err10 != nil {
		fmt.Println(err10)
	}

	_, _, symbols := get_snp500_companies()

	fmt.Println(symbols)
	file, err := os.Create("llm_doc.txt")
	if err != nil {
		fmt.Println(err)
	}
	n_file, err2 := os.Create(("llm_cashflow.txt"))
	if err2 != nil {
		fmt.Println(err)
	}
	defer n_file.Close()

	defer file.Close()
	for i := 0; i < 500; i++ {
		wg.Add(2)
		go get_fundamentals(symbols, i, file, &wg)
		go get_cashflow(symbols, i, n_file, &wg, f_log)
		time.Sleep(time.Millisecond * 1500)
	}
	wg.Wait()
	since := time.Since(time_now)
	fmt.Println(since)
	fmt.Println("done")

}
func file_writer(content string, file *os.File) {
	io.WriteString(file, content)

}

func get_snp500_companies() ([]string, []string, []string) {
	data, err := os.Open("sp500_companies.csv")

	if err != nil {
		fmt.Println(err)
	}

	df := dataframe.ReadCSV(data)
	fmt.Println(df)

	symbols := df.Col("Symbol")
	Longname := df.Col("Longname")
	Shortname := df.Col("Shortname")

	list_symbols := []string{}
	list_Longname := []string{}
	list_Shortname := []string{}
	file, err2 := os.Create("llm_companies_contex.txt")

	if err2 != nil {
		fmt.Println(err2)
	}

	for i := 0; i < symbols.Len(); i++ {

		symbol := symbols.Elem(i).String()
		Shortname := Shortname.Elem(i).String()
		Longname := Longname.Elem(i).String()
		llm_companies_contex := fmt.Sprintf("The symbol for the %s is %s and the short name is %s.", Longname, symbol, Shortname)

		io.WriteString(file, llm_companies_contex)

		list_symbols = append(list_symbols, symbol)
		list_Longname = append(list_Longname, Longname)
		list_Shortname = append(list_Shortname, Shortname)

	}

	defer file.Close()
	defer data.Close()

	return list_Longname, list_Shortname, list_symbols

}

func get_fundamentals(symbols []string, i int, file *os.File, wg *sync.WaitGroup) {
	symbol := symbols[i] // Symbol of the stock you want to fetch data for
	url := fmt.Sprintf("https://alpha-vantage.p.rapidapi.com/query?function=OVERVIEW&symbol=%s&datatype=json", symbol)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", "cf0a4a5f7amshc1e476b658e05e3p140de0jsnc4c8d2dcf811")
	req.Header.Add("X-RapidAPI-Host", "alpha-vantage.p.rapidapi.com")
	response, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	databytes, err2 := io.ReadAll(response.Body)

	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(string(databytes))

	var stock_datamap = make(map[string]interface{})

	err3 := json.Unmarshal(databytes, &stock_datamap)

	if err3 != nil {
		fmt.Println(err3)
	}

	pb := stock_datamap["PriceToBookRatio"]
	yearHigh := stock_datamap["52WeekHigh"]
	yearlow := stock_datamap["52WeekLow"]
	EVToEBITDA := stock_datamap["EVToEBITDA"]
	EVToRevenue := stock_datamap["EVToRevenue"]
	EPS := stock_datamap["EPS"]
	PERatio := stock_datamap["PERatio"]
	DividendYield := stock_datamap["DividendYield"]
	Sector := stock_datamap["Sector"]
	Industry := stock_datamap["Industry"]
	MarketCapitalization := stock_datamap["MarketCapitalization"]
	BookValue := stock_datamap["BookValue"]
	ProfitMargin := stock_datamap["ProfitMargin"]

	db, err_db := sql.Open("mysql", "root:Karma100%@tcp(localhost:3306)/finflo_base_db")

	if err_db != nil {
		fmt.Println(err_db)
	}
	time_st := time.Now()

	unix_time := time_st.Unix()
	mut.Lock()
	db.Exec("call stp_InsertStockData(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", pb, yearHigh, yearlow, EVToEBITDA, EVToRevenue, EPS, PERatio, DividendYield, Sector, Industry, MarketCapitalization, BookValue, ProfitMargin, unix_time, symbol)
	defer db.Close()
	defer response.Body.Close()
	llm_content := fmt.Sprintf("The pb or price to bookvalue for %s is %s, the year high for %s is %s, the year low for %s is%s,the EVToEBITDA for %s is %s,the EVToRevenue for %s is %s,the EPS for %s is %s,the pe or pe ratio or price to earning for %s is %s,the DividendYield for %s is %s,the sector for %s is %s,the industry for %s is %s,the marketcap or market capitalization for %s is %s, the bookvalue for %s is %s, the profitmargin for %s is %s.", symbol, pb, symbol, yearHigh, symbol, yearlow, symbol, EVToEBITDA, symbol, EVToRevenue, symbol, EPS, symbol, PERatio, symbol, DividendYield, symbol, Sector, symbol, Industry, symbol, MarketCapitalization, symbol, BookValue, symbol, ProfitMargin)
	fmt.Println(llm_content)
	file_writer(llm_content, file)

	data, err := os.ReadFile("llm_doc.txt")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
	mut.Unlock()
	wg.Done()
}

func get_cashflow(symbols []string, i int, n_file *os.File, wg *sync.WaitGroup, log_f *os.File) {

	type AutoGenerated_transfor_tools struct {
		Symbol        string `json:"symbol"`
		AnnualReports []struct {
			FiscalDateEnding                                          string `json:"fiscalDateEnding"`
			ReportedCurrency                                          string `json:"reportedCurrency"`
			OperatingCashflow                                         string `json:"operatingCashflow"`
			PaymentsForOperatingActivities                            string `json:"paymentsForOperatingActivities"`
			ProceedsFromOperatingActivities                           string `json:"proceedsFromOperatingActivities"`
			ChangeInOperatingLiabilities                              string `json:"changeInOperatingLiabilities"`
			ChangeInOperatingAssets                                   string `json:"changeInOperatingAssets"`
			DepreciationDepletionAndAmortization                      string `json:"depreciationDepletionAndAmortization"`
			CapitalExpenditures                                       string `json:"capitalExpenditures"`
			ChangeInReceivables                                       string `json:"changeInReceivables"`
			ChangeInInventory                                         string `json:"changeInInventory"`
			ProfitLoss                                                string `json:"profitLoss"`
			CashflowFromInvestment                                    string `json:"cashflowFromInvestment"`
			CashflowFromFinancing                                     string `json:"cashflowFromFinancing"`
			ProceedsFromRepaymentsOfShortTermDebt                     string `json:"proceedsFromRepaymentsOfShortTermDebt"`
			PaymentsForRepurchaseOfCommonStock                        string `json:"paymentsForRepurchaseOfCommonStock"`
			PaymentsForRepurchaseOfEquity                             string `json:"paymentsForRepurchaseOfEquity"`
			PaymentsForRepurchaseOfPreferredStock                     string `json:"paymentsForRepurchaseOfPreferredStock"`
			DividendPayout                                            string `json:"dividendPayout"`
			DividendPayoutCommonStock                                 string `json:"dividendPayoutCommonStock"`
			DividendPayoutPreferredStock                              string `json:"dividendPayoutPreferredStock"`
			ProceedsFromIssuanceOfCommonStock                         string `json:"proceedsFromIssuanceOfCommonStock"`
			ProceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet string `json:"proceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet"`
			ProceedsFromIssuanceOfPreferredStock                      string `json:"proceedsFromIssuanceOfPreferredStock"`
			ProceedsFromRepurchaseOfEquity                            string `json:"proceedsFromRepurchaseOfEquity"`
			ProceedsFromSaleOfTreasuryStock                           string `json:"proceedsFromSaleOfTreasuryStock"`
			ChangeInCashAndCashEquivalents                            string `json:"changeInCashAndCashEquivalents"`
			ChangeInExchangeRate                                      string `json:"changeInExchangeRate"`
			NetIncome                                                 string `json:"netIncome"`
		} `json:"annualReports"`
		QuarterlyReports []struct {
			FiscalDateEnding                                          string `json:"fiscalDateEnding"`
			ReportedCurrency                                          string `json:"reportedCurrency"`
			OperatingCashflow                                         string `json:"operatingCashflow"`
			PaymentsForOperatingActivities                            string `json:"paymentsForOperatingActivities"`
			ProceedsFromOperatingActivities                           string `json:"proceedsFromOperatingActivities"`
			ChangeInOperatingLiabilities                              string `json:"changeInOperatingLiabilities"`
			ChangeInOperatingAssets                                   string `json:"changeInOperatingAssets"`
			DepreciationDepletionAndAmortization                      string `json:"depreciationDepletionAndAmortization"`
			CapitalExpenditures                                       string `json:"capitalExpenditures"`
			ChangeInReceivables                                       string `json:"changeInReceivables"`
			ChangeInInventory                                         string `json:"changeInInventory"`
			ProfitLoss                                                string `json:"profitLoss"`
			CashflowFromInvestment                                    string `json:"cashflowFromInvestment"`
			CashflowFromFinancing                                     string `json:"cashflowFromFinancing"`
			ProceedsFromRepaymentsOfShortTermDebt                     string `json:"proceedsFromRepaymentsOfShortTermDebt"`
			PaymentsForRepurchaseOfCommonStock                        string `json:"paymentsForRepurchaseOfCommonStock"`
			PaymentsForRepurchaseOfEquity                             string `json:"paymentsForRepurchaseOfEquity"`
			PaymentsForRepurchaseOfPreferredStock                     string `json:"paymentsForRepurchaseOfPreferredStock"`
			DividendPayout                                            string `json:"dividendPayout"`
			DividendPayoutCommonStock                                 string `json:"dividendPayoutCommonStock"`
			DividendPayoutPreferredStock                              string `json:"dividendPayoutPreferredStock"`
			ProceedsFromIssuanceOfCommonStock                         string `json:"proceedsFromIssuanceOfCommonStock"`
			ProceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet string `json:"proceedsFromIssuanceOfLongTermDebtAndCapitalSecuritiesNet"`
			ProceedsFromIssuanceOfPreferredStock                      string `json:"proceedsFromIssuanceOfPreferredStock"`
			ProceedsFromRepurchaseOfEquity                            string `json:"proceedsFromRepurchaseOfEquity"`
			ProceedsFromSaleOfTreasuryStock                           string `json:"proceedsFromSaleOfTreasuryStock"`
			ChangeInCashAndCashEquivalents                            string `json:"changeInCashAndCashEquivalents"`
			ChangeInExchangeRate                                      string `json:"changeInExchangeRate"`
			NetIncome                                                 string `json:"netIncome"`
		} `json:"quarterlyReports"`
	}

	symbol := symbols[i]
	url := fmt.Sprintf("https://alpha-vantage.p.rapidapi.com/query?function=CASH_FLOW&symbol=%s&datatype=json&output_size=compact", symbol)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("X-RapidAPI-Key", "a36ccd5d62mshee86e99fbcb9664p11139ejsn7777332728b5")
	req.Header.Add("X-RapidAPI-Host", "alpha-vantage.p.rapidapi.com")

	response, err2 := http.DefaultClient.Do(req)
	if err2 != nil {
		fmt.Println(err2)
	}

	databytes, err3 := io.ReadAll(response.Body)
	if err3 != nil {
		fmt.Println(err3)
	}
	mut.Lock()
	var jsonunfolded AutoGenerated_transfor_tools
	err4 := json.Unmarshal(databytes, &jsonunfolded)
	if err4 != nil {
		fmt.Println(err4)
	}
	// fmt.Println(string(databytes))
	resp_un := jsonunfolded.AnnualReports
	// fmt.Println(resp_un)
	if len(resp_un) > 0 {
		one_year := resp_un[0]

		ocf := one_year.OperatingCashflow
		capex := one_year.CapitalExpenditures
		CashflowFromInvestment := one_year.CashflowFromInvestment
		CashflowFromFinancing := one_year.CashflowFromFinancing
		FiscalDateEnding := one_year.FiscalDateEnding
		NetIncome := one_year.NetIncome
		ChangeInInventory := one_year.ChangeInInventory
		llm_cashflow := fmt.Sprintf("The CashflowFromInvestment for %s is %s,the CashflowFromFinancing for %s is %s,the FiscalDateEnding for %s is %s,the NetIncome for %s is %s,the ChangeInInventory for %s is %s,the OperatingCashflow for %s is %s,the CapitalExpenditures for %s is %s.", symbol, CashflowFromInvestment, symbol, CashflowFromFinancing, symbol, FiscalDateEnding, symbol, NetIncome, symbol, ChangeInInventory, symbol, ocf, symbol, capex)
		io.WriteString(n_file, llm_cashflow)
		response.Body.Close()

	} else {
		fmt.Println("symbol is :", symbol)
		log.SetOutput(log_f)
		log.Println("list is greater for symbol:", symbol)
	}
	mut.Unlock()
	wg.Done()
}
