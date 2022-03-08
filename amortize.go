package main

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/http"
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

// getAmorHtml renders HTML from Postgres via queryAmor().
func getAmorHtml(w http.ResponseWriter, r *http.Request) {

	tmf, err := fs.Sub(embFS, "templates")
	if err != nil {
		log.Println("Error returning embedded directory templates:", err.Error())
	}

	f := "layoutAmor.html"

	_, err = fs.Stat(tmf, f)
	if err != nil {
		log.Println("Error locating", f, err.Error())
		http.Error(w, "HTML template issue", http.StatusInternalServerError)
		return
	}

	recs, err := queryAmor()
	if err != nil {
		log.Println("error calling queryAmor(): ", err.Error())
		http.Error(w, "Database connection issue", http.StatusServiceUnavailable)
		return
	}

	tmpl := template.Must(template.ParseFS(tmf, f))
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
	rows, err := db.Query("SELECT * FROM amortize ORDER BY payment_number;")
	if err != nil {
		log.Println("Error querying amortize table: ", err.Error())
		return nil, err
	}

	defer rows.Close()
	ctr := 0
	// amortize table will always have exactly 360 records
	snbs := make([]Amortize, 360)

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

		snbs[ctr] = snb
		ctr++
	}

	return snbs, nil
}
