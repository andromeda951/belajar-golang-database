package belajar_golang_database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestExecSql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "INSERT INTO customer(id, name) VALUES('joko', 'Joko');"
	_, err := db.ExecContext(ctx, script)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success Insert Data to Database")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	rows, err := db.QueryContext(ctx, "SELECT id, name FROM customer")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Println("id:", id)
		fmt.Println("name:", name)
	}
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	query := "SELECT id, name, email, balance, rating, brith_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float32
		var birthDate sql.NullTime
		var createAt time.Time
		var married bool
		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createAt)
		if err != nil {
			panic(err)
		}

		fmt.Println("=======================")
		fmt.Println("id:", id)
		fmt.Println("name:", name)
		if email.Valid {
			fmt.Println("email:", email.String)
		}
		fmt.Println("balance:", balance)
		fmt.Println("rating:", rating)
		if birthDate.Valid {
			fmt.Println("birth date:", birthDate.Time)
		}
		fmt.Println("married:", married)
		fmt.Println("create at:", createAt)
	}

	defer rows.Close()
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	
	username := "admin" // "admin'; #"
	password := "admin"

	script := "SELECT username FROM user WHERE username = '" + username +"' AND password = '" + password + "' LIMIT 1"
	rows, err := db.QueryContext(ctx, script)
	if err != err {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Sukses login")
	} else {
		fmt.Println("Gagal login")
	}
}

func TestSqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	
	username := "admin'; #"
	password := "admin"

	script := "SELECT username FROM user WHERE username = ? AND password = ? LIMIT 1"
	rows, err := db.QueryContext(ctx, script, username, password)
	if err != err {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var username string
		err := rows.Scan(&username)
		if err != nil {
			panic(err)
		}
		fmt.Println("Sukses login")
	} else {
		fmt.Println("Gagal login")
	}
}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	username := "eko; DROP TABLE user; #"
	password := "eko"
	script := "INSERT INTO user(username, password) VALUES(?, ?);"
	_, err := db.ExecContext(ctx, script, username, password)
	if err != nil {
		panic(err)
	}

	fmt.Println("Success Insert New User")
}

func TestAutoIncrement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "andromeda@gmail.com"
	comment := "test koment"
	sqlQuery := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	result, err := db.ExecContext(ctx, sqlQuery, email, comment)
	if err != nil {
		panic(err)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Last insert id:", insertId)
}

func TestPrepareStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	stmt, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	for i:=0; i<10; i++ {
		email := "andromeda" + strconv.Itoa(i) + "@gmail.com"
		comment := "ini komen ke " + strconv.Itoa(i)
		result, err := stmt.ExecContext(ctx, email, comment)
		if err != nil {
			panic(err)
		}

		lastInsertId, _ := result.LastInsertId()
		fmt.Println("Comment id:", lastInsertId)
	}
}

func TestTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}
	
	// do transaction
	ctx := context.Background()
	sqlQuery := "INSERT INTO comments(email, comment) VALUES(?, ?)"
	for i:=0; i<10; i++ {
		email := "andromeda" + strconv.Itoa(i) + "@gmail.com"
		comment := "ini komen ke " + strconv.Itoa(i)
		result, err := tx.ExecContext(ctx, sqlQuery, email, comment)
		if err != nil {
			panic(err)
		}

		lastInsertId, _ := result.LastInsertId()
		fmt.Println("Comment id:", lastInsertId)
	}

	err = tx.Rollback()
	if err != nil {
		panic(err)
	}
}