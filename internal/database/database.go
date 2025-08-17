// file: internal/database/database.go

package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Connect terhubung ke database MariaDB/MySQL
func Connect() {
	dsn := "root:@tcp(127.0.0.1:3306)/prima?parseTime=true"

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Tidak bisa membuka koneksi database: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Tidak bisa terhubung ke database: %v", err)
	}

	// --- KONFIGURASI CONNECTION POOL ---
	// Atur masa pakai maksimum koneksi (misalnya 3 menit)
	DB.SetConnMaxLifetime(time.Minute * 3)
	// Atur jumlah maksimum koneksi yang terbuka
	DB.SetMaxOpenConns(10)
	// Atur jumlah maksimum koneksi yang diam (tidak terpakai)
	DB.SetMaxIdleConns(10)
	// Atur waktu diam maksimum koneksi. Koneksi yang diam lebih dari 1 menit akan ditutup.
	DB.SetConnMaxIdleTime(time.Minute * 1) // <-- INI ADALAH PERBAIKAN UTAMANYA

	fmt.Println("Berhasil terhubung ke database!")
}
