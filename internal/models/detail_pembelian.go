package models

// DetailPembelian merepresentasikan tabel 'detail_pembelian' (item dalam transaksi)
type DetailPembelian struct {
	DetailPembelianID int64   `json:"detail_pembelian_id"`
	PembelianID       int64   `json:"pembelian_id"`
	ProdukID          int64   `json:"produk_id"`
	Jumlah            int     `json:"jumlah"`
	HargaBeliSatuan   float64 `json:"harga_beli_satuan"`
	Subtotal          float64 `json:"subtotal"`
}
