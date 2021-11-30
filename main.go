package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type JSONResponse struct {
	Status string
}

type Web_info struct {
	Total_amount           string
	Apr                    string
	Paid_thru              string
	Current_balance        string
	Principal_paid         string
	Percent_principal_paid string
	Interest_saved         string
	Payment_date           string
	Payment                string
	Principal              string
	Interest               string
	Balance                string
	Payment_number         string
	Percent_principal      string
	Percent_interest       string
}

var db *sql.DB
var err error

func init() {

	dbh := os.Getenv("DB_HOST")
	dbn := os.Getenv("DB_NAME")
	dbp := os.Getenv("DB_PASS")
	dbt := os.Getenv("DB_TYPE")
	dbu := os.Getenv("DB_USER")
	connStr := dbt + "://" + dbu + ":" + dbp + "@" + dbh + "/" + dbn + "?sslmode=disable"
	db, err = sql.Open("postgres", connStr)

}

func main() {

	http.HandleFunc("/", retrieveRecord)
	http.HandleFunc("/json", retrieveJSON)
	http.HandleFunc("/test/", test)
	fs := http.FileServer(http.Dir("./gopher"))
	http.Handle("/gopher/", http.StripPrefix("/gopher/", fs))
	http.ListenAndServe(":80", nil)

}

func test(w http.ResponseWriter, r *http.Request) {

	jsonResponse := JSONResponse{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(jsonResponse)
}

func retrieveJSON(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT * FROM web_info")
	if err != nil {
		fmt.Println("Error querying web_info")
	}
	defer rows.Close()
	snbs := make([]Web_info, 0)

	for rows.Next() {
		snb := Web_info{}
		err := rows.Scan(
			&snb.Total_amount,
			&snb.Apr,
			&snb.Paid_thru,
			&snb.Current_balance,
			&snb.Principal_paid,
			&snb.Percent_principal_paid,
			&snb.Interest_saved,
			&snb.Payment_date,
			&snb.Payment,
			&snb.Principal,
			&snb.Interest,
			&snb.Balance,
			&snb.Payment_number,
			&snb.Percent_principal,
			&snb.Percent_interest,
		)
		if err != nil {
			log.Println(err)
			return
		}
		snbs = append(snbs, snb)
	}
	// below works
	//json.NewEncoder(w).Encode(snbs)
	//below works with indent
	resp := json.NewEncoder(w)
	resp.SetIndent("", "    ")
	resp.Encode(snbs)

}

func retrieveRecord(w http.ResponseWriter, r *http.Request) {

	rows, err := db.Query("SELECT * FROM web_info")
	if err != nil {
		fmt.Println("Error querying web_info")
	}
	defer rows.Close()
	snbs := make([]Web_info, 0)

	for rows.Next() {
		snb := Web_info{}
		err := rows.Scan(
			&snb.Total_amount,
			&snb.Apr,
			&snb.Paid_thru,
			&snb.Current_balance,
			&snb.Principal_paid,
			&snb.Percent_principal_paid,
			&snb.Interest_saved,
			&snb.Payment_date,
			&snb.Payment,
			&snb.Principal,
			&snb.Interest,
			&snb.Balance,
			&snb.Payment_number,
			&snb.Percent_principal,
			&snb.Percent_interest,
		)
		if err != nil {
			log.Println(err)
			return
		}
		snbs = append(snbs, snb)
	}
	// fmt.Print(snbs[0].Apr)  DEBUG ITEM
	tmpl := template.Must(template.ParseFiles("layout.html"))
	tmpl.Execute(w, snbs[0])
}
