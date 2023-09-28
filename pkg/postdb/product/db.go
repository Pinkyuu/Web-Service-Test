package postdb_product

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type item struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Unit_cost int    `json:"unit_cost"`
	Measure   int    `json:"measure`
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

type Storage interface {
	GET(int) item
	GETALL() []item
	POST(string, int, int) int
	DELETE(int) error
	PUT(int, string, int, int) error
}

type MemoryPostgreSQL struct {
	data map[int]item
}

func NewMemoryPostgreSQL() *MemoryPostgreSQL {
	return &MemoryPostgreSQL{
		data: make(map[int]item),
	}
}

func (s *MemoryPostgreSQL) GET(ID int) (Name string, Quantity int, Unit_cost int, Measure int) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var p item

	row := conn.QueryRow(context.Background(), `select "ID", "Name", "Quantity", "Unit_coast", "Measure" FROM "items" WHERE "ID" = $1`, ID)

	err = row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Unit_cost, &p.Measure)
	if err != nil {
		panic(err)
	}

	return p.Name, p.Quantity, p.Unit_cost, p.Measure
}

func (s *MemoryPostgreSQL) GetAll() []item {
	var product = []item{}
	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var p item

	rows, err := conn.Query(context.Background(), "SELECT * FROM items")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Unit_cost, &p.Measure)
		if err != nil {
			panic(err)
		}
		product = append(product, p)
	}
	return product
}

func (s *MemoryPostgreSQL) Post(Name string, Quantity int, Unit_cost int, Measure int) (ID int) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var newProduct item

	newProduct.Name = Name
	newProduct.Quantity = Quantity
	newProduct.Unit_cost = Unit_cost
	newProduct.Measure = Measure

	row := conn.QueryRow(context.Background(), `insert into "items"("Name", "Quantity", "Unit_coast", "Measure") values($1, $2, $3, $4) RETURNING "ID"`, newProduct.Name, newProduct.Quantity, newProduct.Unit_cost, newProduct.Measure)

	err = row.Scan(&newProduct.ID)
	if err != nil {
		panic(err)
	}
	return newProduct.ID
}

func (s *MemoryPostgreSQL) Delete(ID int) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}

	defer closeDBConnection(conn)

	conn.Exec(context.Background(), `delete from "items" where "ID"=$1`, ID)
}

func (s *MemoryPostgreSQL) Put(ID int, Name string, Quantity int, Unit_cost int) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var changeProduct item

	changeProduct.Name = Name
	changeProduct.Quantity = Quantity
	changeProduct.Unit_cost = Unit_cost

	conn.Exec(context.Background(), `update "items" set "Name"=$1, "Quantity"=$2, "Unit_coast"=$3 "Measure"=$4 where "ID"=$5`, changeProduct.Name, changeProduct.Quantity, changeProduct.Unit_cost, changeProduct.Measure, ID)
}
