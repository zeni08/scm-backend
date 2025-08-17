// file: scm-api/internal/models/pembelian.go

package models

import "database/sql"

// Pembelian merepresentasikan tabel 'pembelian' (header transaksi)
type Pembelian struct {
	PembelianID  int64           `json:"pembelian_id"`
	SupplierID   int64           `json:"supplier_id"`
	TanggalPesan string          `json:"tanggal_pesan"` // Menggunakan string untuk kemudahan
	EstimasiTiba sql.NullString  `json:"estimasi_tiba"`
	TotalBiaya   sql.NullFloat64 `json:"total_biaya"`
	Status       string          `json:"status"`
}
