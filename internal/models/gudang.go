// file: scm-api/internal/models/gudang.go
package models

// Gudang merepresentasikan tabel gudang di database
type Gudang struct {
	GudangID   int64  `json:"gudang_id"`
	NamaGudang string `json:"nama_gudang"`
	Lokasi     string `json:"lokasi"`
}
