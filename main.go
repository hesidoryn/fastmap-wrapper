package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

const (
	DB_USER     = "sidoringi"
	DB_PASSWORD = "pass"
	DB_NAME     = "gis"
)

func init() {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		DB_USER, DB_PASSWORD, DB_NAME)
	db, _ = sqlx.Open("postgres", dbinfo)
}

func main() {
	http.HandleFunc("/api/0.6/map", foo)
	http.ListenAndServe(":3001", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(fastmapQuery)
	checkErr(err)

	nodeCount := 48343
	xmlNodes := make([]string, nodeCount)
	i := 0
	for rows.Next() {
		rows.Scan(&xmlNodes[i])
		i++
	}

	xml := strings.Join(xmlNodes, "")
	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(xml))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
