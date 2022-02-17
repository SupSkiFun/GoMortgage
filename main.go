package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	// "github.com/prometheus/client_golang/prometheus"
	// "github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Amortize is used for a) Postgres data and b) Marshalling JSON.
type Amortize struct {
	Payment_date      string `json:"payment_date"`
	Payment           string `json:"payment"`
	Principal         string `json:"principal"`
	Interest          string `json:"interest"`
	Total_interest    string `json:"total_interest"`
	Balance           string `json:"balance"`
	Payment_number    string `json:"payment_number"`
	Percent_principal string `json:"percent_principal"`
	Percent_interest  string `json:"percent_interest"`
}

// JSONResponse works with test() to provide a simple JSON response.
type JSONResponse struct {
	Status string `json:"status"`
}

// WebInfo is used for a) Postgres data and b) Marshalling JSON.
type WebInfo struct {
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

// connectDB initializes a connection to Postgres.
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

// runWEB registers handlers and starts the web server.
func runWEB() {
	var err error
	m := http.NewServeMux()
	s := &http.Server{
		Handler:           m,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	m.Handle("/webinfo", http.HandlerFunc(getWebiHtml))
	// m.Handle("/bogus", http.RedirectHandler("/webinfo", http.StatusMovedPermanently))
	m.Handle("/webinfojson", http.HandlerFunc(getWebiJson))
	m.Handle("/json1", http.HandlerFunc(getWebiJson))
	m.Handle("/amortize", http.HandlerFunc(getAmorHtml))
	m.Handle("/amortizejson", http.HandlerFunc(getAmorJson))
	m.Handle("/json2", http.HandlerFunc(getAmorJson))
	m.Handle("/test/", http.HandlerFunc(test))
	m.Handle("/metrics", promhttp.Handler())

	fr := http.FileServer(http.Dir("./htdocs"))
	m.Handle("/", fr)

	fs := http.FileServer(http.Dir("./proverbs"))
	m.Handle("/proverbs/", http.StripPrefix("/proverbs", fs))

	err = s.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}

// test provides a basic JSON response.
func test(w http.ResponseWriter, r *http.Request) {
	jsonResponse := JSONResponse{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(jsonResponse)
}
