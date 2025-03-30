package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	mycache "sray/my-cache"
)

func main() {
	http.HandleFunc("/order", createOrderHandler)
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody struct {
		IdempotencyKey string `json:"idempotency_key"`
		ItemName       string `json:"item_name"`
		Amount         int    `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	orderID, err := createOrGetOrder(requestBody.IdempotencyKey, requestBody.ItemName, requestBody.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"order_id": orderID,
	})
}

func createOrGetOrder(idempotencyKey, itemName string, amount int) (string, error) {
	if orderID := mycache.Get(idempotencyKey); orderID != "" {
		return orderID, fmt.Errorf("duplicate request: already processed")
	}

	orderID := fmt.Sprintf("ORD-%s%d", itemName[:3], amount)
	mycache.Put(idempotencyKey, orderID)
	return orderID, nil
}
