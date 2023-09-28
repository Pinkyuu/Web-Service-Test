package database_measure

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Measure struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

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

func closeDBConnection(conn *pgx.Conn) {
	conn.Close(context.Background())
}

func GET(ID int) (Name string) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var p Measure

	row := conn.QueryRow(context.Background(), `select "ID", "Name" FROM "items" WHERE "ID" = $1`, ID)

	err = row.Scan(&p.ID, &p.Name)
	if err != nil {
		panic(err)
	}

	return p.Name
}

func GETALL() []Measure {

	var Units = []Measure{}
	/*conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var p Units

	rows, err := conn.Query(context.Background(), "SELECT * FROM items")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Unit_cost)
		if err != nil {
			panic(err)
		}
		product = append(product, p)
	}*/
	return Units
}

func POST(Name string, Quantity int, Unit_cost int) (ID int) {

	/*conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var newProduct Item

	newProduct.Name = Name
	newProduct.Quantity = Quantity
	newProduct.Unit_cost = Unit_cost

	row := conn.QueryRow(context.Background(), `insert into "items"("Name", "Quantity", "Unit_coast") values($1, $2, $3) RETURNING "ID"`, newProduct.Name, newProduct.Quantity, newProduct.Unit_cost)

	err = row.Scan(&newProduct.ID)
	if err != nil {
		panic(err)
	}
	return newProduct.ID*/
	return 1
}

func DELETE(ID int) {

	/*conn, err := getDBConnection()
		if err != nil {
			panic(err)
		}

		defer closeDBConnection(conn)

		conn.Exec(context.Background(), `delete from "items" where "ID"=$1`, ID)
	}

	func () PUT(ID int, Name string, Quantity int, Unit_cost int) {

		conn, err := getDBConnection()
		if err != nil {
			panic(err)
		}
		defer closeDBConnection(conn)

		var changeProduct Item

		changeProduct.Name = Name
		changeProduct.Quantity = Quantity
		changeProduct.Unit_cost = Unit_cost

		conn.Exec(context.Background(), `update "items" set "Name"=$1, "Quantity"=$2, "Unit_coast"=$3 where "ID"=$4`, changeProduct.Name, changeProduct.Quantity, changeProduct.Unit_cost, ID)*/
}
