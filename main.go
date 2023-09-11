package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
		if check.ID == 0 || len(check.Name) == 0 {
			return true
		} else {
			return false
		}
	default:
		return false
	}
}

func main() {

	http.HandleFunc("/product", personHandler)
	http.HandleFunc("/health", healthCheckHandler)
	log.Println("Server start listen port 8080!")
	err := http.ListenAndServe("localhost:8080", nil)
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
	case http.MethodPut:
		putProduct(w, r)
	case http.MethodDelete:
		deleteProduct(w, r)
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}

func getProductAll(w http.ResponseWriter, r *http.Request) { // GET
	jsonBytes, err := json.Marshal(product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func postProduct(w http.ResponseWriter, r *http.Request) { // POST
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
	fmt.Fprintf(w, "post new product: '%v'", newProduct)
}

func putProduct(w http.ResponseWriter, r *http.Request) { // PUT

	var changeProduct item
	var productIndex int = -1

	err := json.NewDecoder(r.Body).Decode(&changeProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if CheckValid(changeProduct, AllField) { // Проверка на пустые поля
		fmt.Fprintf(w, "Invalid parameters!")
		return
	}

	for i, p := range product {
		if p.ID == changeProduct.ID {
			productIndex = i
			break
		}
	}

	if productIndex < 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Product not found! ")
		return
	}

	product[productIndex] = changeProduct
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Put product: '%v'", changeProduct)

}

func deleteProduct(w http.ResponseWriter, r *http.Request) {

	var deleteProduct item
	var productIndex int = -1

	err := json.NewDecoder(r.Body).Decode(&deleteProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if CheckValid(deleteProduct, OnlyID) { // Проверка на ID, пустой или нет
		fmt.Fprintf(w, "Invalid parameters!")
		return
	}

	for i, p := range product {
		if p.ID == deleteProduct.ID {
			productIndex = i
			break
		}
	}

	if productIndex < 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Product not found! ")
		return
	}

	fmt.Fprintf(w, "delete product: '%v'", product[productIndex])
	w.WriteHeader(http.StatusOK)
	product = append(product[:productIndex], product[productIndex+1:]...)

}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Web-Service is working in normal mode! ")
}
