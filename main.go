/*
Mortgage Application produces a container that interacts with postgreSQL.
Metrics are exported to Prometheus.
*/
package main

import (
	// "context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

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

// test provides a basic JSON response.
func test(w http.ResponseWriter, r *http.Request) {
	jsonResponse := JSONResponse{
		Status: "OK",
	}
	json.NewEncoder(w).Encode(jsonResponse)
}
