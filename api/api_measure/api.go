package api_measure

import (
	valid "Web-Service/pkg/function_check_valid"
	postdb_measure "Web-Service/pkg/postdb/measure"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type measure struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func PersonHandler(w http.ResponseWriter, r *http.Request) { // switch GET, POST
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

	Units := postdb_measure.GetAll()
	jsonBytes, err := json.Marshal(Units)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func postMeasure(w http.ResponseWriter, r *http.Request) { // POST - создать новую запись о продукте

	var newUnits measure

	err := json.NewDecoder(r.Body).Decode(&newUnits)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if valid.CheckName(newUnits.Name) { // Проверка на пустые поля
		fmt.Fprintf(w, "Invalid parameters!")
		return
	}

	var ID int = postdb_measure.Post(newUnits.Name)

	jsonBytes, err := json.Marshal(ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonBytes)
}

func PersonHandlerByIndex(w http.ResponseWriter, r *http.Request) { // switch GET, PUT, DELETE
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

	var Unit measure

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
	Unit.Name = postdb_measure.Get(Unit.ID)

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

	var changeUnit measure

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

	postdb_measure.Put(changeUnit.ID, changeUnit.Name)
}

func DeleteMeasureByIndex(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	number, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postdb_measure.Delete(number)
}
