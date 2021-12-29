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

var (
	db *sql.DB
	//err error  Maybe use for init() sql.Open Statement?  Or just declare in init()?
)

func init() {

	dbh := os.Getenv("DB_HOST")
	dbn := os.Getenv("DB_NAME")
	dbp := os.Getenv("DB_PASS")
	dbt := os.Getenv("DB_TYPE")
	dbu := os.Getenv("DB_USER")
	connStr := dbt + "://" + dbu + ":" + dbp + "@" + dbh + "/" + dbn + "?sslmode=disable"
	db, _ = sql.Open("postgres", connStr)

}

func main() {

	http.HandleFunc("/", retrieveHTML)
	http.HandleFunc("/json", retrieveJSON)
	http.HandleFunc("/test/", test)
	fs := http.FileServer(http.Dir("./gopher"))
	http.Handle("/gopher/", http.StripPrefix("/gopher/", fs))
	http.ListenAndServe(":80", nil)

}

func queryDB() ([]Web_info, error) {

	rows, err := db.Query("SELECT * FROM web_info")
	if err != nil {
		fmt.Println("Error querying web_info")
		return nil, err
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
			return nil, err
		}
		snbs = append(snbs, snb)
	}

	return snbs, nil

}

func test(w http.ResponseWriter, r *http.Request) {

	jsonResponse := JSONResponse{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(jsonResponse)

}

func retrieveJSON(w http.ResponseWriter, r *http.Request) {

	recs, err := queryDB()
	if err != nil {
		fmt.Println("Error querying web_info")
		http.Error(w, "Database connection issue", http.StatusServiceUnavailable)
		return
	}

	// for no indent use: json.NewEncoder(w).Encode(snbs) - below has indent
	resp := json.NewEncoder(w)
	resp.SetIndent("", "    ")
	resp.Encode(recs)

}

func retrieveHTML(w http.ResponseWriter, r *http.Request) {

	f := "layout.html"
	_, err := os.Stat(f)
	if err != nil {
		fmt.Println("layout.html not found.")
		http.Error(w, "HTML template issue", http.StatusServiceUnavailable)
		return
	}

	recs, err := queryDB()
	if err != nil {
		fmt.Println("Error querying web_info")
		http.Error(w, "Database connection issue", http.StatusServiceUnavailable)
		return
	}
	// fmt.Print(snbs[0].Apr)  DEBUG ITEM
	tmpl := template.Must(template.ParseFiles(f))
	tmpl.Execute(w, recs[0])

}
