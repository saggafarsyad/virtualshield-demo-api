package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
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

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
