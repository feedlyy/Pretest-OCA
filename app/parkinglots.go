package app

import "time"

type Car struct {
	Plat_nomor    string `form:"plat_nomor" json:"plat_nomor" validate:"required"`
	Warna         string `form:"warna" json:"warna" validate:"required"`
	Tipe          string `form:"tipe" json:"tipe" validate:"required"`
	Tanggal_masuk time.Time `form:"tanggal_masuk" json:"tanggal_masuk"`
}
