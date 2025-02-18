package mocks

import (
	"golang-api/models"

	"github.com/stretchr/testify/mock"
)

// MockDatabase struct untuk menggantikan database asli
type MockDatabase struct {
	mock.Mock
}

// 🔹 Mock QueryRow untuk User
func (m *MockDatabase) QueryRow(query string, args ...interface{}) *models.User {
	argsMock := m.Called(query, args)
	return argsMock.Get(0).(*models.User)
}

// 🔹 Mock QueryRow untuk Transaksi
func (m *MockDatabase) QueryRowTransaction(query string, args ...interface{}) *models.Transaction {
	argsMock := m.Called(query, args)
	return argsMock.Get(0).(*models.Transaction)
}

// 🔹 Mock QueryRow untuk Produk
func (m *MockDatabase) QueryRowProduct(query string, args ...interface{}) *models.FakeStoreProduct {
	argsMock := m.Called(query, args)
	return argsMock.Get(0).(*models.FakeStoreProduct)
}
