/*
Mortgage Application produces a container that interacts with postgreSQL.
Metrics are exported to Prometheus.
*/
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/lib/pq"
)

// JSONResponse works with test() to provide a simple JSON response.
type JSONResponse struct {
	Status string `json:"status"`
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
	// var err error	Remove this if I keep using go routines
	rtr := mux.NewRouter()
	rtr.Use(prometheusMiddleware)
	s := &http.Server{
		Handler:           rtr,
		IdleTimeout:       30 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		ReadHeaderTimeout: 2 * time.Second,
	}

	rtr.HandleFunc("/mortgage", getWebiHtml)
	rtr.HandleFunc("/mortgagejson", getWebiJson)
	rtr.HandleFunc("/json1", getWebiJson)
	rtr.HandleFunc("/amortize", getAmorHtml)
	rtr.HandleFunc("/amortizejson", getAmorJson)
	rtr.HandleFunc("/json2", getAmorJson)
	rtr.HandleFunc("/test", test)

	rtr.Path("/prometheus").Handler(promhttp.Handler())
	rtr.Path("/webinfo").Handler(http.RedirectHandler("/mortgage", http.StatusMovedPermanently))
	rtr.Path("/webinfojson").Handler(http.RedirectHandler("/mortgagejson", http.StatusMovedPermanently))
	rtr.Path("/cake").Handler(http.RedirectHandler("/cake.html", http.StatusMovedPermanently))

	fr := http.FileServer(http.Dir("./proverbs"))
	rtr.PathPrefix("/proverbs/").Handler(http.StripPrefix("/proverbs/", fr))
	fs := http.FileServer(http.Dir("./htdocs"))
	rtr.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
			return
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	sig := <-sigs
	secs := 3 * time.Second
	log.Println("Terminating in", secs, ": Received signal", sig)
	time.Sleep(secs)
	s.Shutdown(context.TODO())

	/*  Original Startup is below - can replace above goroutine / closure
	    and also the sigs / signal channel portion directly above

			err = s.ListenAndServe()
			if err != nil {
				log.Fatal(err)
			}
	*/
}

// test provides a basic JSON response.
func test(w http.ResponseWriter, r *http.Request) {
	jsonResponse := JSONResponse{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(jsonResponse)
}
