// file: scm-api/internal/models/supplier.go

package models

import "database/sql"

// Supplier merepresentasikan tabel supplier di database
type Supplier struct {
	SupplierID    int64           `json:"supplier_id"`
	NamaSupplier  string          `json:"nama_supplier"`
	Alamat        sql.NullString  `json:"alamat"`
	Kontak        sql.NullString  `json:"kontak"`
	ContactPerson sql.NullString  `json:"contact_person"`
	Rating        sql.NullFloat64 `json:"rating"`
}
