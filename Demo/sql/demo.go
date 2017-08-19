package sql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	. "simpleSql"
	"fmt"
)

/*
demo来源于https://github.com/go-sql-driver/mysql/wiki/Examples
将来参数做啦一下修改
 */

const (
	testDBurl="mhj:mhj123@tcp(192.168.5.15:3306)/mhj?charset=utf8"

	TABLE_UserEntity="User_Entity"
	MoUserEntity_Username="Username"
	MoUserEntity_Password="Password"
)


func main() {
	// Open database connection
	db, err := sql.Open("mysql", testDBurl)
	if err != nil {
		panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
	}
	defer db.Close()

	sSql := NewCommand(TABLE_UserEntity).
		Select(ALL).
		Where(
		MoUserEntity_Username, EQ, PARAM, AND,
		MoUserEntity_Password, EQ, PARAM,
	).Args("mhj", "123")

	fmt.Println("GenerateCommand:",sSql.GenerateCommand())
	fmt.Println("args:",sSql.GetArgs())
	// Execute the query
	rows, err := db.Query(sSql.GenerateCommand(),sSql.GetArgs()...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// Make a slice for the values
	values := make([]sql.RawBytes, len(columns))

	// rows.Scan wants '[]interface{}' as an argument, so we must copy the
	// references into such a slice
	// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Fetch rows
	for rows.Next() {
		// get RawBytes from data
		err = rows.Scan(scanArgs...)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Now do something with the data.
		// Here we just print each column as a string.
		var value string
		for i, col := range values {
			// Here we can check if the value is nil (NULL value)
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			fmt.Println(columns[i], ": ", value)
		}
		fmt.Println("-----------------------------------")
	}
	if err = rows.Err(); err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

}