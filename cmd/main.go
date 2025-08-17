package main

import (
	"database/sql"
	"log"
	"net/http"
	"scm-api/internal/database"
	"scm-api/internal/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// =================================================================
// DEFINISI STRUCT UNTUK RESPON API
// =================================================================

type PembelianResponse struct {
	PembelianID  int64           `json:"pembelian_id"`
	SupplierID   int64           `json:"supplier_id"`
	NamaSupplier string          `json:"nama_supplier"`
	TanggalPesan string          `json:"tanggal_pesan"`
	EstimasiTiba sql.NullString  `json:"estimasi_tiba"`
	TotalBiaya   sql.NullFloat64 `json:"total_biaya"`
	Status       string          `json:"status"`
}

type DetailPembelianResponse struct {
	ProdukID        int64   `json:"produk_id"`
	NamaProduk      string  `json:"nama_produk"`
	Jumlah          int     `json:"jumlah"`
	HargaBeliSatuan float64 `json:"harga_beli_satuan"`
	Subtotal        float64 `json:"subtotal"`
}

type PembelianDenganDetailResponse struct {
	PembelianID  int64                     `json:"pembelian_id"`
	SupplierID   int64                     `json:"supplier_id"`
	NamaSupplier string                    `json:"nama_supplier"`
	TanggalPesan string                    `json:"tanggal_pesan"`
	EstimasiTiba sql.NullString            `json:"estimasi_tiba"`
	TotalBiaya   sql.NullFloat64           `json:"total_biaya"`
	Status       string                    `json:"status"`
	Details      []DetailPembelianResponse `json:"details"`
}

// StokResponse adalah struct untuk menampung data gabungan stok, produk, dan gudang
type StokResponse struct {
	StokID        int64  `json:"stok_id"`
	ProdukID      int64  `json:"produk_id"`
	NamaProduk    string `json:"nama_produk"`
	GudangID      int64  `json:"gudang_id"`
	NamaGudang    string `json:"nama_gudang"`
	Jumlah        int    `json:"jumlah"`
	TanggalUpdate string `json:"tanggal_update"`
}

// DashboardStats adalah struct untuk menampung data ringkasan dashboard
type DashboardStats struct {
	JumlahProduk    int `json:"jumlah_produk"`
	JumlahSupplier  int `json:"jumlah_supplier"`
	JumlahPembelian int `json:"jumlah_pembelian"`
	JumlahGudang    int `json:"jumlah_gudang"`
}

// (di bawah struct-struct lainnya)

// StokChartResponse adalah struct untuk data grafik stok per produk
type StokChartResponse struct {
	Labels []string `json:"labels"` // Untuk nama produk
	Data   []int    `json:"data"`   // Untuk jumlah stok
}

// =================================================================
// FUNGSI UTAMA (MAIN)
// =================================================================

func main() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://127.0.0.1:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	database.Connect()

	api := router.Group("/api")
	{
		// --- Rute-rute Produk ---
		api.GET("/produk", getProdukHandler)
		api.GET("/produk/:id", getProdukByIdHandler)
		api.POST("/produk", createProdukHandler)
		api.PUT("/produk/:id", updateProdukHandler)
		api.DELETE("/produk/:id", deleteProdukHandler)

		// --- Rute-rute Supplier ---
		api.GET("/supplier", getSuppliersHandler)
		api.GET("/supplier/:id", getSupplierByIdHandler)
		api.POST("/supplier", createSupplierHandler)
		api.PUT("/supplier/:id", updateSupplierHandler)
		api.DELETE("/supplier/:id", deleteSupplierHandler)

		// --- Rute-rute Pembelian ---
		api.GET("/pembelian", getPembelianHandler)
		api.GET("/pembelian/:id", getPembelianByIdHandler)
		api.POST("/pembelian", createPembelianHandler)
		api.DELETE("/pembelian/:id", deletePembelianHandler)
		api.PUT("/pembelian/:id/terima", terimaPembelianHandler)

		// --- Rute-rute Gudang (BARU) ---
		api.GET("/gudang", getGudangHandler)
		api.GET("/gudang/:id", getGudangByIdHandler)
		api.POST("/gudang", createGudangHandler)
		api.PUT("/gudang/:id", updateGudangHandler)
		api.DELETE("/gudang/:id", deleteGudangHandler)

		// --- Rute-rute Stok (BARU) ---
		api.GET("/stok", getStokHandler)
		api.POST("/stok/adjust", adjustStokHandler)

		// --- Rute-rute Dashboard (BARU) ---
		api.GET("/dashboard/stats", getDashboardStatsHandler)
		api.GET("/dashboard/stok-per-produk", getStokChartHandler)
		api.GET("/dashboard/pembelian-terakhir", getPembelianTerakhirHandler)
	}

	log.Println("Server berjalan di http://localhost:8080")
	router.Run(":8080")
}

// =================================================================
// HANDLER UNTUK MODUL PRODUK
// =================================================================

func getProdukHandler(c *gin.Context) {
	rows, err := database.DB.Query("SELECT produk_id, sku, nama_produk, deskripsi, kategori, satuan, harga_jual, berat_kg, gambar_produk, supplier_id FROM produk")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}
	defer rows.Close()
	daftarProduk := make([]models.Produk, 0)
	for rows.Next() {
		var p models.Produk
		err := rows.Scan(&p.ProdukID, &p.SKU, &p.NamaProduk, &p.Deskripsi, &p.Kategori, &p.Satuan, &p.HargaJual, &p.BeratKg, &p.GambarProduk, &p.SupplierID)
		if err != nil {
			log.Printf("Error scanning row produk: %v", err)
			continue
		}
		daftarProduk = append(daftarProduk, p)
	}
	c.JSON(http.StatusOK, daftarProduk)
}

func getProdukByIdHandler(c *gin.Context) {
	id := c.Param("id")
	var p models.Produk
	query := "SELECT produk_id, sku, nama_produk, deskripsi, kategori, satuan, harga_jual, berat_kg, gambar_produk, supplier_id FROM produk WHERE produk_id = ?"
	row := database.DB.QueryRow(query, id)
	err := row.Scan(&p.ProdukID, &p.SKU, &p.NamaProduk, &p.Deskripsi, &p.Kategori, &p.Satuan, &p.HargaJual, &p.BeratKg, &p.GambarProduk, &p.SupplierID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan internal"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func createProdukHandler(c *gin.Context) {
	var req struct {
		SKU        string   `json:"sku"`
		NamaProduk string   `json:"nama_produk"`
		Deskripsi  *string  `json:"deskripsi"`
		Kategori   *string  `json:"kategori"`
		Satuan     string   `json:"satuan"`
		HargaJual  float64  `json:"harga_jual"`
		BeratKg    *float64 `json:"berat_kg"`
		SupplierID *int64   `json:"supplier_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data JSON tidak valid: " + err.Error()})
		return
	}
	produkBaru := models.Produk{SKU: req.SKU, NamaProduk: req.NamaProduk, Satuan: req.Satuan, HargaJual: req.HargaJual}
	if req.Kategori != nil {
		produkBaru.Kategori = sql.NullString{String: *req.Kategori, Valid: true}
	}
	if req.Deskripsi != nil {
		produkBaru.Deskripsi = sql.NullString{String: *req.Deskripsi, Valid: true}
	}
	if req.BeratKg != nil {
		produkBaru.BeratKg = sql.NullFloat64{Float64: *req.BeratKg, Valid: true}
	}
	if req.SupplierID != nil {
		produkBaru.SupplierID = sql.NullInt64{Int64: *req.SupplierID, Valid: true}
	}
	query := `INSERT INTO produk (sku, nama_produk, deskripsi, kategori, satuan, harga_jual, berat_kg, supplier_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, produkBaru.SKU, produkBaru.NamaProduk, produkBaru.Deskripsi, produkBaru.Kategori, produkBaru.Satuan, produkBaru.HargaJual, produkBaru.BeratKg, produkBaru.SupplierID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan produk ke database"})
		return
	}
	id, _ := result.LastInsertId()
	produkBaru.ProdukID = id
	c.JSON(http.StatusCreated, produkBaru)
}

func updateProdukHandler(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		SKU        string  `json:"sku"`
		NamaProduk string  `json:"nama_produk"`
		Kategori   *string `json:"kategori"`
		Satuan     string  `json:"satuan"`
		HargaJual  float64 `json:"harga_jual"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data JSON tidak valid: " + err.Error()})
		return
	}
	query := `UPDATE produk SET sku = ?, nama_produk = ?, kategori = ?, satuan = ?, harga_jual = ? WHERE produk_id = ?`
	_, err := database.DB.Exec(query, req.SKU, req.NamaProduk, req.Kategori, req.Satuan, req.HargaJual, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate produk"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil diupdate"})
}

func deleteProdukHandler(c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM produk WHERE produk_id = ?`
	_, err := database.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}

// =================================================================
// HANDLER UNTUK MODUL SUPPLIER
// =================================================================

func getSuppliersHandler(c *gin.Context) {
	rows, err := database.DB.Query("SELECT supplier_id, nama_supplier, alamat, kontak, contact_person, rating FROM supplier")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data supplier"})
		return
	}
	defer rows.Close()
	daftarSupplier := make([]models.Supplier, 0)
	for rows.Next() {
		var s models.Supplier
		err := rows.Scan(&s.SupplierID, &s.NamaSupplier, &s.Alamat, &s.Kontak, &s.ContactPerson, &s.Rating)
		if err != nil {
			log.Printf("Error scanning row supplier: %v", err)
			continue
		}
		daftarSupplier = append(daftarSupplier, s)
	}
	c.JSON(http.StatusOK, daftarSupplier)
}

func getSupplierByIdHandler(c *gin.Context) {
	id := c.Param("id")
	var s models.Supplier
	query := "SELECT supplier_id, nama_supplier, alamat, kontak, contact_person, rating FROM supplier WHERE supplier_id = ?"
	row := database.DB.QueryRow(query, id)
	err := row.Scan(&s.SupplierID, &s.NamaSupplier, &s.Alamat, &s.Kontak, &s.ContactPerson, &s.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Supplier tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan internal"})
		return
	}
	c.JSON(http.StatusOK, s)
}

func createSupplierHandler(c *gin.Context) {
	var req struct {
		NamaSupplier  string   `json:"nama_supplier"`
		Alamat        *string  `json:"alamat"`
		Kontak        *string  `json:"kontak"`
		ContactPerson *string  `json:"contact_person"`
		Rating        *float64 `json:"rating"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data JSON tidak valid: " + err.Error()})
		return
	}
	supplierBaru := models.Supplier{NamaSupplier: req.NamaSupplier}
	if req.Alamat != nil {
		supplierBaru.Alamat = sql.NullString{String: *req.Alamat, Valid: true}
	}
	if req.Kontak != nil {
		supplierBaru.Kontak = sql.NullString{String: *req.Kontak, Valid: true}
	}
	if req.ContactPerson != nil {
		supplierBaru.ContactPerson = sql.NullString{String: *req.ContactPerson, Valid: true}
	}
	if req.Rating != nil {
		supplierBaru.Rating = sql.NullFloat64{Float64: *req.Rating, Valid: true}
	}
	query := `INSERT INTO supplier (nama_supplier, alamat, kontak, contact_person, rating) VALUES (?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, supplierBaru.NamaSupplier, supplierBaru.Alamat, supplierBaru.Kontak, supplierBaru.ContactPerson, supplierBaru.Rating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan supplier ke database"})
		return
	}
	id, _ := result.LastInsertId()
	supplierBaru.SupplierID = id
	c.JSON(http.StatusCreated, supplierBaru)
}

func updateSupplierHandler(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		NamaSupplier  string   `json:"nama_supplier"`
		Alamat        *string  `json:"alamat"`
		Kontak        *string  `json:"kontak"`
		ContactPerson *string  `json:"contact_person"`
		Rating        *float64 `json:"rating"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data JSON tidak valid"})
		return
	}
	query := `UPDATE supplier SET nama_supplier = ?, alamat = ?, kontak = ?, contact_person = ?, rating = ? WHERE supplier_id = ?`
	_, err := database.DB.Exec(query, req.NamaSupplier, req.Alamat, req.Kontak, req.ContactPerson, req.Rating, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate supplier"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supplier berhasil diupdate"})
}

func deleteSupplierHandler(c *gin.Context) {
	id := c.Param("id")
	query := `DELETE FROM supplier WHERE supplier_id = ?`
	_, err := database.DB.Exec(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus supplier"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Supplier berhasil dihapus"})
}

// =================================================================
// HANDLER UNTUK MODUL PEMBELIAN
// =================================================================

func getPembelianHandler(c *gin.Context) {
	query := `SELECT p.pembelian_id, p.supplier_id, s.nama_supplier, p.tanggal_pesan, p.estimasi_tiba, p.total_biaya, p.status FROM pembelian p JOIN supplier s ON p.supplier_id = s.supplier_id ORDER BY p.tanggal_pesan DESC`
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pembelian"})
		return
	}
	defer rows.Close()
	daftarPembelian := make([]PembelianResponse, 0)
	for rows.Next() {
		var p PembelianResponse
		err := rows.Scan(&p.PembelianID, &p.SupplierID, &p.NamaSupplier, &p.TanggalPesan, &p.EstimasiTiba, &p.TotalBiaya, &p.Status)
		if err != nil {
			log.Printf("Error scanning row pembelian: %v", err)
			continue
		}
		daftarPembelian = append(daftarPembelian, p)
	}
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan internal"})
		return
	}
	c.JSON(http.StatusOK, daftarPembelian)
}

func createPembelianHandler(c *gin.Context) {
	var req struct {
		SupplierID   int64   `json:"supplier_id"`
		TanggalPesan string  `json:"tanggal_pesan"`
		Status       string  `json:"status"`
		TotalBiaya   float64 `json:"total_biaya"`
		Details      []struct {
			ProdukID        int64   `json:"produk_id"`
			Jumlah          int     `json:"jumlah"`
			HargaBeliSatuan float64 `json:"harga_beli_satuan"`
			Subtotal        float64 `json:"subtotal"`
		} `json:"details"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data JSON tidak valid: " + err.Error()})
		return
	}
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi database"})
		return
	}
	queryHeader := `INSERT INTO pembelian (supplier_id, tanggal_pesan, total_biaya, status) VALUES (?, ?, ?, ?)`
	result, err := tx.Exec(queryHeader, req.SupplierID, req.TanggalPesan, req.TotalBiaya, "Dipesan")
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data pembelian"})
		return
	}
	pembelianID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mendapatkan ID pembelian"})
		return
	}
	queryDetail := `INSERT INTO detail_pembelian (pembelian_id, produk_id, jumlah, harga_beli_satuan, subtotal) VALUES (?, ?, ?, ?, ?)`
	for _, detail := range req.Details {
		_, err := tx.Exec(queryDetail, pembelianID, detail.ProdukID, detail.Jumlah, detail.HargaBeliSatuan, detail.Subtotal)
		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan detail produk pembelian"})
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan transaksi"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Pesanan pembelian berhasil dibuat", "pembelian_id": pembelianID})
}

func getPembelianByIdHandler(c *gin.Context) {
	id := c.Param("id")
	var response PembelianDenganDetailResponse
	queryHeader := `SELECT p.pembelian_id, p.supplier_id, s.nama_supplier, p.tanggal_pesan, p.estimasi_tiba, p.total_biaya, p.status FROM pembelian p JOIN supplier s ON p.supplier_id = s.supplier_id WHERE p.pembelian_id = ?`
	row := database.DB.QueryRow(queryHeader, id)
	err := row.Scan(&response.PembelianID, &response.SupplierID, &response.NamaSupplier, &response.TanggalPesan, &response.EstimasiTiba, &response.TotalBiaya, &response.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan pembelian tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data header pembelian"})
		return
	}
	queryDetail := `SELECT d.produk_id, pr.nama_produk, d.jumlah, d.harga_beli_satuan, d.subtotal FROM detail_pembelian d JOIN produk pr ON d.produk_id = pr.produk_id WHERE d.pembelian_id = ?`
	rows, err := database.DB.Query(queryDetail, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail produk pembelian"})
		return
	}
	defer rows.Close()
	details := make([]DetailPembelianResponse, 0)
	for rows.Next() {
		var d DetailPembelianResponse
		if err := rows.Scan(&d.ProdukID, &d.NamaProduk, &d.Jumlah, &d.HargaBeliSatuan, &d.Subtotal); err != nil {
			log.Printf("Gagal scan detail pembelian: %v", err)
			continue
		}
		details = append(details, d)
	}
	response.Details = details
	c.JSON(http.StatusOK, response)
}

// HANDLER UNTUK DELETE PEMBELIAN
// ===============================
func deletePembelianHandler(c *gin.Context) {
	id := c.Param("id")

	// Mulai transaksi
	tx, err := database.DB.Begin()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi"})
		return
	}

	// Hapus dulu semua baris di tabel detail yang terkait
	_, err = tx.Exec("DELETE FROM detail_pembelian WHERE pembelian_id = ?", id)
	if err != nil {
		tx.Rollback() // Batalkan jika gagal
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus detail pembelian"})
		return
	}

	// Setelah itu, baru hapus baris di tabel header
	_, err = tx.Exec("DELETE FROM pembelian WHERE pembelian_id = ?", id)
	if err != nil {
		tx.Rollback() // Batalkan jika gagal
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus pembelian"})
		return
	}

	// Jika semua berhasil, simpan perubahan
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pesanan pembelian berhasil dihapus"})
}

// =================================================================
// HANDLER UNTUK MODUL GUDANG
// =================================================================

func getGudangHandler(c *gin.Context) {
	rows, err := database.DB.Query("SELECT gudang_id, nama_gudang, lokasi FROM gudang")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data gudang"})
		return
	}
	defer rows.Close()
	daftarGudang := make([]models.Gudang, 0)
	for rows.Next() {
		var g models.Gudang
		if err := rows.Scan(&g.GudangID, &g.NamaGudang, &g.Lokasi); err != nil {
			log.Printf("Error scanning row gudang: %v", err)
			continue
		}
		daftarGudang = append(daftarGudang, g)
	}
	c.JSON(http.StatusOK, daftarGudang)
}

func getGudangByIdHandler(c *gin.Context) {
	id := c.Param("id")
	var g models.Gudang
	row := database.DB.QueryRow("SELECT gudang_id, nama_gudang, lokasi FROM gudang WHERE gudang_id = ?", id)
	err := row.Scan(&g.GudangID, &g.NamaGudang, &g.Lokasi)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Gudang tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Terjadi kesalahan internal"})
		return
	}
	c.JSON(http.StatusOK, g)
}

func createGudangHandler(c *gin.Context) {
	var g models.Gudang
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data JSON tidak valid"})
		return
	} /*  */
	query := "INSERT INTO gudang (nama_gudang, lokasi) VALUES (?, ?)"
	result, err := database.DB.Exec(query, g.NamaGudang, g.Lokasi)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan gudang"})
		return
	}
	id, _ := result.LastInsertId()
	g.GudangID = id
	c.JSON(http.StatusCreated, g)
}

func updateGudangHandler(c *gin.Context) {
	id := c.Param("id")
	var g models.Gudang
	if err := c.ShouldBindJSON(&g); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data JSON tidak valid"})
		return
	}
	query := "UPDATE gudang SET nama_gudang = ?, lokasi = ? WHERE gudang_id = ?"
	_, err := database.DB.Exec(query, g.NamaGudang, g.Lokasi, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate gudang"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Gudang berhasil diupdate"})
}

func deleteGudangHandler(c *gin.Context) {
	id := c.Param("id")
	_, err := database.DB.Exec("DELETE FROM gudang WHERE gudang_id = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus gudang"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Gudang berhasil dihapus"})
}

// =================================================================
// HANDLER UNTUK MODUL STOK
// =================================================================

func getStokHandler(c *gin.Context) {
	query := `
        SELECT 
            s.stok_id, s.produk_id, p.nama_produk, 
            s.gudang_id, g.nama_gudang, s.jumlah, s.tanggal_update
        FROM stok s
        JOIN produk p ON s.produk_id = p.produk_id
        JOIN gudang g ON s.gudang_id = g.gudang_id
    `
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data stok"})
		return
	}
	defer rows.Close()

	daftarStok := make([]StokResponse, 0)
	for rows.Next() {
		var s StokResponse
		err := rows.Scan(
			&s.StokID, &s.ProdukID, &s.NamaProduk,
			&s.GudangID, &s.NamaGudang, &s.Jumlah, &s.TanggalUpdate,
		)
		if err != nil {
			log.Printf("Error scanning row stok: %v", err)
			continue
		}
		daftarStok = append(daftarStok, s)
	}
	c.JSON(http.StatusOK, daftarStok)
}

// HANDLER UNTUK PENYESUAIAN STOK (UPSERT)
// =======================================
func adjustStokHandler(c *gin.Context) {
	var req struct {
		ProdukID int64 `json:"produk_id"`
		GudangID int64 `json:"gudang_id"`
		Jumlah   int   `json:"jumlah"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data JSON tidak valid"})
		return
	}

	// Query UPSERT: Insert data baru, tapi jika terjadi duplikasi pada unique key
	// (produk_id, gudang_id), maka update kolom jumlah.
	query := `
        INSERT INTO stok (produk_id, gudang_id, jumlah, tanggal_update)
        VALUES (?, ?, ?, NOW())
        ON DUPLICATE KEY UPDATE jumlah = VALUES(jumlah), tanggal_update = NOW()
    `

	_, err := database.DB.Exec(query, req.ProdukID, req.GudangID, req.Jumlah)
	if err != nil {
		log.Printf("Error upsert stok: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyesuaikan stok"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stok berhasil disesuaikan"})
}

// HANDLER UNTUK MENERIMA PESANAN PEMBELIAN & UPDATE STOK
// ======================================================
func terimaPembelianHandler(c *gin.Context) {
	pembelianID := c.Param("id")

	// Mulai Transaksi
	tx, err := database.DB.Begin()
	if err := database.DB.Ping(); err != nil {
		log.Printf("Koneksi database terputus sebelum transaksi: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Koneksi ke database terputus, coba lagi.", "detail": err.Error()})
		return
	}

	// 1. Ambil semua item dari detail_pembelian untuk pesanan ini
	rows, err := tx.Query("SELECT produk_id, jumlah, gudang_id FROM detail_pembelian WHERE pembelian_id = ?", pembelianID)
	// Asumsi sementara kita terima barang di gudang_id = 1. Nanti ini bisa dibuat lebih dinamis.
	// Untuk amannya, kita modifikasi query di atas jika tabel detail_pembelian belum punya gudang_id
	// Mari kita asumsikan barang masuk ke Gudang ID 1
	rows, err = tx.Query("SELECT produk_id, jumlah FROM detail_pembelian WHERE pembelian_id = ?", pembelianID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail pesanan"})
		return
	}
	defer rows.Close()

	// 2. Looping setiap item untuk mengupdate (UPSERT) tabel stok
	for rows.Next() {
		var detail struct {
			ProdukID int64
			Jumlah   int
		}
		if err := rows.Scan(&detail.ProdukID, &detail.Jumlah); err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal scan detail item"})
			return
		}

		// Query UPSERT untuk menambah stok. Jika sudah ada, tambahkan jumlahnya.
		// Asumsi barang masuk ke Gudang ID 1
		queryStok := `
            INSERT INTO stok (produk_id, gudang_id, jumlah, tanggal_update)
            VALUES (?, 1, ?, NOW())
            ON DUPLICATE KEY UPDATE jumlah = jumlah + VALUES(jumlah), tanggal_update = NOW()
        `
		_, err := tx.Exec(queryStok, detail.ProdukID, detail.Jumlah)
		if err != nil {
			tx.Rollback()
			log.Printf("Gagal upsert stok untuk produk ID %d: %v", detail.ProdukID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate stok produk", "detail": err.Error()})
			return
		}
	}

	// 3. Update status pesanan pembelian menjadi "Diterima"
	_, err = tx.Exec("UPDATE pembelian SET status = 'Diterima' WHERE pembelian_id = ?", pembelianID)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status pembelian"})
		return
	}

	// 4. Jika semua berhasil, commit transaksi
	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyelesaikan transaksi"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pesanan berhasil diterima dan stok telah diperbarui"})
}

// =================================================================
// HANDLER UNTUK MODUL DASHBOARD
// =================================================================

func getDashboardStatsHandler(c *gin.Context) {
	var stats DashboardStats
	var err error

	// Hitung jumlah produk
	err = database.DB.QueryRow("SELECT COUNT(*) FROM produk").Scan(&stats.JumlahProduk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung produk"})
		return
	}

	// Hitung jumlah supplier
	err = database.DB.QueryRow("SELECT COUNT(*) FROM supplier").Scan(&stats.JumlahSupplier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung supplier"})
		return
	}

	// Hitung jumlah pembelian
	err = database.DB.QueryRow("SELECT COUNT(*) FROM pembelian").Scan(&stats.JumlahPembelian)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung pembelian"})
		return
	}

	// Hitung jumlah gudang
	err = database.DB.QueryRow("SELECT COUNT(*) FROM gudang").Scan(&stats.JumlahGudang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung gudang"})
		return
	}

	// Jika semua berhasil, kirim data statistik sebagai respons JSON
	c.JSON(http.StatusOK, stats)
}

// HANDLER UNTUK DATA GRAFIK STOK
// ===============================
func getStokChartHandler(c *gin.Context) {
	query := `
        SELECT p.nama_produk, SUM(s.jumlah) as total_stok
        FROM stok s
        JOIN produk p ON s.produk_id = p.produk_id
        GROUP BY p.nama_produk
        ORDER BY total_stok DESC
    `
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data stok untuk grafik"})
		return
	}
	defer rows.Close()

	// Siapkan struct untuk diisi
	var chartData StokChartResponse
	chartData.Labels = make([]string, 0)
	chartData.Data = make([]int, 0)

	for rows.Next() {
		var namaProduk string
		var totalStok int
		if err := rows.Scan(&namaProduk, &totalStok); err != nil {
			log.Printf("Error scanning row stok chart: %v", err)
			continue
		}
		chartData.Labels = append(chartData.Labels, namaProduk)
		chartData.Data = append(chartData.Data, totalStok)
	}

	c.JSON(http.StatusOK, chartData)
}

// HANDLER UNTUK 5 PEMBELIAN TERAKHIR
// ===================================
func getPembelianTerakhirHandler(c *gin.Context) {
	// Query kita batasi hanya 5 hasil, diurutkan dari yang paling baru
	query := `
        SELECT 
            p.pembelian_id, p.supplier_id, s.nama_supplier, 
            p.tanggal_pesan, p.estimasi_tiba, p.total_biaya, p.status
        FROM pembelian p
        JOIN supplier s ON p.supplier_id = s.supplier_id
        ORDER BY p.tanggal_pesan DESC
        LIMIT 5
    `
	rows, err := database.DB.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pembelian"})
		return
	}
	defer rows.Close()

	// Kita bisa gunakan lagi struct PembelianResponse yang sudah ada
	daftarPembelian := make([]PembelianResponse, 0)

	for rows.Next() {
		var p PembelianResponse
		err := rows.Scan(&p.PembelianID, &p.SupplierID, &p.NamaSupplier, &p.TanggalPesan, &p.EstimasiTiba, &p.TotalBiaya, &p.Status)
		if err != nil {
			log.Printf("Error scanning row pembelian: %v", err)
			continue
		}
		daftarPembelian = append(daftarPembelian, p)
	}
	c.JSON(http.StatusOK, daftarPembelian)
}

// PASTIKAN ANDA MENYALIN SAMPAI BARIS INI
