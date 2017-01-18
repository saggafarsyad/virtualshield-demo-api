package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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

// ChartResponse struct
type ChartResponse struct {
	Result struct {
		MQ2Data   []ChartData `json:"mq2"`
		MQ135Data []ChartData `json:"mq135"`
	} `json:"result"`
	LastCreated int64 `json:"last_created"`
}

// ErrorResponse struct
type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	ErrorCode  string `json:"error_code"`
}

// SuccessResponse struct
type SuccessResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

var db *sql.DB

func routeIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, "{\"status\": \"ok\"}")
}

func routePush(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Read parameters
	dataMq135, err := strconv.ParseFloat(r.FormValue("getValue_mq135"), 32)
	if err != nil {
		writeErrorResponse(w, err, 400, "Invalid MQ135 Data", "INVALID_MQ135_DATA")
		return
	}
	dataMq2, err := strconv.ParseFloat(r.FormValue("getValue_mq2"), 32)
	if err != nil {
		writeErrorResponse(w, err, 400, "Invalid MQ2 Data", "INVALID_MQ2_DATA")
		return
	}
	log.Println("dataMq135:", dataMq135)
	log.Println("dataMq2:", dataMq2)
	// Insert to database
	insertData("data_mq135", float32(dataMq135))
	insertData("data_mq2", float32(dataMq2))
	// Write success response
	writeSuccessResponse(w, "Success")
}

func routeGetChart(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get latest timestamp
	queryParams := r.URL.Query()
	var timestamp int64
	timestamp, err := strconv.ParseInt(queryParams.Get("timestamp"), 10, 64)
	if err != nil {
		timestamp = 0
	}
	log.Println("Latest Data Timestamp:", timestamp)
	// Preparing response
	var chartData ChartResponse
	// Query MQ135 Chart data
	chartData.Result.MQ135Data = queryLatestData("data_mq135", timestamp)
	chartData.Result.MQ2Data = queryLatestData("data_mq2", timestamp)
	// If no data return
	if len(chartData.Result.MQ135Data) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Set timestamp
	chartData.LastCreated = chartData.Result.MQ135Data[len(chartData.Result.MQ135Data)-1].Created.Unix()
	// Convert chart data to string json
	chartJSON, _ := json.Marshal(chartData)
	// Write response
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(chartJSON))
}

func main() {
	// Get datasource string from env
	var datasource string
	datasource = os.Getenv("VIRTUALSHIELD_DATASOURCE")
	if datasource == "" {
		datasource = "root:root@/virtualshield_demo?parseTime=true"
	}
	log.Println(datasource)
	// Open database
	var err error
	db, err = sql.Open("mysql", datasource)
	checkErr(err)
	// Close database on finished
	defer db.Close()

	// Init Router
	router := httprouter.New()
	router.GET("/", routeIndex)
	router.POST("/data", routePush)
	router.GET("/chart", routeGetChart)
	log.Println("Listening to port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func insertData(tableName string, data float32) {
	// Prepare insert statement
	stmt, err := db.Prepare("INSERT `" + tableName + "` SET data=?")
	checkErr(err)
	// Execute statement
	_, err = stmt.Exec(data)
	checkErr(err)
}

func queryLatestData(tableName string, timestamp int64) []ChartData {
	// Prepare query statement
	stmt, err := db.Prepare("SELECT * FROM `" + tableName + "` WHERE UNIX_TIMESTAMP(created) > ? ORDER BY created ASC")
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

func writeSuccessResponse(w http.ResponseWriter, message string) {
	response := SuccessResponse{200, message}
	responseJSON, _ := json.Marshal(response)
	// Setup header
	w.Header().Set("Content-Type", "application/json")
	// Write response
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(responseJSON))
}

func writeErrorResponse(w http.ResponseWriter, err error, statusCode int, message string, errorCode string) {
	log.Println("[ERROR]", err)
	// Create error response
	errResponse := ErrorResponse{statusCode, message, errorCode}
	errResponseJSON, _ := json.Marshal(errResponse)
	// Setup header
	w.Header().Set("Content-Type", "application/json")
	// Write response
	w.WriteHeader(statusCode)
	fmt.Fprint(w, string(errResponseJSON))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
