package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

// getWebiHtml renders HTML from Postgres via queryWebi().
func getWebiHtml(w http.ResponseWriter, r *http.Request) {

	f := "./templates/layoutWebi.html"
	_, err := os.Stat(f)
	if err != nil {
		log.Println(f, "not found: ", err.Error())
		http.Error(w, "HTML template issue", http.StatusInternalServerError)
		return
	}

	recs, err := queryWebi()
	if err != nil {
		log.Println("error calling queryWebi(): ", err.Error())
		http.Error(w, "Database connection issue", http.StatusServiceUnavailable)
		return
	}

	tmpl := template.Must(template.ParseFiles(f))
	tmpl.Execute(w, recs) //recs[0])

}

// getWebiJson renders JSON from Postgres via queryWebi().
func getWebiJson(w http.ResponseWriter, r *http.Request) {

	recs, err := queryWebi()
	if err != nil {
		log.Println("error calling queryAmor(): ", err.Error())
		http.Error(w, "Database connection issue", http.StatusServiceUnavailable)
		return
	}

	// for no indent use: json.NewEncoder(w).Encode(snbs) - below has indent
	resp := json.NewEncoder(w)
	resp.SetIndent("", "    ")
	resp.Encode(recs)

}

// queryWebi queries and returns results from Postgres view web_info.
func queryWebi() (WebInfo, error) {

	snb := WebInfo{}
	row := db.QueryRow("SELECT * FROM web_info")

	if err := row.Scan(
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
	); err != nil {
		log.Println("error querying web_info table: ", err.Error())
		return snb, err
	}

	return snb, nil

}
