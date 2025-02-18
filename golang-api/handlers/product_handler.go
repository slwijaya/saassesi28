package handlers

import (
	"database/sql"
	"encoding/json"
	"golang-api/config"
	"golang-api/models"
	"io"
	"log"
	"net/http"
)

const FAKESTORE_API_URL = "https://fakestoreapi.com/products"

// 🔹 SyncProductsHandler untuk sinkronisasi produk dari API FakeStore
func SyncProductsHandler(w http.ResponseWriter, r *http.Request) {
	response, err := http.Get(FAKESTORE_API_URL)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, "Failed to read API response", http.StatusInternalServerError)
		return
	}

	var products []models.FakeStoreProduct
	err = json.Unmarshal(body, &products)
	if err != nil {
		http.Error(w, "Failed to parse API response", http.StatusInternalServerError)
		return
	}

	for _, product := range products {
		var exists bool

		err := config.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM fakestore_products WHERE external_id=$1)", product.ID).Scan(&exists)
		if err != nil {
			log.Printf("Database error when checking existence: %v", err)
			http.Error(w, "Database error", http.StatusInternalServerError)
			return
		}

		if exists {
			_, err = config.DB.Exec(
				"UPDATE fakestore_products SET rating=$1, rating_count=$2 WHERE external_id=$3",
				product.Rating.Rate, product.Rating.Count, product.ID)
			if err != nil {
				log.Printf("Failed to update product: %v", err)
				http.Error(w, "Failed to update product", http.StatusInternalServerError)
				return
			}
		} else {
			_, err = config.DB.Exec(
				"INSERT INTO fakestore_products (external_id, title, price, category, image_url, description, rating, rating_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
				product.ID, product.Title, product.Price, product.Category, product.ImageURL, product.Description, product.Rating.Rate, product.Rating.Count)
			if err != nil {
				log.Printf("Failed to insert product: %v", err)
				http.Error(w, "Failed to save products", http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Products synced successfully"})
}

// // 🔹 SearchProductsHandler untuk mencari produk berdasarkan keyword
// func SearchProductsHandler(w http.ResponseWriter, r *http.Request) {
// 	queryParam := r.URL.Query().Get("q")

// 	var rows *sql.Rows
// 	var err error

// 	if queryParam == "" {
// 		rows, err = config.DB.Query("SELECT id, title, price, category, image_url, description FROM fakestore_products")
// 	} else {
// 		searchQuery := "%" + queryParam + "%"
// 		rows, err = config.DB.Query("SELECT id, title, price, category, image_url, description FROM fakestore_products WHERE title ILIKE $1", searchQuery)
// 	}

// 	if err != nil {
// 		http.Error(w, "Database query error", http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var products []models.FakeStoreProduct
// 	for rows.Next() {
// 		var product models.FakeStoreProduct
// 		if err := rows.Scan(&product.ID, &product.Title, &product.Price, &product.Category, &product.ImageURL, &product.Description); err != nil {
// 			http.Error(w, "Failed to parse database result", http.StatusInternalServerError)
// 			return
// 		}
// 		products = append(products, product)
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(products)
// }

// 🔹 Handler untuk mencari produk
func SearchProductsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("🔍 Memproses request SearchProductsHandler...") // Log awal

	// Pastikan koneksi database tidak nil
	if config.DB == nil {
		log.Println("❌ Database connection is nil")
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}

	queryParam := r.URL.Query().Get("q")

	var rows *sql.Rows
	var err error

	if queryParam == "" {
		log.Println("🔍 Fetching all products...")
		rows, err = config.DB.Query("SELECT id, title, price, category, image_url, description FROM fakestore_products")
	} else {
		log.Println("🔍 Searching for products with keyword:", queryParam)
		searchQuery := "%" + queryParam + "%"
		rows, err = config.DB.Query("SELECT id, title, price, category, image_url, description FROM fakestore_products WHERE title ILIKE $1", searchQuery)
	}

	if err != nil {
		log.Println("❌ Database query error:", err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.FakeStoreProduct
	for rows.Next() {
		var product models.FakeStoreProduct
		if err := rows.Scan(&product.ID, &product.Title, &product.Price, &product.Category, &product.ImageURL, &product.Description); err != nil {
			log.Println("❌ Error scanning row:", err)
			http.Error(w, "Failed to parse database result", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	log.Println("✅ Search completed, returning results.")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// 🔹 TopProductsHandler untuk mendapatkan 10 produk terbaik berdasarkan rating
func TopProductsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT id, title, price, category, image_url, description, rating
		FROM fakestore_products
		ORDER BY rating DESC
		LIMIT 10`)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.FakeStoreProduct
	for rows.Next() {
		var product models.FakeStoreProduct
		if err := rows.Scan(&product.ID, &product.Title, &product.Price, &product.Category, &product.ImageURL, &product.Description, &product.Rating.Rate); err != nil {
			http.Error(w, "Failed to parse database result", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// 🔹 RecommendProductsHandler untuk merekomendasikan produk berdasarkan kategori
func RecommendProductsHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	if category == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}

	rows, err := config.DB.Query(`
		SELECT id, title, price, category, image_url, description, rating
		FROM fakestore_products
		WHERE category=$1
		ORDER BY rating DESC
		LIMIT 5`, category)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []models.FakeStoreProduct
	for rows.Next() {
		var product models.FakeStoreProduct
		if err := rows.Scan(&product.ID, &product.Title, &product.Price, &product.Category, &product.ImageURL, &product.Description, &product.Rating.Rate); err != nil {
			http.Error(w, "Failed to parse database result", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// CreateProductHandler menangani pembuatan produk baru
func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var product models.FakeStoreProduct
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Simpan produk ke database
	query := `INSERT INTO fakestore_products (external_id, title, price, category, image_url, description) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := config.DB.Exec(query, product.ID, product.Title, product.Price, product.Category, product.ImageURL, product.Description)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Product created successfully",
	})
}
