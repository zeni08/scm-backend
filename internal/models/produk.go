package models

import "database/sql"

// Produk merepresentasikan tabel produk
type Produk struct {
    ProdukID      int64           `json:"produk_id"`
    SKU           string          `json:"sku"`
    NamaProduk    string          `json:"nama_produk"`
    Deskripsi     sql.NullString  `json:"deskripsi"`
    Kategori      sql.NullString  `json:"kategori"`
    Satuan        string          `json:"satuan"`
    HargaJual     float64         `json:"harga_jual"`
    BeratKg       sql.NullFloat64 `json:"berat_kg"`
    GambarProduk  sql.NullString  `json:"gambar_produk"`
    SupplierID    sql.NullInt64   `json:"supplier_id"`
}