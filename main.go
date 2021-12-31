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

	http.HandleFunc("/", getWebiHtml)
	http.HandleFunc("/jsonWebInfo", getWebiJson)
	http.HandleFunc("/amortize", getAmorHtml)
	http.HandleFunc("/jsonAmortize", getAmorJson)
	http.HandleFunc("/test/", test)
	http.HandleFunc("/json1", getWebiJson)
	http.HandleFunc("/json2", getAmorJson)
	fs := http.FileServer(http.Dir("./gopher"))
	http.Handle("/gopher/", http.StripPrefix("/gopher/", fs))
	http.ListenAndServe(":80", nil)

}

func getAmorHtml(w http.ResponseWriter, r *http.Request) {

	f := "layoutAmor.html"
	_, err := os.Stat(f)
	if err != nil {
		fmt.Println(f, "not found.")
		http.Error(w, "HTML template issue", http.StatusServiceUnavailable)
		return
	}

	recs, err := queryAmor()
	if err != nil {
		fmt.Println("Error querying web_info")
		http.Error(w, "Database connection issue", http.StatusServiceUnavailable)
		return
	}

	// fmt.Printf("%+v", recs)

	tmpl := template.Must(template.ParseFiles(f))
	tmpl.Execute(w, recs) //recs[0])

}

func getAmorJson(w http.ResponseWriter, r *http.Request) {

	recs, err := queryAmor()
	if err != nil {
		fmt.Println("Error querying amortize")
		http.Error(w, "Database connection issue", http.StatusServiceUnavailable)
		return
	}

	// for no indent use: json.NewEncoder(w).Encode(snbs) - below has indent
	resp := json.NewEncoder(w)
	resp.SetIndent("", "    ")
	resp.Encode(recs)

}

func getWebiHtml(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		fmt.Println(r.URL.Path, "not found.")
		http.Error(w, r.URL.Path+" not found", http.StatusNotFound)
		return
	}

	f := "layoutWebi.html"
	_, err := os.Stat(f)
	if err != nil {
		fmt.Println(f, "not found.")
		http.Error(w, "HTML template issue", http.StatusServiceUnavailable)
		return
	}

	recs, err := queryWebi()
	if err != nil {
		fmt.Println("Error querying web_info")
		http.Error(w, "Database connection issue", http.StatusServiceUnavailable)
		return
	}

	tmpl := template.Must(template.ParseFiles(f))
	tmpl.Execute(w, recs) //recs[0])

}

func getWebiJson(w http.ResponseWriter, r *http.Request) {

	recs, err := queryWebi()
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

func queryAmor() ([]Amortize, error) {

	rows, err := db.Query("SELECT * FROM amortize")
	if err != nil {
		fmt.Println("Error querying amortize")
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
			log.Println(err)
			return nil, err
		}
		snbs = append(snbs, snb)
	}

	return snbs, nil

}

func queryWebi() (Web_info, error) {

	snb := Web_info{}
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
		fmt.Println("Error querying web_info")
		return snb, err
	}

	return snb, nil

}

func test(w http.ResponseWriter, r *http.Request) {

	jsonResponse := JSONResponse{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(jsonResponse)

}
