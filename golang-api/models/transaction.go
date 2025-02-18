package models

import "time"

// 🔹 Struct untuk Request Pembuatan Invoice
type CreateInvoiceRequest struct {
	ExternalID  string `json:"external_id"`
	Amount      int    `json:"amount"`
	PayerEmail  string `json:"payer_email"`
	Description string `json:"description"`
}

// 🔹 Struct untuk Response Invoice yang Dihasilkan
type CreateInvoiceResponse struct {
	InvoiceID  string `json:"id"`
	InvoiceURL string `json:"invoice_url"`
}

// 🔹 Struct untuk Xendit Callback dari Webhook
type XenditCallback struct {
	InvoiceID string `json:"id"`
	Status    string `json:"status"`
}

// 🔹 Struct untuk Transaksi
type Transaction struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	InvoiceID string    `json:"invoice_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
