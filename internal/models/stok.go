// file: scm-api/internal/models/stok.go
package models

// Stok merepresentasikan tabel stok di database
type Stok struct {
	StokID        int64  `json:"stok_id"`
	ProdukID      int64  `json:"produk_id"`
	GudangID      int64  `json:"gudang_id"`
	Jumlah        int    `json:"jumlah"`
	TanggalUpdate string `json:"tanggal_update"`
}
