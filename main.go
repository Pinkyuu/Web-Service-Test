package main

import (
	"context"
	"encoding/json"
	"fmt"
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

var product = []item{}

func CheckValid(check item, flags int) bool { // true - ошибка, false - ошибок нет
	switch flags {
	case 1: // Check valid ID
		if check.ID == 0 {
			return true
		} else {
			return false
		}
	case 2: // Check valid ID и Name
		if check.ID == 0 || len(check.Name) == 0 || check.Quantity == 0 || check.Unit_coast == 0 {
			return true
		} else {
			return false
		}
	case 3:
		if len(check.Name) == 0 || check.Quantity == 0 || check.Unit_coast == 0 { // Проверка для записи
			return true
		} else {
			return false
		}
	default:
		return false
	}
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

	var p item

	for rows.Next() {
		err = rows.Scan(&p.ID, &p.Name, &p.Quantity, &p.Unit_coast)
		if err != nil {
			panic(err)
		}
		product = append(product, p)
	}

	jsonBytes, err := json.Marshal(product) // todo:Проверять, пустой ли Product
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

	newProduct.ID = len(product)

	if CheckValid(newProduct, Body) { // Проверка на пустые поля
		fmt.Fprintf(w, "Invalid parameters!")
		return
	}

	product = append(product, newProduct)

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

	vars := mux.Vars(r)
	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, p := range product {
		if p.ID == number {
			jsonBytes, err := json.Marshal(p) // todo:Проверять, пустой ли Product
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				w.Header().Set("Content-Type", "application/json")
				w.Write(jsonBytes)
				break
			}
		}
	}

}

func PutProductByIndex(w http.ResponseWriter, r *http.Request) { // PUT

	var changeProduct item

	err := json.NewDecoder(r.Body).Decode(&changeProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r) // Извлекаем ID
	number, err := strconv.Atoi(vars["id"])
	if err != nil { // Проверяем на ошибки
		log.Fatal(err)
		return
	}

	changeProduct.ID = number

	if CheckValid(changeProduct, AllField) { // Проверка на пустые поля
		fmt.Fprintf(w, "Invalid parameters!")
		return
	}

	for i, p := range product {
		if p.ID == number {
			product[i] = changeProduct
			return
		}
	}

	fmt.Fprintf(w, "Product id: '%v' not found ", number)
}

func DeleteProductByIndex(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Fatal(err)
		return
	}

	for i, p := range product {
		if p.ID == number {
			product = append(product[:i], product[i+1:]...)
			for a := i; a < len(product); a++ {
				product[a].ID--
			}
			break
		}
	}
}
