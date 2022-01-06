package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Amortize struct {
	Payment_date      string `json:"payment_date"`
	Payment           string `json:"payment"`
	Principal         string `json:"principal"`
	Interest          string `json:"interest "`
	Total_interest    string `json:"total_interest"`
	Balance           string `json:"balance"`
	Payment_number    string `json:"payment_number"`
	Percent_principal string `json:"percent_principal"`
	Percent_interest  string `json:"percent_interest"`
}

type JSONResponse struct {
	Status string `json:"status"`
}

type Web_info struct {
	Total_amount           string `json:"total_amount"`
	Apr                    string `json:"apr"`
	Paid_thru              string `json:"paid_thru"`
	Current_balance        string `json:"current_balance"`
	Principal_paid         string `json:"principal_paid"`
	Percent_principal_paid string `json:"percent_principal_paid"`
	Interest_saved         string `json:"interest_saved"`
	Payment_date           string `json:"payment_date"`
	Payment                string `json:"payment"`
	Principal              string `json:"principal"`
	Interest               string `json:"interest"`
	Balance                string `json:"balance"`
	Payment_number         string `json:"payment_number"`
	Percent_principal      string `json:"percent_principal"`
	Percent_interest       string `json:"percent_interest"`
}

var db *sql.DB

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {

	connectDB()
	runWEB()

}

func connectDB() {
	dbh := os.Getenv("DB_HOST")
	dbn := os.Getenv("DB_NAME")
	dbp := os.Getenv("DB_PASS")
	dbt := os.Getenv("DB_TYPE")
	dbu := os.Getenv("DB_USER")
	connStr := dbt + "://" + dbu + ":" + dbp + "@" + dbh + "/" + dbn + "?sslmode=disable"
	db, _ = sql.Open("postgres", connStr)
	err := db.Ping()
	if err != nil {
		log.Println("error opening database connection: ", err.Error())
	} else {
		log.Println("database ping ok")
	}
}

func runWEB() {
	http.HandleFunc("/", getWebiHtml)
	http.HandleFunc("/webinfo", getWebiHtml)
	http.HandleFunc("/webinfojson", getWebiJson)
	http.HandleFunc("/json1", getWebiJson)
	http.HandleFunc("/amortize", getAmorHtml)
	http.HandleFunc("/amortizejson", getAmorJson)
	http.HandleFunc("/json2", getAmorJson)
	http.HandleFunc("/test/", test)
	fs := http.FileServer(http.Dir("./gopher"))
	http.Handle("/gopher/", http.StripPrefix("/gopher/", fs))
	log.Fatal(http.ListenAndServe(":80", nil))
}

func test(w http.ResponseWriter, r *http.Request) {

	jsonResponse := JSONResponse{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(jsonResponse)

}
