package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// ChartData struct
type ChartData struct {
	ID      int32
	Data    float32
	Created time.Time
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
		fmt.Println(row)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
