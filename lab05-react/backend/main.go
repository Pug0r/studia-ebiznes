package main

import (
	"encoding/json"
	"net/http"
)

type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type Payment struct {
	Email  string  `json:"email"`
	Amount float64 `json:"amount"`
}

var products = []Product{
	{ID: 1, Name: "Laptop", Price: 3999.0},
	{ID: 2, Name: "Mysz", Price: 149.0},
	{ID: 3, Name: "Klawiatura", Price: 299.0},
}

var payments = make([]Payment, 0)

func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next(w, r)
	}
}

func productsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func paymentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items := make([]Payment, len(payments))
		copy(items, payments)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	case http.MethodPost:
		var p Payment
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid payload", http.StatusBadRequest)
			return
		}
		payments = append(payments, p)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	http.HandleFunc("/products", withCORS(productsHandler))
	http.HandleFunc("/payments", withCORS(paymentsHandler))
	http.ListenAndServe(":8080", nil)
}
