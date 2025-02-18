package routes

import (
	"golang-api/handlers"
	"golang-api/middlewares"

	"github.com/gorilla/mux"
)

// SetupRoutes mendaftarkan semua endpoint API
func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Tambahkan Middleware Monitoring ke seluruh request
	router.Use(middlewares.MonitoringMiddleware)

	// 🔹 Auth Routes
	router.HandleFunc("/register", handlers.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")

	// 🔹 Product Routes
	router.HandleFunc("/sync-products", handlers.SyncProductsHandler).Methods("GET")
	router.HandleFunc("/search-products", handlers.SearchProductsHandler).Methods("GET")
	router.HandleFunc("/create-product", handlers.CreateProductHandler).Methods("POST")
	router.HandleFunc("/top-products", handlers.TopProductsHandler).Methods("GET")
	router.HandleFunc("/recommend-products", handlers.RecommendProductsHandler).Methods("GET")

	// 🔹 Transaction Routes
	router.HandleFunc("/create-invoice", handlers.CreateInvoiceHandler).Methods("POST")
	router.HandleFunc("/xendit-callback", handlers.XenditCallbackHandler).Methods("POST")
	router.HandleFunc("/get-transaction/{invoice_id}", handlers.GetTransactionHandler).Methods("GET")
	router.HandleFunc("/transaction-report", handlers.TransactionReportHandler).Methods("GET")

	// 🔹 Access Routes
	router.HandleFunc("/get-product-access", handlers.GetProductAccessHandler).Methods("GET")

	return router
}
