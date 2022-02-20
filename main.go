package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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
	r := mux.NewRouter()
	s := &http.Server{
		Handler:           r,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	r.HandleFunc("/mortgage", getWebiHtml)
	r.HandleFunc("/mortgagejson", getWebiJson)
	r.HandleFunc("/json1", getWebiJson)
	r.HandleFunc("/amortize", getAmorHtml)
	r.HandleFunc("/amortizejson", getAmorJson)
	r.HandleFunc("/json2", getAmorJson)
	r.HandleFunc("/test", test)

	r.Path("/prometheus").Handler(promhttp.Handler())

	r.Path("/webinfo").Handler(http.RedirectHandler("/mortgage", http.StatusMovedPermanently))
	r.Path("/webinfojson").Handler(http.RedirectHandler("/mortgagejson", http.StatusMovedPermanently))

	fr := http.FileServer(http.Dir("./htdocs"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fr))

	fs := http.FileServer(http.Dir("./proverbs"))
	r.PathPrefix("/proverbs/").Handler(http.StripPrefix("/proverbs", fs))

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
