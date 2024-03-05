package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tools/models"
)

func isStringInMap(m map[string]int, str string) bool {
	_, exists := m[str]
	return exists
}

func buyCandyHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if order.Money < 0 || order.CandyCount <= 0 || !isStringInMap(models.ValidType, order.CandyType) {
		errResponse := "Invalid input"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResponse)
		return
	}
	if val := models.ValidType[order.CandyType]; order.Money >= (order.CandyCount * val) {
		response := models.ProperResponse{
			Change: order.Money - val*order.CandyCount,
			Thanks: "Thank you!",
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(fmt.Sprintf("You need %d more money!", val*order.CandyCount-order.Money))
	}
}

func main() {
	PORT := ":8001"
	fmt.Println("Using port number: ", PORT)
	http.HandleFunc("/buy_candy", buyCandyHandler)
	err := http.ListenAndServe(PORT, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
