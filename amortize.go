package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
)

// getAmorHtml renders HTML from Postgres via queryAmor().
func getAmorHtml(w http.ResponseWriter, r *http.Request) {

	f := "layoutAmor.html"
	_, err := os.Stat(f)
	if err != nil {
		log.Println(f, "not found: ", err.Error())
		http.Error(w, "HTML template issue", http.StatusInternalServerError)
		return
	}

	recs, err := queryAmor()
	if err != nil {
		log.Println("error calling queryAmor(): ", err.Error())
		http.Error(w, "Database connection issue", http.StatusServiceUnavailable)
		return
	}

	tmpl := template.Must(template.ParseFiles(f))
	tmpl.Execute(w, recs) //recs[0])

}

// getAmorJson renders JSON from Postgres via queryAmor().
func getAmorJson(w http.ResponseWriter, r *http.Request) {

	recs, err := queryAmor()
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

// queryAmor queries and returns results from Postgres table amortize.
func queryAmor() ([]Amortize, error) {

	rows, err := db.Query("SELECT * FROM amortize")
	if err != nil {
		log.Println("Error querying amortize table: ", err.Error())
		return nil, err
	}

	defer rows.Close()

	snbs := make([]Amortize, 0)

	for rows.Next() {
		snb := Amortize{}
		err := rows.Scan(
			&snb.Payment_date,
			&snb.Payment,
			&snb.Principal,
			&snb.Interest,
			&snb.Total_interest,
			&snb.Balance,
			&snb.Payment_number,
			&snb.Percent_principal,
			&snb.Percent_interest,
		)
		if err != nil {
			log.Println("error processing rows from amortize table: ", err.Error())
			return nil, err
		}
		snbs = append(snbs, snb)
	}

	return snbs, nil

}
