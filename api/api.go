package server

import (
	database_measure "Web-Service/pkg/database_measurement"
	database_product "Web-Service/pkg/database_product"
	valid "Web-Service/pkg/function_check_valid"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Item struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Quantity  int    `json:"quantity"`
	Unit_cost int    `json:"unit_cost"`
}

type Measure struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func ServerRun() {
	r := mux.NewRouter()
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	r.HandleFunc("/product", personHandler).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/product/{id}", personHandlerByIndex).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/measure", measureHandler).Methods("GET", "POST", "OPTIONS") // Настроить динамическую маршрутизацию
	r.HandleFunc("/measure/{id}", measureHandlerByIndex).Methods("GET", "PUT", "DELETE")
	log.Println("Server start listen port 8080!")
	err := http.ListenAndServe("localhost:8080", handlers.CORS(headers, methods, origins)(r))
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

	var newProduct Item

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

	var prod Item

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

	var changeProduct Item

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

func measureHandler(w http.ResponseWriter, r *http.Request) { // switch GET, POST
	switch r.Method {
	case http.MethodGet:
		getMeasureAll(w, r)
	case http.MethodPost:
		postMeasure(w, r)
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}

func getMeasureAll(w http.ResponseWriter, r *http.Request) { // GET - получить список всех продуктов

	Units := database_measure.GETALL()
	jsonBytes, err := json.Marshal(Units)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func postMeasure(w http.ResponseWriter, r *http.Request) { // POST - создать новую запись о продукте

	var newUnits Measure

	err := json.NewDecoder(r.Body).Decode(&newUnits)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if valid.CheckName(newUnits.Name) { // Проверка на пустые поля
		fmt.Fprintf(w, "Invalid parameters!")
		return
	}

	var ID int = database_measure.POST(newUnits.Name)

	jsonBytes, err := json.Marshal(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func measureHandlerByIndex(w http.ResponseWriter, r *http.Request) { // switch GET, PUT, DELETE
	switch r.Method {
	case http.MethodGet:
		getMeasureByIndex(w, r)
	case http.MethodPut:
		PutMeasureByIndex(w, r)
	case http.MethodDelete:
		DeleteMeasureByIndex(w, r)
	default:
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
	}
}

func getMeasureByIndex(w http.ResponseWriter, r *http.Request) { // GET - Вывод продукта с индентификатором i

	var Unit Measure

	vars := mux.Vars(r)
	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if valid.CheckID(number) {
		fmt.Fprintf(w, "Not correct ID: '%v'", number)
	}

	Unit.ID = number
	Unit.Name = database_measure.GET(Unit.ID)

	jsonBytes, err := json.Marshal(Unit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBytes)

	}

}

func PutMeasureByIndex(w http.ResponseWriter, r *http.Request) { // PUT

	var changeUnit Measure

	err := json.NewDecoder(r.Body).Decode(&changeUnit)
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

	changeUnit.ID = number

	if valid.CheckName(changeUnit.Name) { // Проверка на пустые поля
		fmt.Fprintf(w, "Invalid parameters!")
		return
	}

	database_measure.PUT(changeUnit.ID, changeUnit.Name)
}

func DeleteMeasureByIndex(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	database_measure.DELETE(number)
}
