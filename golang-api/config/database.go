// package config

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"

// 	_ "github.com/lib/pq"
// )

// var DB *sql.DB

// func InitDB() {
// 	var err error
// 	connStr := "postgres://postgres:postgres@localhost:5432/saas_platform?sslmode=disable"
// 	DB, err = sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("Failed to connect to database:", err)
// 	}

// 	err = DB.Ping()
// 	if err != nil {
// 		log.Fatal("Database is not reachable:", err)
// 	}

// 	fmt.Println("Connected to PostgreSQL database successfully!")
// }

package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// Inisialisasi koneksi database
func InitDB() {
	var err error
	// connStr := "postgres://postgres:postgres@localhost:5432/saas_platform?sslmode=disable"
	// GANTI localhost MENJADI db (nama service di docker-compose.yml)
	// connStr := "postgres://postgres:postgres@db:5432/saas_platform?sslmode=disable"

	// GANTI "localhost" menjadi Public IP dari Cloud SQL
	connStr := "postgres://postgres:postgres@34.128.100.213:5432/saas_platform?sslmode=disable"

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Failed to connect to database:", err)
	}

	// Cek koneksi database
	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Database is not reachable:", err)
	}

	fmt.Println("✅ Connected to PostgreSQL database successfully!")
}
