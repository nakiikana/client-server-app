package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "cow.c"

#cgo CFLAGS: -Wall
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"net/http"
	"tools/models"
	"unsafe"
)

func isStringInMap(m map[string]int, str string) bool {
	_, exists := m[str]
	return exists
}

func buyCandyHandler(w http.ResponseWriter, r *http.Request) {

	var order models.Order
	var cowMes string

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
	}
	if val := models.ValidType[order.CandyType]; order.Money >= (order.CandyCount * val) {
		response := models.ProperResponse{
			Change: order.Money - val*order.CandyCount,
			Thanks: "Thank you!",
		}
		w.WriteHeader(http.StatusCreated)
		message, _ := json.Marshal(response)
		cowMes = string(message)
	} else {
		w.WriteHeader(http.StatusCreated)
		cowMes = fmt.Sprintf("You need %d more money!", val*order.CandyCount-order.Money)
	}

	cowPtr := C.ask_cow(C.CString(cowMes))

	defer C.free(unsafe.Pointer(cowPtr))

	cow := C.GoString(cowPtr)
	json.NewEncoder(w).Encode(cow)
}

func main() {
	PORT := ":8200"
	fmt.Println("Using port number: ", PORT)
	http.HandleFunc("/buy_candy", buyCandyHandler)
	err := http.ListenAndServeTLS(PORT, "../candy_server2/cert.pem", "../candy_server2/key.pem", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}

//из src
// curl -s --key ./candy_client/minica/localhost/key.pem --cert ./candy_client/minica/localhost/cert.pem --cacert ./candy_client/minica/minica.pem -XPOST -H "Content-Type: application/json" -d '{"candyType": "NT", "candyCount": 2, "money": 34}' "https://localhost:8200/buy_candy"
