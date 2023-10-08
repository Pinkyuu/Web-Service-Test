package api_product

import (
	valid "Web-Service/pkg/function_check_valid"
	postdb_product "Web-Service/pkg/postdb/product"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type item = postdb_product.Item
type measure = postdb_product.Measure

type Storage interface {
	Get(int) (string, int, int, int, string)
	GetAll() []item
	Post(string, int, int, int) int
	Delete(int)
	Put(int, string, int, int, int)
}

func PersonHandler(w http.ResponseWriter, r *http.Request) { // switch GET, POST
	storage := postdb_product.NewMemoryPostgreSQL()
	switch r.Method {
	case http.MethodGet:
		getProductAll(w, r, storage)
	case http.MethodPost:
		postProduct(w, r, storage)
	case http.MethodOptions:
		optionalproduct(w, r)
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}

func getProductAll(w http.ResponseWriter, r *http.Request, storage Storage) { // GET - получить список всех продуктов

	product := storage.GetAll()
	jsonBytes, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		panic(err)
	}
}

func postProduct(w http.ResponseWriter, r *http.Request, storage Storage) { // POST - создать новую запись о продукте

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

	var ID int = storage.Post(newProduct.Name, newProduct.Quantity, newProduct.Unit_cost, newProduct.Measure.ID)

	jsonBytes, err := json.Marshal(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		panic(err)
	}
}

func optionalproduct(w http.ResponseWriter, r *http.Request) { // GET - получить список всех продуктов
	w.WriteHeader(http.StatusOK)
}

func PersonHandlerByIndex(w http.ResponseWriter, r *http.Request) { // switch GET, PUT, DELETE
	storage := postdb_product.NewMemoryPostgreSQL()
	switch r.Method {
	case http.MethodGet:
		getProductByIndex(w, r, storage)
	case http.MethodPut:
		PutProductByIndex(w, r, storage)
	case http.MethodDelete:
		DeleteProductByIndex(w, r, storage)
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}

func getProductByIndex(w http.ResponseWriter, r *http.Request, storage Storage) { // GET - Вывод продукта с индентификатором i

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

	prod.ID = number
	prod.Name, prod.Quantity, prod.Unit_cost, prod.Measure.ID, prod.Measure.Value = storage.Get(number)

	jsonBytes, err := json.Marshal(prod)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonBytes)
		if err != nil {
			panic(err)
		}

	}

}

func PutProductByIndex(w http.ResponseWriter, r *http.Request, storage Storage) { // PUT

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

	storage.Put(changeProduct.ID, changeProduct.Name, changeProduct.Quantity, changeProduct.Unit_cost, changeProduct.Measure.ID)
}

func DeleteProductByIndex(w http.ResponseWriter, r *http.Request, storage Storage) {

	vars := mux.Vars(r)
	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	storage.Delete(number)
}
