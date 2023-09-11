package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type item struct {
	ID   int
	Name string
}

var product = []item{
	{ID: 1, Name: "Product 1"},
	{ID: 2, Name: "Product 2"},
	{ID: 3, Name: "Product 3"},
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
		getProduct(w, r)
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

func getProduct(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(product)
	fmt.Fprintf(w, "get product: '%v'", product)
}

func postProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct item
	err := json.NewDecoder(r.Body).Decode(&newProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	product = append(product, newProduct)
	fmt.Fprintf(w, "post new product: '%v'", newProduct)
}

func putProduct(w http.ResponseWriter, r *http.Request) {

	var changeProduct item
	var productIndex int = -1

	err := json.NewDecoder(r.Body).Decode(&changeProduct)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
