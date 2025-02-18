package mocks

import (
	"golang-api/models"

	"github.com/stretchr/testify/mock"
)

// Mock untuk FakeStoreAPI
type MockFakeStoreAPI struct {
	mock.Mock
}

// 🔹 Mock untuk mendapatkan produk dari API FakeStore
func (m *MockFakeStoreAPI) GetProducts() []models.FakeStoreProduct {
	args := m.Called()
	return args.Get(0).([]models.FakeStoreProduct)
}
