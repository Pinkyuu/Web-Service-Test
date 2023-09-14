package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"
)

func getDBConnection() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:123@localhost/Web-Service")
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func closeDBConnection(conn *pgx.Conn) {
	conn.Close(context.Background())
}

const (
	OnlyID   int = 1
	AllField int = 2
	Body     int = 3
)

type item struct {
	ID         int    `json:"id"`
	Name       string `json:"Name"`
	Quantity   int    `json:"Quantity"`
	Unit_coast int    `json:"Unit_coast"`
}

func main() {

	r := mux.NewRouter()
	r.HandleFunc("/product", personHandler).Methods("GET", "POST")
	r.HandleFunc("/product/{id}", personHandlerByIndex).Methods("GET", "PUT", "DELETE")
	log.Println("Server start listen port 8080!")
	err := http.ListenAndServe("localhost:8080", r)
	if err != nil {
		log.Fatal(err)
	}
}

func personHandler(w http.ResponseWriter, r *http.Request) { // switch GET, POST
	switch r.Method {
	case http.MethodGet:
		getProductAll(w, r)
	case http.MethodPost:
		postProduct(w, r)
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}

func getProductAll(w http.ResponseWriter, r *http.Request) { // GET - получить список всех продуктов

	var product []item
	var p item

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	rows, err := conn.Query(context.Background(), "SELECT * FROM items")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Unit_coast)
		if err != nil {
			panic(err)
		}
		product = append(product, p)
	}

	jsonBytes, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func postProduct(w http.ResponseWriter, r *http.Request) { // POST - создать новую запись о продукте
	var newProduct item
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	row := conn.QueryRow(context.Background(), `insert into "items"("Name", "Quantity", "Unit_coast") values($1, $2, $3) RETURNING "ID"`, newProduct.Name, newProduct.Quantity, newProduct.Unit_coast)
	if err != nil {
		panic(err)
	}
	err = row.Scan(&newProduct.ID)
	if err != nil {
		panic(err)
	}

	jsonBytes, err := json.Marshal(newProduct.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func personHandlerByIndex(w http.ResponseWriter, r *http.Request) { // switch GET, PUT, DELETE
	switch r.Method {
	case http.MethodGet:
		getProductByIndex(w, r)
	case http.MethodPut:
		PutProductByIndex(w, r)
	case http.MethodDelete:
		DeleteProductByIndex(w, r)
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}

func getProductByIndex(w http.ResponseWriter, r *http.Request) { // GET - Вывод продукта с индентификатором i

	var p item
	vars := mux.Vars(r)
	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
		return
	}

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	row := conn.QueryRow(context.Background(), `select "ID", "Name", "Quantity", "Unit_coast" FROM "items" WHERE "ID" = $1`, number)
	if err != nil {
		panic(err)
	}

	err = row.Scan(&p.ID, &p.Name, &p.Quantity, &p.Unit_coast)
	if err != nil {
		panic(err)
	}

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func PutProductByIndex(w http.ResponseWriter, r *http.Request) { // PUT

	var changeProduct item

	err := json.NewDecoder(r.Body).Decode(&changeProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r) // Извлекаем ID
	id, err := strconv.Atoi(vars["id"])
	if err != nil { // Проверяем на ошибки
		log.Fatal(err)
		return
	}

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	conn.Exec(context.Background(), `update "items" set "Name"=$1, "Quantity"=$2, "Unit_coast"=$3 where "ID"=$4`, changeProduct.Name, changeProduct.Quantity, changeProduct.Unit_coast, id)
}

func DeleteProductByIndex(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
		return
	}

	conn, err := getDBConnection()
	if err != nil {
		panic(err)
	}
	defer closeDBConnection(conn)

	conn.Exec(context.Background(), `delete from "items" where "ID"=$1`, number)

}
