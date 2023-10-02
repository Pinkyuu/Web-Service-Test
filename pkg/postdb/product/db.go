package postdb_product

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type item struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Unit_cost int     `json:"unit_cost"`
	Measure   measure `json:"measure"`
}

type measure struct {
	ID    int    `json:"id"`
	Value string `json:"name"`
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
	POST(string, int, int, int) int
	DELETE(int) error
	PUT(int, string, int, int, int) error
}

type MemoryPostgreSQL struct {
	data map[int]item
}

func NewMemoryPostgreSQL() *MemoryPostgreSQL {
	return &MemoryPostgreSQL{
		data: make(map[int]item),
	}
}

func (s *MemoryPostgreSQL) GET(ID int) (Name string, Quantity int, Unit_cost int, Measure_ID int, MeasureIDName string) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var p item
	row := conn.QueryRow(context.Background(), `SELECT items.id, items.name, items.quantity, items.unit_coast, items.measure_id, measure.value 
	FROM items 
	JOIN measure ON items.measure_id = measure.id 
	WHERE items.id = $1`, ID)
	err = row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Unit_cost, &p.Measure.ID, &p.Measure.Value)
	if err != nil {
		panic(err)
	}

	return p.Name, p.Quantity, p.Unit_cost, p.Measure.ID, p.Measure.Value
}

func (s *MemoryPostgreSQL) GetAll() []item {
	var product = []item{}
	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var p item

	rows, err := conn.Query(context.Background(), `SELECT items.id, items.name, items.quantity, items.unit_coast, items.measure_id, measure.value     
	FROM items 
	JOIN measure ON items.measure_id = measure.id `)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Unit_cost, &p.Measure.ID, &p.Measure.Value)
		if err != nil {
			panic(err)
		}
		product = append(product, p)
	}
	return product
}

func (s *MemoryPostgreSQL) Post(Name string, Quantity int, Unit_cost int, Measure_ID int) (ID int) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	var newProduct item

	newProduct.Name = Name
	newProduct.Quantity = Quantity
	newProduct.Unit_cost = Unit_cost
	newProduct.Measure.ID = Measure_ID

	row := conn.QueryRow(context.Background(), `insert into "items"(name, quantity, unit_coast, measure_id) values($1, $2, $3, $4) RETURNING "id"`, newProduct.Name, newProduct.Quantity, newProduct.Unit_cost, newProduct.Measure.ID)

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

	_, err = conn.Exec(context.Background(), `delete from items where id=$1`, ID)
	if err != nil {
		panic(err)
	}

}

func (s *MemoryPostgreSQL) Put(ID int, Name string, Quantity int, Unit_cost int, Measure_ID int) {

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	_, err = conn.Exec(context.Background(), `update items set name=$1, quantity=$2, unit_coast=$3, measure_id=$4 where id=$5`, Name, Quantity, Unit_cost, Measure_ID, ID)
	if err != nil {
		panic(err)
	}
}
