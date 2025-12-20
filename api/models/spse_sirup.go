package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Massad/gin-boilerplate/db"
)

// SPSESirupPerencanaan represents enriched planning stage data with SIRUP fields
type SPSESirupPerencanaan struct {
	ID                int64      `db:"id, primarykey, autoincrement" json:"id"`
	KodeRUP           string     `db:"kode_rup" json:"kode_rup"`
	SatuanKerja       string     `db:"satuan_kerja" json:"satuan_kerja"`
	NamaPaket         string     `db:"nama_paket" json:"nama_paket"`
	MetodePemilihan   string     `db:"metode_pemilihan" json:"metode_pemilihan"`
	TanggalPengumuman string     `db:"tanggal_pengumuman" json:"tanggal_pengumuman"`
	RencanaPemilihan  string     `db:"rencana_pemilihan" json:"rencana_pemilihan"`
	PaguRUP           string     `db:"pagu_rup" json:"pagu_rup"`
	KodeSatuanKerja   string     `db:"kode_satuan_kerja" json:"kode_satuan_kerja"`
	CaraPengadaan     string     `db:"cara_pengadaan" json:"cara_pengadaan"`
	JenisPengadaan    string     `db:"jenis_pengadaan" json:"jenis_pengadaan"`
	PDN               string     `db:"pdn" json:"pdn"`
	UMK               string     `db:"umk" json:"umk"`
	SumberDana        string     `db:"sumber_dana" json:"sumber_dana"`
	KodeRUPLokal      string     `db:"kode_rup_lokal" json:"kode_rup_lokal"`
	AkhirPemilihan    string     `db:"akhir_pemilihan" json:"akhir_pemilihan"`
	TipeSwakelola     string     `db:"tipe_swakelola" json:"tipe_swakelola"`

	// Additional SIRUP fields
	Dates             string     `db:"dates" json:"dates"`                        // JSON string for SIRUP dates
	SumberDanaSirup   string     `db:"sumber_dana_sirup" json:"sumber_dana_sirup"` // JSON string for SIRUP funding sources array
	LokasiPekerjaan   string     `db:"lokasi_pekerjaan" json:"lokasi_pekerjaan"`   // JSON string for SIRUP work locations array

	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	LastUpdate        int64      `db:"last_update" json:"last_update"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deleted_at"`
}

func (s SPSESirupPerencanaan) TableName() string {
	return "spse_perencanaansirup"
}

// RunSPSESirupMigrations runs SIRUP-specific database migrations
func RunSPSESirupMigrations() error {
	log.Println("Running SPSE SIRUP database migrations...")

	// Register SPSE SIRUP models with gorp
	dbMap := db.GetDB()
	dbMap.AddTableWithName(SPSESirupPerencanaan{}, "spse_perencanaansirup").SetKeys(true, "id")

	// Create tables using gorp
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		return fmt.Errorf("failed to create SPSE SIRUP tables: %v", err)
	}

	// Add unique constraints
	tables := []struct {
		name   string
		fields []string
	}{
		{
			name:   "spse_perencanaansirup",
			fields: []string{"kode_rup", "nama_paket"},
		},
	}

	for _, table := range tables {
		constraintName := fmt.Sprintf("unique_sirup_%s_kode_nama", table.name)

		// Drop existing partial index if it exists
		dropSQL := fmt.Sprintf("DROP INDEX IF EXISTS %s", constraintName)
		_, dropErr := dbMap.Db.Exec(dropSQL)
		if dropErr != nil {
			log.Printf("Warning: Failed to drop existing index for %s: %v", table.name, dropErr)
		}

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

	// Add indexes for performance
	indexes := []struct {
		name   string
		table  string
		fields []string
	}{
		{"idx_sirup_perencanaan_kode_rup", "spse_perencanaansirup", []string{"kode_rup"}},
		{"idx_sirup_perencanaan_kode_satker", "spse_perencanaansirup", []string{"kode_satuan_kerja"}},
		{"idx_sirup_perencanaan_deleted_at", "spse_perencanaansirup", []string{"deleted_at"}},
	}

	for _, idx := range indexes {
		indexSQL := fmt.Sprintf(
			"CREATE INDEX IF NOT EXISTS %s ON %s (%s)",
			idx.name,
			idx.table,
			strings.Join(idx.fields, ", "),
		)

		_, err := dbMap.Db.Exec(indexSQL)
		if err != nil {
			log.Printf("Warning: Failed to create index %s: %v", idx.name, err)
		}
	}

	log.Println("SPSE SIRUP database migrations completed successfully")
	return nil
}

// CreateSPSESirupTables creates the SPSE SIRUP tables directly without the migration system
func CreateSPSESirupTables() error {
	log.Println("Creating SPSE SIRUP database tables...")
	return RunSPSESirupMigrations()
}
