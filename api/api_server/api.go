package server

import (
	"Web-Service/api/api_measure"
	"Web-Service/api/api_product"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func ServerRun() {
	r := mux.NewRouter()
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	r.HandleFunc("/product", api_product.PersonHandler).Methods("GET", "POST", "OPTIONS")
	r.HandleFunc("/product/{id}", api_product.PersonHandlerByIndex).Methods("GET", "PUT", "DELETE")
	r.HandleFunc("/measure", api_measure.PersonHandler).Methods("GET", "POST", "OPTIONS") // Настроить динамическую маршрутизацию
	r.HandleFunc("/measure/{id}", api_measure.PersonHandlerByIndex).Methods("GET", "PUT", "DELETE")
	log.Println("Server start listen port 8080!")
	err := http.ListenAndServe("localhost:8080", handlers.CORS(headers, methods, origins)(r))
	if err != nil {
		log.Fatal(err)
	}
}
