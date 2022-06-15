package test

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func GetConnection() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/test?parseDate=true")
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(10 * time.Minute)
	return db
}

func TestExec(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customer(id, name) VALUES('dyaksa', 'Dyaksa')"

	_, err := db.ExecContext(ctx, script)
	if err != nil {
		panic(err)
	}
	fmt.Println("Executed")
}

func TestQuery(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT * FROM customer"

	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Println(id, name)
	}
	defer rows.Close()
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin"
	password := "admin"

	script := "SELECT username, password FROM customer WHERE username='" + username + "' AND password= '" + password + "' LIMIT 1 "
	rows, err := db.QueryContext(ctx, script)

	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Username: ", username)
	} else {
		fmt.Println("No user found")
	}
	defer rows.Close()
}

func TestSQLWithParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "admin"
	password := "admin"

	script := "SELECT username, password FROM customer WHERE username = ? AND password = ? LIMIT 1"

	rows, err := db.QueryContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Username: ", username)
	} else {
		fmt.Println("Data Not Found")
	}

	defer rows.Close()
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "eko"
	password := "mas"

	script := "INSERT INTO customer(username, password) VALUES(? ,?)"
	result, err := db.ExecContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}
	fmt.Println("Last Insert ID: ", lastId)
}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "INSERT INTO customer(email, password) VALUES(? , ?)"
	stmt, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for i := 0; i < 10; i++ {
		email := "dyaksa" + strconv.Itoa(i) + "@gmail.com"
		password := "random"

		result, err := stmt.ExecContext(ctx, email, password)
		if err != nil {
			panic(err)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("last id : ", lastId)

	}
}

func TestTransactionDb(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	script := "INSERT INTO customer(email, password) VALUES(? , ?)"

	for i := 0; i < 10; i++ {
		email := "dyaksaa" + strconv.Itoa(i) + "@gmail.com"
		password := "random"
		result, err := db.ExecContext(ctx, script, email, password)
		if err != nil {
			panic(err)
		}

		lastId, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}
		fmt.Println("last id :", lastId)
	}

	tx.Commit()
}
