package api_product

import (
	"Web-Service/pkg/database_product"
	valid "Web-Service/pkg/function_check_valid"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type item struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Unit_cost int    `json:"unit_cost"`
}

func PersonHandler(w http.ResponseWriter, r *http.Request) { // switch GET, POST
	switch r.Method {
	case http.MethodGet:
		getProductAll(w, r)
	case http.MethodPost:
		postProduct(w, r)
	case http.MethodOptions:
		optionalproduct(w, r)
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}

func getProductAll(w http.ResponseWriter, r *http.Request) { // GET - получить список всех продуктов

	storage := database_product.NewMemoryPostgreSQL()
	product := storage.GETALL()
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

	if valid.CheckBody(newProduct.Name, newProduct.Quantity, newProduct.Unit_cost) { // Проверка на пустые поля
		fmt.Fprintf(w, "Invalid parameters!")
		return
	}

	storage := database_product.NewMemoryPostgreSQL()
	var ID int = storage.POST(newProduct.Name, newProduct.Quantity, newProduct.Unit_cost)

	jsonBytes, err := json.Marshal(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func optionalproduct(w http.ResponseWriter, r *http.Request) { // GET - получить список всех продуктов
	w.WriteHeader(http.StatusOK)
}

func PersonHandlerByIndex(w http.ResponseWriter, r *http.Request) { // switch GET, PUT, DELETE
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

	var prod item

	vars := mux.Vars(r)
	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if valid.CheckID(number) {
		fmt.Fprintf(w, "Not correct ID: '%v'", number)
	}

	storage := database_product.NewMemoryPostgreSQL()
	prod.ID = number
	prod.Name, prod.Quantity, prod.Unit_cost = storage.GET(number)

	jsonBytes, err := json.Marshal(prod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	changeProduct.ID = number

	if valid.CheckBody(changeProduct.Name, changeProduct.Quantity, changeProduct.Unit_cost) { // Проверка на пустые поля
		fmt.Fprintf(w, "Invalid parameters!")
		return
	}

	storage := database_product.NewMemoryPostgreSQL()

	storage.PUT(changeProduct.ID, changeProduct.Name, changeProduct.Quantity, changeProduct.Unit_cost)
}

func DeleteProductByIndex(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	storage := database_product.NewMemoryPostgreSQL()

	storage.DELETE(number)
}
