package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Struct untuk menyimpan request invoice
type CreateInvoiceRequest struct {
	ExternalID  string `json:"external_id"`
	Amount      int    `json:"amount"`
	PayerEmail  string `json:"payer_email"`
	Description string `json:"description"`
}

// Struct untuk menyimpan response Xendit
type CreateInvoiceResponse struct {
	InvoiceID  string `json:"id"`
	InvoiceURL string `json:"invoice_url"`
}

// 🔹 Test untuk Xendit API menggunakan mock server
func TestMockXenditAPI(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Simulasi HTTP 200
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":"mock-invoice-id","invoice_url":"https://mock-invoice-url.com"}`)) // Simulasi response Xendit
	}))
	defer mockServer.Close()

	// Simulasi request ke Mock API
	requestBody, _ := json.Marshal(CreateInvoiceRequest{
		ExternalID:  "invoice-123",
		Amount:      100000,
		PayerEmail:  "test@example.com",
		Description: "Test Invoice",
	})

	resp, err := http.Post(mockServer.URL, "application/json", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Validasi response
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var invoice CreateInvoiceResponse
	err = json.NewDecoder(resp.Body).Decode(&invoice)
	assert.NoError(t, err)

	// Cek apakah invoice yang dikembalikan sesuai
	assert.Equal(t, "mock-invoice-id", invoice.InvoiceID)
	assert.Equal(t, "https://mock-invoice-url.com", invoice.InvoiceURL)
}
