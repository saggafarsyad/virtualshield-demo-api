package main

import (
	"database/sql"
	"fmt"
	"time"

	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
)

// ChartData struct
type ChartData struct {
	ID      int32     `json:"id"`
	Data    float32   `json:"data"`
	Created time.Time `json:"created"`
}

func main() {
	// Open database
	db, err := sql.Open("mysql", "root:root@/virtualshield_demo?parseTime=true")
	checkErr(err)
	// Close database on finished
	defer db.Close()
	// Prepare statement
	stmt, err := db.Prepare("INSERT `virtualshield_demo`.`data_mq135` SET data=?")
	checkErr(err)
	// Execute statement
	res, err := stmt.Exec(1)
	checkErr(err)
	// Get Id of latest data
	id, err := res.LastInsertId()
	checkErr(err)
	// Query data
	stmt, err = db.Prepare("SELECT * FROM data_mq135 WHERE id=?")
	checkErr(err)
	rows, err := stmt.Query(id)
	checkErr(err)
	// Scan rows
	for rows.Next() {
		var id int32
		var data float32
		var created time.Time
		err = rows.Scan(&id, &data, &created)
		checkErr(err)
		row := ChartData{id, data, created}
		// Print to JSON
		response, _ := json.Marshal(row)
		fmt.Println(string(response))
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
