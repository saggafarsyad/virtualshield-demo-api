package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
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

func main() {
	// Open database
	db, err := sql.Open("mysql", "root:root@/virtualshield_demo?parseTime=true")
	checkErr(err)
	// Close database on finished
	defer db.Close()

	insertData(db, "data_mq135", 23.34)
	insertData(db, "data_mq2", 15.79)

	var res Response
	res.Result.MQ2Data = queryLatestData(db, "data_mq2", 0)
	res.Result.MQ135Data = queryLatestData(db, "data_mq135", 0)

	response, _ := json.Marshal(res)
	fmt.Println(string(response))
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
