package handlers

import (
	"encoding/json"
	"net/http"

	"golang-api/config"
	"golang-api/models"

	"github.com/gorilla/mux"
)

// 🔹 CreateInvoiceHandler untuk membuat invoice menggunakan Xendit
func CreateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	var request models.CreateInvoiceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Simpan transaksi ke database
	_, err := config.DB.Exec(`
		INSERT INTO transactions (user_id, invoice_id, amount, status, created_at) 
		VALUES ($1, $2, $3, 'PENDING', NOW())`,
		1, request.ExternalID, request.Amount)
	if err != nil {
		http.Error(w, "Failed to save transaction", http.StatusInternalServerError)
		return
	}

	// Response ke frontend
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.CreateInvoiceResponse{
		InvoiceID:  request.ExternalID,
		InvoiceURL: "https://xendit.co/invoice/" + request.ExternalID,
	})
}

// 🔹 XenditCallbackHandler untuk menangani callback dari Xendit
func XenditCallbackHandler(w http.ResponseWriter, r *http.Request) {
	var callback models.XenditCallback
	if err := json.NewDecoder(r.Body).Decode(&callback); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var exists bool
	err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM transactions WHERE invoice_id=$1)", callback.InvoiceID).Scan(&exists)
	if err != nil || !exists {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	_, err = config.DB.Exec(`
		UPDATE transactions SET status = $1, updated_at = NOW() 
		WHERE invoice_id = $2`, callback.Status, callback.InvoiceID)

	if err != nil {
		http.Error(w, "Failed to update transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Transaction status updated successfully"})
}

// 🔹 GetTransactionHandler untuk mendapatkan transaksi berdasarkan invoice_id
func GetTransactionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	invoiceID := vars["invoice_id"]
	if invoiceID == "" {
		http.Error(w, "Invoice ID is required", http.StatusBadRequest)
		return
	}

	var transaction models.Transaction
	err := config.DB.QueryRow(`
		SELECT id, user_id, invoice_id, amount, status, created_at, updated_at
		FROM transactions WHERE invoice_id = $1`, invoiceID).
		Scan(&transaction.ID, &transaction.UserID, &transaction.InvoiceID, &transaction.Amount, &transaction.Status, &transaction.CreatedAt, &transaction.UpdatedAt)

	if err != nil {
		http.Error(w, "Transaction not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transaction)
}

// 🔹 TransactionReportHandler untuk mendapatkan laporan transaksi
func TransactionReportHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT id, user_id, invoice_id, amount, status, created_at, updated_at 
		FROM transactions ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(&t.ID, &t.UserID, &t.InvoiceID, &t.Amount, &t.Status, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			http.Error(w, "Failed to parse database result", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}
