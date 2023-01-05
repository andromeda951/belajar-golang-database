package belajar_golang_database

import (
	"context"
	"fmt"
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
	sql := "SELECT id, name, email, balance, rating, brith_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, sql)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, name, email string
		var balance int32
		var rating float32
		var birthDate, createAt time.Time
		var married bool
		err := rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createAt)
		if err != nil {
			panic(err)
		}

		fmt.Println("=======================")
		fmt.Println("id:", id)
		fmt.Println("name:", name)
		fmt.Println("email:", email)
		fmt.Println("balance:", balance)
		fmt.Println("rating:", rating)
		fmt.Println("birth date:", birthDate)
		fmt.Println("married:", married)
		fmt.Println("create at:", createAt)
	}

	defer rows.Close()

}
