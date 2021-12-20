package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Web_info struct {
	Total_amount           string
	Apr                    string
	Paid_thru              string
	Current_balance        string
	Principal_paid         string
	Percent_principal_paid string
	Interest_saved         string
	// Payment_date           string
	// Payment        string
	// Principal      string
	// Interest       string
	// Balance        string
	// Payment_number string
	// Percent_principal      string
	// Percent_interest string
}

func procUboo(body string) {

	var w []Web_info
	err := json.Unmarshal([]byte(body), &w)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", w)
}

func getUboo() string {

	resp, err := http.Get("http://uboo.supskifun.net/json")
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		print(err)
	}

	return (string(body))

}

func main() {

	info := getUboo()
	procUboo(info)

}
