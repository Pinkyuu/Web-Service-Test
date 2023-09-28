package database_measure

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Measure struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Открытие базы данных
func getDBConnection() (*pgx.Conn, error) {
	config, err := pgx.ParseConfig("postgres://postgres:123@localhost/Web-Service")
	if err != nil {
		return nil, err
	}
	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// Закрытие базы данных
func closeDBConnection(conn *pgx.Conn) {
	conn.Close(context.Background())
}

// Вывод по ID единицу измерения
func GET(ID int) (Name string) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var p Measure

	row := conn.QueryRow(context.Background(), `select "ID", "Name" FROM "measure" WHERE "ID" = $1`, ID)

	err = row.Scan(&p.ID, &p.Name)
	if err != nil {
		panic(err)
	}

	return p.Name
}

// Вывод всех ед.измерений
func GETALL() []Measure {

	var Units = []Measure{}

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var p Measure

	rows, err := conn.Query(context.Background(), "SELECT * FROM measure")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&p.ID, &p.Name)
		if err != nil {
			panic(err)
		}
		Units = append(Units, p)
	}
	return Units
}

// Добавление единицы измерения
func POST(Name string) (ID int) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var newUnit Measure

	newUnit.Name = Name

	row := conn.QueryRow(context.Background(), `insert into "measure"("Name") values($1) RETURNING "ID"`, newUnit.Name)

	err = row.Scan(&newUnit.ID)
	if err != nil {
		panic(err)
	}
	return newUnit.ID
}

// Удаление единицы измерения
func DELETE(ID int) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}

	defer closeDBConnection(conn)

	conn.Exec(context.Background(), `delete from "measure" where "ID"=$1`, ID)
}

// Изменение единицы измерения
func PUT(ID int, Name string) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var changeUnit Measure

	changeUnit.Name = Name

	conn.Exec(context.Background(), `update "measure" set "Name"=$1 where "ID"=$2`, changeUnit.Name, ID)
}
