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
	dbUser = "heorhi"
	dbName = "openstreetmap"
)

func init() {
	dbinfo := fmt.Sprintf("user=%s dbname=%s sslmode=disable",
		dbUser, dbName)
	db, _ = sqlx.Open("postgres", dbinfo)
}

func main() {
	http.HandleFunc("/api/0.6/map", fastmap)
	http.ListenAndServe(":3001", nil)
}

func fastmap(w http.ResponseWriter, r *http.Request) {
	qp := r.URL.Query()
	bbox := qp.Get("bbox")
	if len(bbox) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - bbox is requierd!"))
	}

	cs := strings.Split(bbox, ",")
	if len(cs) < 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - minlon, minlat, maxlon, maxlat!"))
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
