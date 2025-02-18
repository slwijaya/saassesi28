package mocks

import (
	"golang-api/models"

	"github.com/stretchr/testify/mock"
)

// Mock untuk Xendit API
type MockXenditAPI struct {
	mock.Mock
}

// 🔹 Mock untuk Create Invoice
func (m *MockXenditAPI) CreateInvoice(request models.CreateInvoiceRequest) models.CreateInvoiceResponse {
	args := m.Called(request)
	return args.Get(0).(models.CreateInvoiceResponse)
}
