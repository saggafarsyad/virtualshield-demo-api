package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/julienschmidt/httprouter"
)

// ChartData struct
type ChartData struct {
	ID      int32     `json:"id"`
	Data    float32   `json:"data"`
	Created time.Time `json:"created"`
}

// Response struct
type Response struct {
	Result struct {
		MQ2Data   []ChartData `json:"mq2"`
		MQ135Data []ChartData `json:"mq135"`
	} `json:"result"`
}

// ErrorResponse struct
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	ErrorCode  string `json:"error_code"`
}

var db *sql.DB

func routeIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "{\"status\": \"ok\"}")
}

func main() {
	// Open database
	db, err := sql.Open("mysql", "root:root@/virtualshield_demo?parseTime=true")
	checkErr(err)
	// Close database on finished
	defer db.Close()

	// Init Router
	router := httprouter.New()
	router.GET("/", routeIndex)
	log.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func insertData(db *sql.DB, tableName string, data float32) {
	// Prepare insert statement
	stmt, err := db.Prepare("INSERT `" + tableName + "` SET data=?")
	checkErr(err)
	// Execute statement
	_, err = stmt.Exec(data)
	checkErr(err)
}

func queryLatestData(db *sql.DB, tableName string, timestamp int64) []ChartData {
	// Prepare query statement
	stmt, err := db.Prepare("SELECT * FROM `" + tableName + "` WHERE UNIX_TIMESTAMP(created) > ?")
	checkErr(err)
	// Execute Query
	rows, err := stmt.Query(timestamp)
	checkErr(err)
	// Scan Rows
	var dataset []ChartData
	for rows.Next() {
		row := ChartData{}
		rows.Scan(&row.ID, &row.Data, &row.Created)
		dataset = append(dataset, row)
	}
	return dataset
}

func writeErrorResponse(w http.ResponseWriter, err error, statusCode int, message string, errorCode string) {
	log.Println("[ERROR]", err)
	// Create error response
	errResponse := struct {
		Response ErrorResponse `json:"error"`
	}{
		Response: ErrorResponse{statusCode, message, errorCode},
	}
	errResponseJSON, _ := json.Marshal(errResponse)
	// Setup header
	w.Header().Set("Content-Type", "application/json")
	// Write response
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprint(w, string(errResponseJSON))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
