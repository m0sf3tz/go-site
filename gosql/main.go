package main

import "database/sql"
import _ "github.com/go-sql-driver/mysql"
import "fmt"

//import "strconv"

var db *sql.DB //global handle to our handle

func new_day(str string) {

	//str := "INSERT INTO worker(id,name) VALUES (" + strconv.Itoa(id) + ", 'big dick')"
	str_test := `CREATE TABLE IF NOT EXISTS TODAY_BITCH_2_3(
		worker_id INT PRIMARY KEY,
		name      VARCHAR(255))`

	fmt.Println(str_test)
	_, err := db.Exec(str_test)
	fmt.Println(err)
}

func add(id int) {

	//str := "INSERT INTO worker(id,name) VALUES (" + strconv.Itoa(id) + ", 'big dick')"
	str := "SELECT * from worker"

	fmt.Println(db)
	db.Exec(str)
}

func main() {

	db, _ = sql.Open("mysql", "sam:pw@/workers")
	//i	if err != nil {
	//		fmt.Println("Could not connect to MariaDB!")
	//		panic(err)
	//	}
	err := db.Ping()
	if err != nil {
		fmt.Println("Could not ping maria DB")
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	new_day("sam")

	stmtOut, err := db.Prepare("SELECT id FROM worker where id > 20")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer stmtOut.Close()

	rows, err := stmtOut.Query() // execute our select statement

	add(500)

	for rows.Next() {
		var title string
		var id int

		rows.Scan(&id)
		fmt.Println("Title of tutorial is :", title, "and the ID is: ", id)
	}

}
