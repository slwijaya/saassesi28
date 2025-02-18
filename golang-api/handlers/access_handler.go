package handlers

import (
	"encoding/json"
	"net/http"
)

// GetProductAccessHandler untuk mendapatkan daftar paket akses produk
func GetProductAccessHandler(w http.ResponseWriter, r *http.Request) {
	packages := []map[string]interface{}{
		{"name": "Basic", "features": []string{"search-products"}},
		{"name": "Pro", "features": []string{"search-products", "top-products"}},
		{"name": "Enterprise", "features": []string{"search-products", "top-products", "recommend-products"}},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(packages)
}
