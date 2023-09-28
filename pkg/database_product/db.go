package database_product

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Item struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Unit_cost int    `json:"unit_cost"`
}

var product = []Item{
	{ID: 0, Name: "Product 1", Quantity: 10, Unit_cost: 100},
	{ID: 1, Name: "Product 2", Quantity: 20, Unit_cost: 150},
	{ID: 2, Name: "Product 3", Quantity: 100, Unit_cost: 10},
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
	GET(int) Item
	GETALL() []Item
	POST(string, int, int) int
	DELETE(int) error
	PUT(int, string, int, int) error
}

type MemoryStorage struct {
	data map[int]Item
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		data: make(map[int]Item),
	}
}

func (s *MemoryStorage) GET(ID int) (Name string, Quantity int, Unit_cost int) {
	for _, p := range product {
		if p.ID == ID {
			return p.Name, p.Quantity, p.Unit_cost
		}
	}
	return
}

func (s *MemoryStorage) GETALL() []Item {
	return product
}

func (s *MemoryStorage) POST(Name string, Quantity int, Unit_cost int) (ID int) {
	var NewProduct Item
	NewProduct.ID = product[len(product)-1].ID + 1
	NewProduct.Name = Name
	NewProduct.Quantity = Quantity
	NewProduct.Unit_cost = Unit_cost
	product = append(product, NewProduct)
	return NewProduct.ID
}

func (s *MemoryStorage) DELETE(ID int) {
	for i, p := range product {
		if p.ID == ID {
			product = append(product[:i], product[i+1:]...)
			break
		}
	}
}

func (s *MemoryStorage) PUT(ID int, Name string, Quantity int, Unit_cost int) {

	for i, p := range product {
		if p.ID == ID {
			product[i].Name = Name
			product[i].Quantity = Quantity
			product[i].Unit_cost = Unit_cost
		}
	}
}

func (s *MemoryStorage) GETLEN() int {
	return len(product)
}

type MemoryPostgreSQL struct {
	data map[int]Item
}

func NewMemoryPostgreSQL() *MemoryPostgreSQL {
	return &MemoryPostgreSQL{
		data: make(map[int]Item),
	}
}

func (s *MemoryPostgreSQL) GET(ID int) (Name string, Quantity int, Unit_cost int) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var p Item

	row := conn.QueryRow(context.Background(), `select "ID", "Name", "Quantity", "Unit_coast" FROM "items" WHERE "ID" = $1`, ID)

	err = row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Unit_cost)
	if err != nil {
		panic(err)
	}

	return p.Name, p.Quantity, p.Unit_cost
}

func (s *MemoryPostgreSQL) GETALL() []Item {
	var product = []Item{}
	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var p Item

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
	}
	return product
}

func (s *MemoryPostgreSQL) POST(Name string, Quantity int, Unit_cost int) (ID int) {

	conn, err := getDBConnection()
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
	return newProduct.ID
}

func (s *MemoryPostgreSQL) DELETE(ID int) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}

	defer closeDBConnection(conn)

	conn.Exec(context.Background(), `delete from "items" where "ID"=$1`, ID)
}

func (s *MemoryPostgreSQL) PUT(ID int, Name string, Quantity int, Unit_cost int) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var changeProduct Item

	changeProduct.Name = Name
	changeProduct.Quantity = Quantity
	changeProduct.Unit_cost = Unit_cost

	conn.Exec(context.Background(), `update "items" set "Name"=$1, "Quantity"=$2, "Unit_coast"=$3 where "ID"=$4`, changeProduct.Name, changeProduct.Quantity, changeProduct.Unit_cost, ID)
}
