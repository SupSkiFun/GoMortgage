/*
Mortgage Application produces a container that interacts with postgreSQL.
Metrics are exported to Prometheus.
*/
package main

import (
	// "context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	connectDB()
	runWEB()
}
