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
	qp := r.URL.Query()
	bbox := qp.Get("bbox")
	if len(bbox) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - bbox is requierd!"))
	}

	cs := strings.Split(bbox, ",")
	if len(cs) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - minlot, minlat, maxlot, maxlot!"))
	}

	fastmapQuery := fmt.Sprintf(fastmapQueryf, cs[0], cs[1], cs[2], cs[3])
	rows, err := db.Query(fastmapQuery)
	checkErr(err)

	xmlNodes := make([]string, 0)
	for rows.Next() {
		node := ""
		rows.Scan(&node)
		xmlNodes = append(xmlNodes, node)
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
