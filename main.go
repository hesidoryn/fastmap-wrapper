package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	dbUser = "heorhi"
	dbName = "openstreetmap"
)

func init() {
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable",
		dbUser, dbName)
	db, _ = sql.Open("postgres", dbinfo)
}

func main() {
	http.HandleFunc("/api/0.6/map", fastmap)
	http.ListenAndServe(":3001", nil)
}

func fastmap(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()
	bbox := qp.Get("bbox")
	cs := strings.Split(bbox, ",")
	if len(cs) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - minlon, minlat, maxlon, maxlat!"))
		return
	}

	fastmapQuery := fmt.Sprintf(fastmapQueryf, cs[0], cs[1], cs[2], cs[3])
	rows, err := db.Query(fastmapQuery)
	checkErr(err)

	w.Header().Set("Content-Type", "application/xml")
	for rows.Next() {
		node := make([]byte, 0)
		rows.Scan(&node)
		w.Write(node)
	}
	rows.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
