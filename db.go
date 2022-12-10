package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

var (
	db  *sql.DB
	err error
)

type Case struct {
	Id     int
	Text   string
	Date   string
	Rating int
}

var list Case

func connect() error {
	db, err = sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.PgHost, cfg.PgPort, cfg.PgUser, cfg.PgPass, cfg.PgName))
	if err != nil {
		log.Fatal("Error - connect to DB", err)
		return err
	}
	return nil
}

func selectCase() []Case {
	rows, err := db.Query(`Select * FROM "ToDo"`)
	if err != nil {
		fmt.Println("Error - selectCase()")
	}
	employee := Case{}
	employees := []Case{}

	for rows.Next() {
		var id, rating int
		var text, date string
		err = rows.Scan(&id, &text, &date, &rating)
		if err != nil {
			fmt.Println("Error = selectCase() rows.Next()  db.go")
			panic(err)
		}
		employee.Id = id
		employee.Text = text
		employee.Date = date
		employee.Rating = rating
		employees = append(employees, employee)
	}
	return employees
}

func addCase(list Case) error {
	query := `INSERT INTO "ToDo"("Id", "Text", "Date", "Priority") VALUES ($1, $2, $3, $4)`
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	stmt, err := db.PrepareContext(ctx, query)
	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, list.Id, list.Text, list.Date, list.Rating)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	log.Printf("%d products created ", rows)
	return nil
}

func delCase(id int) error {
	query := fmt.Sprintf(`DELETE FROM "ToDo" WHERE "Id" = %d`, id)
	res, err := db.Exec(query)
	if err == nil {
		_, err := res.RowsAffected()
		if err != nil {
			fmt.Println("Error - delCase()")
		}
		return nil
	}
	return err
}

func selectById(ID int, cases string) error {
	qry := fmt.Sprintf(`Select * FROM "ToDo" WHERE "Id" = %d`, ID)
	rows, err := db.Query(qry)
	if err != nil {
		fmt.Println("Error - selectById() Select")
		return err
	}
	var id, rating int
	var text, date string
	for rows.Next() {
		err = rows.Scan(&id, &text, &date, &rating)
		if err != nil {
			fmt.Println("Error = selectCase() rows.Next()  db.go")
			panic(err)
		}
		fmt.Println("Старая задача:", id, text, date, rating)
		query := fmt.Sprintf(`UPDATE "ToDo" SET "Text" = '%s' WHERE "Id" =  %d`, cases, ID)
		_, err = db.Query(query)
		if err != nil {
			fmt.Println("Error - selectById() UPDATE", err)
			return err
		}
		fmt.Println("Обновленная задача:", id, cases, date, rating)
	}
	return nil
}

func changeCase() {}

func incId() int {
	rows, err := db.Query(`SELECT * FROM "ToDo" ORDER BY "Id" DESC LIMIT 1`)
	if err != nil {
		fmt.Println("Error - incId()")
	}
	var id, rating int
	var text, date string
	for rows.Next() {

		err = rows.Scan(&id, &text, &date, &rating)
		if err != nil {
			fmt.Println("Error = selectCase() rows.Next()  db.go")
			panic(err)
		}
	}
	return id
}
