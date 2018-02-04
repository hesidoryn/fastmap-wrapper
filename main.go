package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

var db *sql.DB

const (
	dbUser = "heorhi"
	dbName = "openstreetmap"

	minLatitude  = -90.0
	maxLatitude  = 90.0
	minLongitude = -180.0
	maxLongitude = 180.0

	errorBboxIsRequired         = "The parameter bbox is required, and must be of the form min_lon,min_lat,max_lon,max_lat."
	errorWrongLongitudeLatitude = "The latitudes must be between -90 and 90, longitudes between -180 and 180 and the minima must be less than the maxima."
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
	if len(bbox) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorBboxIsRequired))
		return
	}
	cs := strings.Split(bbox, ",")
	if len(cs) != 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorBboxIsRequired))
		return
	}

	minlonString, minlatString, maxlonString, maxlatString := cs[0], cs[1], cs[2], cs[3]
	minlon, err := strconv.ParseFloat(minlonString, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorBboxIsRequired))
		return
	}
	minlat, err := strconv.ParseFloat(minlatString, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorBboxIsRequired))
		return
	}
	maxlon, err := strconv.ParseFloat(maxlonString, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorBboxIsRequired))
		return
	}
	maxlat, err := strconv.ParseFloat(maxlatString, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorBboxIsRequired))
		return
	}

	if minlon > maxlon || minlat > maxlat ||
		minlon < minLongitude || minlon > maxLongitude ||
		maxlon < minLongitude || maxlon > maxLongitude ||
		minlat < minLatitude || minlat > minLatitude ||
		maxlat < minLatitude || maxlat > minLatitude {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorWrongLongitudeLatitude))
		return
	}

	fastmapQuery := fmt.Sprintf(fastmapQueryf,
		minlon,
		minlat,
		maxlon,
		maxlat,
	)
	rows, err := db.Query(fastmapQuery)
	checkErr(err)
	defer rows.Close()
	w.Header().Set("Content-Type", "application/xml")
	for rows.Next() {
		node := make([]byte, 0)
		rows.Scan(&node)
		w.Write(node)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
