package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const OnlyID int = 1
const AllField int = 2

type item struct {
	ID         int
	Name       string
	Quantity   int
	Unit_coast int
}

var product = []item{
	{ID: 1, Name: "Product 1", Quantity: 10, Unit_coast: 100},
	{ID: 2, Name: "Product 2", Quantity: 20, Unit_coast: 150},
	{ID: 3, Name: "Product 3", Quantity: 100, Unit_coast: 10},
}

func CheckValid(check item, flags int) bool { // true - ошибка, false - ошибок нет
	switch flags {
	case 1: // Check valid ID
		if check.ID == 0 {
			return true
		} else {
			return false
		}
	case 2: // Check valid ID и Name
		if len(check.Name) == 0 || check.Quantity == 0 || check.Unit_coast == 0 {
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

func personHandler(w http.ResponseWriter, r *http.Request) {
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

	if CheckValid(newProduct, AllField) { // Проверка на пустые поля
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

func personHandlerByIndex(w http.ResponseWriter, r *http.Request) {
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
	var productIndex int = -1

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
	fmt.Fprintf(w, "Тут id: '%v'", number)

	if CheckValid(changeProduct, AllField) { // Проверка на пустые поля
		fmt.Fprintf(w, "Invalid parameters!")
		return
	}

	for i, p := range product {
		if p.ID == number {
			productIndex = i
			break
		}
	}

	if productIndex < 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Product not found! ")
		return
	}

	jsonBytes, err := json.Marshal(changeProduct) // todo:Проверять, пустой ли Product
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)

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
			break
		}
	}
}
