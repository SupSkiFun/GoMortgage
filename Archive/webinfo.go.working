package main

import (
	"encoding/json"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

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

// getWebiHtml renders HTML from Postgres via queryWebi().
func getWebiHtml(w http.ResponseWriter, r *http.Request) {

	tmf, err := fs.Sub(embFS, "templates")
	if err != nil {
		log.Println("Error returning embedded directory templates:", err.Error())
	}

	f := "layoutWebi.html"

	_, err = fs.ReadFile(tmf, f)
	if err != nil {
		log.Println("Error reading", f, err.Error())
		http.Error(w, "HTML template issue", http.StatusInternalServerError)
		return
	}

	recs, err := queryWebi()
	if err != nil {
		log.Println("error calling queryWebi(): ", err.Error())
		http.Error(w, "Database connection issue", http.StatusServiceUnavailable)
		return
	}

	tmpl := template.Must(template.ParseFS(tmf, f))
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
	// pgsql-view web_info only has one record
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
