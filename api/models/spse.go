package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Massad/gin-boilerplate/db"
)

// SPSEPerencanaan represents planning stage data
type SPSEPerencanaan struct {
	ID                int64     `db:"id, primarykey, autoincrement" json:"id"`
	KodeRUP           string    `db:"kode_rup" json:"kode_rup"`
	SatuanKerja       string    `db:"satuan_kerja" json:"satuan_kerja"`
	NamaPaket         string    `db:"nama_paket" json:"nama_paket"`
	MetodePemilihan   string    `db:"metode_pemilihan" json:"metode_pemilihan"`
	TanggalPengumuman string    `db:"tanggal_pengumuman" json:"tanggal_pengumuman"`
	RencanaPemilihan  string    `db:"rencana_pemilihan" json:"rencana_pemilihan"`
	PaguRUP           string    `db:"pagu_rup" json:"pagu_rup"`
	KodeSatuanKerja   string    `db:"kode_satuan_kerja" json:"kode_satuan_kerja"`
	CaraPengadaan     string    `db:"cara_pengadaan" json:"cara_pengadaan"`
	JenisPengadaan    string    `db:"jenis_pengadaan" json:"jenis_pengadaan"`
	PDN               string    `db:"pdn" json:"pdn"`
	UMK               string    `db:"umk" json:"umk"`
	SumberDana        string    `db:"sumber_dana" json:"sumber_dana"`
	KodeRUPLokal      string    `db:"kode_rup_lokal" json:"kode_rup_lokal"`
	AkhirPemilihan    string    `db:"akhir_pemilihan" json:"akhir_pemilihan"`
	TipeSwakelola     string    `db:"tipe_swakelola" json:"tipe_swakelola"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	LastUpdate        int64     `db:"last_update" json:"last_update"`
}

func (s SPSEPerencanaan) TableName() string {
	return "spse_perencanaan"
}

// SPSEPersiapan represents preparation stage data
type SPSEPersiapan struct {
	ID               int64     `db:"id, primarykey, autoincrement" json:"id"`
	KodeRUP          string    `db:"kode_rup" json:"kode_rup"`
	SatuanKerja      string    `db:"satuan_kerja" json:"satuan_kerja"`
	NamaPaket        string    `db:"nama_paket" json:"nama_paket"`
	MetodePemilihan  string    `db:"metode_pemilihan" json:"metode_pemilihan"`
	TanggalBuatPaket string    `db:"tanggal_buat_paket" json:"tanggal_buat_paket"`
	NilaiPaguRUP     string    `db:"nilai_pagu_rup" json:"nilai_pagu_rup"`
	NilaiPaguPaket   string    `db:"nilai_pagu_paket" json:"nilai_pagu_paket"`
	KodeSatuanKerja  string    `db:"kode_satuan_kerja" json:"kode_satuan_kerja"`
	CaraPengadaan    string    `db:"cara_pengadaan" json:"cara_pengadaan"`
	JenisPengadaan   string    `db:"jenis_pengadaan" json:"jenis_pengadaan"`
	PDN              string    `db:"pdn" json:"pdn"`
	UMK              string    `db:"umk" json:"umk"`
	SumberDana       string    `db:"sumber_dana" json:"sumber_dana"`
	KodeRUPLokal     string    `db:"kode_rup_lokal" json:"kode_rup_lokal"`
	MetodePengadaan  string    `db:"metode_pengadaan" json:"metode_pengadaan"`
	TipeSwakelola    string    `db:"tipe_swakelola" json:"tipe_swakelola"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	LastUpdate       int64     `db:"last_update" json:"last_update"`
}

func (s SPSEPersiapan) TableName() string {
	return "spse_persiapan"
}

// SPSEPemilihan represents selection/contract stage data
type SPSEPemilihan struct {
	ID               int64     `db:"id, primarykey, autoincrement" json:"id"`
	KodeRUP          string    `db:"kode_rup" json:"kode_rup"`
	SatuanKerja      string    `db:"satuan_kerja" json:"satuan_kerja"`
	NamaPaket        string    `db:"nama_paket" json:"nama_paket"`
	MetodePemilihan  string    `db:"metode_pemilihan" json:"metode_pemilihan"`
	RencanaPemilihan string    `db:"rencana_pemilihan" json:"rencana_pemilihan"`
	TanggalPemilihan string    `db:"tanggal_pemilihan" json:"tanggal_pemilihan"`
	NilaiHPS         string    `db:"nilai_hps" json:"nilai_hps"`
	StatusPaket      string    `db:"status_paket" json:"status_paket"`
	KodeSatuanKerja  string    `db:"kode_satuan_kerja" json:"kode_satuan_kerja"`
	CaraPengadaan    string    `db:"cara_pengadaan" json:"cara_pengadaan"`
	JenisPengadaan   string    `db:"jenis_pengadaan" json:"jenis_pengadaan"`
	PDN              string    `db:"pdn" json:"pdn"`
	UMK              string    `db:"umk" json:"umk"`
	SumberDana       string    `db:"sumber_dana" json:"sumber_dana"`
	KodeRUPLokal     string    `db:"kode_rup_lokal" json:"kode_rup_lokal"`
	MetodePengadaan  string    `db:"metode_pengadaan" json:"metode_pengadaan"`
	PaguRUP          string    `db:"pagu_rup" json:"pagu_rup"`
	TipeSwakelola    string    `db:"tipe_swakelola" json:"tipe_swakelola"`
	AkhirPemilihan   string    `db:"akhir_pemilihan" json:"akhir_pemilihan"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	LastUpdate       int64     `db:"last_update" json:"last_update"`
}

func (s SPSEPemilihan) TableName() string {
	return "spse_pemilihan"
}

// RunSPESEMigrations runs SPSE-specific database migrations
func RunSPESEMigrations() error {
	log.Println("Running SPSE database migrations...")

	// Register SPSE models with gorp
	dbMap := db.GetDB()
	dbMap.AddTableWithName(SPSEPerencanaan{}, "spse_perencanaan").SetKeys(true, "id")
	dbMap.AddTableWithName(SPSEPersiapan{}, "spse_persiapan").SetKeys(true, "id")
	dbMap.AddTableWithName(SPSEPemilihan{}, "spse_pemilihan").SetKeys(true, "id")

	// Create tables using gorp
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		return fmt.Errorf("failed to create SPSE tables: %v", err)
	}

	// Add unique constraints
	tables := []struct {
		name   string
		fields []string
	}{
		{
			name:   "spse_perencanaan",
			fields: []string{"kode_rup", "nama_paket"},
		},
		{
			name:   "spse_persiapan",
			fields: []string{"kode_rup", "nama_paket"},
		},
		{
			name:   "spse_pemilihan",
			fields: []string{"kode_rup", "nama_paket"},
		},
	}

	for _, table := range tables {
		constraintName := fmt.Sprintf("unique_%s_kode_nama", table.name)
		constraintSQL := fmt.Sprintf(
			"CREATE UNIQUE INDEX IF NOT EXISTS %s ON %s (%s)",
			constraintName,
			table.name,
			strings.Join(table.fields, ", "),
		)

		_, err := dbMap.Db.Exec(constraintSQL)
		if err != nil {
			log.Printf("Warning: Failed to create unique constraint for %s: %v", table.name, err)
		}
	}

	log.Println("SPSE database migrations completed successfully")
	return nil
}

// CreateSPASETables creates the SPSE tables directly without the migration system
func CreateSPASETables() error {
	log.Println("Creating SPSE database tables...")

	return RunSPESEMigrations()
}
