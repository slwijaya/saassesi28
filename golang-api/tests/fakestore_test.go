package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Struct untuk menyimpan response produk
type FakeStoreProduct struct {
	ID    int     `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

// 🔹 Test untuk FakeStoreAPI menggunakan mock server
func TestMockFakeStoreAPI(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Simulasi HTTP 200
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"id":1, "title":"Mock Product", "price":99.99}]`)) // Simulasi response FakeStoreAPI
	}))
	defer mockServer.Close()

	// Simulasi request ke Mock Server
	resp, err := http.Get(mockServer.URL)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Validasi response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var products []FakeStoreProduct
	err = json.NewDecoder(resp.Body).Decode(&products)
	assert.NoError(t, err)

	// Cek apakah produk yang dikembalikan sesuai
	assert.Len(t, products, 1)
	assert.Equal(t, "Mock Product", products[0].Title)
	assert.Equal(t, 99.99, products[0].Price)
}
