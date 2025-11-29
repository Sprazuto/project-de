package models

import (
	"fmt"
	"log"

	"github.com/Massad/gin-boilerplate/db"
)

// DeDetailAnggaran ...
type DeDetailAnggaran struct {
	ID                 int64   `db:"id, primarykey, autoincrement" json:"id"`
	IdRankingOpd       int64   `db:"id_ranking_opd" json:"id_ranking_opd"`
	CAnggaranTarget    float64 `db:"c_anggaran_target" json:"c_anggaran_target"`
	CAnggaranRealisasi float64 `db:"c_anggaran_realisasi" json:"c_anggaran_realisasi"`
	KAnggaranTarget    float64 `db:"k_anggaran_target" json:"k_anggaran_target"`
	KAnggaranRealisasi float64 `db:"k_anggaran_realisasi" json:"k_anggaran_realisasi"`
	PAnggaranTarget    float64 `db:"p_anggaran_target" json:"p_anggaran_target"`
	PAnggaranRealisasi float64 `db:"p_anggaran_realisasi" json:"p_anggaran_realisasi"`
	LastUpdate         int64   `db:"last_update" json:"last_update"`
}

func (d DeDetailAnggaran) TableName() string {
	return "de_detail_anggaran"
}

// DeDetailBarjas ...
type DeDetailBarjas struct {
	ID                    int64   `db:"id, primarykey, autoincrement" json:"id"`
	IdRankingOpd          int64   `db:"id_ranking_opd" json:"id_ranking_opd"`
	CBarjasTarget         float64 `db:"c_barjas_target" json:"c_barjas_target"`
	CBarjasRealisasi      float64 `db:"c_barjas_realisasi" json:"c_barjas_realisasi"`
	KBarjasTarget         float64 `db:"k_barjas_target" json:"k_barjas_target"`
	KBarjasRealisasi      float64 `db:"k_barjas_realisasi" json:"k_barjas_realisasi"`
	PBarjasTarget         float64 `db:"p_barjas_target" json:"p_barjas_target"`
	PBarjasRealisasi      float64 `db:"p_barjas_realisasi" json:"p_barjas_realisasi"`
	CPerencanaanSelesai   int64   `db:"c_perencanaan_selesai" json:"c_perencanaan_selesai"`
	CPerencanaanTerlambat int64   `db:"c_perencanaan_terlambat" json:"c_perencanaan_terlambat"`
	CPerencanaanTarget    int64   `db:"c_perencanaan_target" json:"c_perencanaan_target"`
	CPemilihanSelesai     int64   `db:"c_pemilihan_selesai" json:"c_pemilihan_selesai"`
	CPemilihanTerlambat   int64   `db:"c_pemilihan_terlambat" json:"c_pemilihan_terlambat"`
	CPemilihanTarget      int64   `db:"c_pemilihan_target" json:"c_pemilihan_target"`
	CPengadaanSelesai     int64   `db:"c_pengadaan_selesai" json:"c_pengadaan_selesai"`
	CPengadaanTerlambat   int64   `db:"c_pengadaan_terlambat" json:"c_pengadaan_terlambat"`
	CPengadaanTarget      int64   `db:"c_pengadaan_target" json:"c_pengadaan_target"`
	CPenyerahanSelesai    int64   `db:"c_penyerahan_selesai" json:"c_penyerahan_selesai"`
	CPenyerahanTerlambat  int64   `db:"c_penyerahan_terlambat" json:"c_penyerahan_terlambat"`
	CPenyerahanTarget     int64   `db:"c_penyerahan_target" json:"c_penyerahan_target"`
	KPerencanaanSelesai   int64   `db:"k_perencanaan_selesai" json:"k_perencanaan_selesai"`
	KPerencanaanTerlambat int64   `db:"k_perencanaan_terlambat" json:"k_perencanaan_terlambat"`
	KPerencanaanTarget    int64   `db:"k_perencanaan_target" json:"k_perencanaan_target"`
	KPemilihanSelesai     int64   `db:"k_pemilihan_selesai" json:"k_pemilihan_selesai"`
	KPemilihanTerlambat   int64   `db:"k_pemilihan_terlambat" json:"k_pemilihan_terlambat"`
	KPemilihanTarget      int64   `db:"k_pemilihan_target" json:"k_pemilihan_target"`
	KPengadaanSelesai     int64   `db:"k_pengadaan_selesai" json:"k_pengadaan_selesai"`
	KPengadaanTerlambat   int64   `db:"k_pengadaan_terlambat" json:"k_pengadaan_terlambat"`
	KPengadaanTarget      int64   `db:"k_pengadaan_target" json:"k_pengadaan_target"`
	KPenyerahanSelesai    int64   `db:"k_penyerahan_selesai" json:"k_penyerahan_selesai"`
	KPenyerahanTerlambat  int64   `db:"k_penyerahan_terlambat" json:"k_penyerahan_terlambat"`
	KPenyerahanTarget     int64   `db:"k_penyerahan_target" json:"k_penyerahan_target"`
	PPerencanaanSelesai   int64   `db:"p_perencanaan_selesai" json:"p_perencanaan_selesai"`
	PPerencanaanTerlambat int64   `db:"p_perencanaan_terlambat" json:"p_perencanaan_terlambat"`
	PPerencanaanTarget    int64   `db:"p_perencanaan_target" json:"p_perencanaan_target"`
	PPemilihanSelesai     int64   `db:"p_pemilihan_selesai" json:"p_pemilihan_selesai"`
	PPemilihanTerlambat   int64   `db:"p_pemilihan_terlambat" json:"p_pemilihan_terlambat"`
	PPemilihanTarget      int64   `db:"p_pemilihan_target" json:"p_pemilihan_target"`
	PPengadaanSelesai     int64   `db:"p_pengadaan_selesai" json:"p_pengadaan_selesai"`
	PPengadaanTerlambat   int64   `db:"p_pengadaan_terlambat" json:"p_pengadaan_terlambat"`
	PPengadaanTarget      int64   `db:"p_pengadaan_target" json:"p_pengadaan_target"`
	PPenyerahanSelesai    int64   `db:"p_penyerahan_selesai" json:"p_penyerahan_selesai"`
	PPenyerahanTerlambat  int64   `db:"p_penyerahan_terlambat" json:"p_penyerahan_terlambat"`
	PPenyerahanTarget     int64   `db:"p_penyerahan_target" json:"p_penyerahan_target"`
	LastUpdate            int64   `db:"last_update" json:"last_update"`
}

func (d DeDetailBarjas) TableName() string {
	return "de_detail_barjas"
}

// DeDetailFisik ...
type DeDetailFisik struct {
	ID              int64   `db:"id, primarykey, autoincrement" json:"id"`
	IdRankingOpd    int64   `db:"id_ranking_opd" json:"id_ranking_opd"`
	CFisikTarget    float64 `db:"c_fisik_target" json:"c_fisik_target"`
	CFisikRealisasi float64 `db:"c_fisik_realisasi" json:"c_fisik_realisasi"`
	KFisikTarget    float64 `db:"k_fisik_target" json:"k_fisik_target"`
	KFisikRealisasi float64 `db:"k_fisik_realisasi" json:"k_fisik_realisasi"`
	PFisikTarget    float64 `db:"p_fisik_target" json:"p_fisik_target"`
	PFisikRealisasi float64 `db:"p_fisik_realisasi" json:"p_fisik_realisasi"`
	LastUpdate      int64   `db:"last_update" json:"last_update"`
}

func (d DeDetailFisik) TableName() string {
	return "de_detail_fisik"
}

// DeDetailKinerja ...
type DeDetailKinerja struct {
	ID                int64   `db:"id, primarykey, autoincrement" json:"id"`
	IdRankingOpd      int64   `db:"id_ranking_opd" json:"id_ranking_opd"`
	CKinerjaTarget    float64 `db:"c_kinerja_target" json:"c_kinerja_target"`
	CKinerjaRealisasi float64 `db:"c_kinerja_realisasi" json:"c_kinerja_realisasi"`
	KKinerjaTarget    float64 `db:"k_kinerja_target" json:"k_kinerja_target"`
	KKinerjaRealisasi float64 `db:"k_kinerja_realisasi" json:"k_kinerja_realisasi"`
	PKinerjaTarget    float64 `db:"p_kinerja_target" json:"p_kinerja_target"`
	PKinerjaRealisasi float64 `db:"p_kinerja_realisasi" json:"p_kinerja_realisasi"`
	LastUpdate        int64   `db:"last_update" json:"last_update"`
}

func (d DeDetailKinerja) TableName() string {
	return "de_detail_kinerja"
}

// DePetaDetail ...
type DePetaDetail struct {
	IdDetail       int64   `db:"id_detail, primarykey, autoincrement" json:"id_detail"`
	KodeGadm       string  `db:"kode_gadm" json:"kode_gadm"`
	Koordinat      string  `db:"koordinat" json:"koordinat"`
	Tahun          int     `db:"tahun" json:"tahun"`
	IdSkpd         int64   `db:"id_skpd" json:"id_skpd"`
	Idsatker       int64   `db:"idsatker" json:"idsatker"`
	NamaOpd        string  `db:"nama_opd" json:"nama_opd"`
	IdKontrak      int64   `db:"id_kontrak" json:"id_kontrak"`
	IdRup          int64   `db:"id_rup" json:"id_rup"`
	NamaPaket      string  `db:"nama_paket" json:"nama_paket"`
	Status         string  `db:"status" json:"status"`
	Progres        float64 `db:"progres" json:"progres"`
	TanggalMulai   string  `db:"tanggal_mulai" json:"tanggal_mulai"`
	TanggalSelesai string  `db:"tanggal_selesai" json:"tanggal_selesai"`
	IsRemoved      int64   `db:"is_removed" json:"is_removed"`
	LastUpdate     int64   `db:"last_update" json:"last_update"`
}

func (d DePetaDetail) TableName() string {
	return "de_peta_detail"
}

// DePetaKecamatan ...
type DePetaKecamatan struct {
	ID                int64   `db:"id, primarykey, autoincrement" json:"id"`
	KodeGadm          string  `db:"kode_gadm" json:"kode_gadm"`
	NamaKecamatan     string  `db:"nama_kecamatan" json:"nama_kecamatan"`
	NamaKabupaten     string  `db:"nama_kabupaten" json:"nama_kabupaten"`
	Latitude          float64 `db:"latitude" json:"latitude"`
	Longitude         float64 `db:"longitude" json:"longitude"`
	TotalPaket        int64   `db:"total_paket" json:"total_paket"`
	PersentaseProgres float64 `db:"persentase_progres" json:"persentase_progres"`
	JumlahOpd         int64   `db:"jumlah_opd" json:"jumlah_opd"`
	LabelStatus       string  `db:"label_status" json:"label_status"`
}

func (d DePetaKecamatan) TableName() string {
	return "de_peta_kecamatan"
}

// DeRankingOpd ...
type DeRankingOpd struct {
	ID                int64   `db:"id, primarykey, autoincrement" json:"id"`
	IdSkpd            int64   `db:"id_skpd" json:"id_skpd"`
	Idsatker          int64   `db:"idsatker" json:"idsatker"`
	NamaOpd           string  `db:"nama_opd" json:"nama_opd"`
	JenisOpd          string  `db:"jenis_opd" json:"jenis_opd"`
	Status            string  `db:"status" json:"status"`
	CapaianOpd        float64 `db:"capaian_opd" json:"capaian_opd"`
	CapaianBarjas     float64 `db:"capaian_barjas" json:"capaian_barjas"`
	CapaianFisik      float64 `db:"capaian_fisik" json:"capaian_fisik"`
	CapaianAnggaran   float64 `db:"capaian_anggaran" json:"capaian_anggaran"`
	CapaianKinerja    float64 `db:"capaian_kinerja" json:"capaian_kinerja"`
	KumulatifOpd      float64 `db:"kumulatif_opd" json:"kumulatif_opd"`
	KumulatifBarjas   float64 `db:"kumulatif_barjas" json:"kumulatif_barjas"`
	KumulatifFisik    float64 `db:"kumulatif_fisik" json:"kumulatif_fisik"`
	KumulatifAnggaran float64 `db:"kumulatif_anggaran" json:"kumulatif_anggaran"`
	KumulatifKinerja  float64 `db:"kumulatif_kinerja" json:"kumulatif_kinerja"`
	PeriodikOpd       float64 `db:"periodik_opd" json:"periodik_opd"`
	PeriodikBarjas    float64 `db:"periodik_barjas" json:"periodik_barjas"`
	PeriodikFisik     float64 `db:"periodik_fisik" json:"periodik_fisik"`
	PeriodikAnggaran  float64 `db:"periodik_anggaran" json:"periodik_anggaran"`
	PeriodikKinerja   float64 `db:"periodik_kinerja" json:"periodik_kinerja"`
	PeringkatOpd      int64   `db:"peringkat_opd" json:"peringkat_opd"`
	PeringkatBarjas   int64   `db:"peringkat_barjas" json:"peringkat_barjas"`
	PeringkatFisik    int64   `db:"peringkat_fisik" json:"peringkat_fisik"`
	PeringkatAnggaran int64   `db:"peringkat_anggaran" json:"peringkat_anggaran"`
	PeringkatKinerja  int64   `db:"peringkat_kinerja" json:"peringkat_kinerja"`
	Tahun             int     `db:"tahun" json:"tahun"`
	Bulan             int     `db:"bulan" json:"bulan"`
	LastUpdate        int64   `db:"last_update" json:"last_update"`
}

func (d DeRankingOpd) TableName() string {
	return "de_ranking_opd"
}

// DeStatusPaket ...
type DeStatusPaket struct {
	IdRup              int64   `db:"id_rup, primarykey" json:"id_rup"`
	Idsatker           int64   `db:"idsatker" json:"idsatker"`
	Tahun              int     `db:"tahun" json:"tahun"`
	Bulan              int     `db:"bulan" json:"bulan"`
	Skor               float64 `db:"skor" json:"skor"`
	StatusPerencanaan  string  `db:"status_perencanaan" json:"status_perencanaan"`
	StatusPemilihan    string  `db:"status_pemilihan" json:"status_pemilihan"`
	StatusPengadaan    string  `db:"status_pengadaan" json:"status_pengadaan"`
	StatusPenyerahan   string  `db:"status_penyerahan" json:"status_penyerahan"`
	OverduePerencanaan string  `db:"overdue_perencanaan" json:"overdue_perencanaan"`
	OverduePemilihan   string  `db:"overdue_pemilihan" json:"overdue_pemilihan"`
	OverduePengadaan   string  `db:"overdue_pengadaan" json:"overdue_pengadaan"`
	OverduePenyerahan  string  `db:"overdue_penyerahan" json:"overdue_penyerahan"`
	LastUpdate         int64   `db:"last_update" json:"last_update"`
	IsRemoved          int64   `db:"is_removed" json:"is_removed"`
}

func (d DeStatusPaket) TableName() string {
	return "de_status_paket"
}

// Migration represents a database migration
type Migration struct {
	Version  int
	Name     string
	UpFunc   func() error
	DownFunc func() error
}

// migrations list - add new migrations here
var migrations = []Migration{
	{
		Version: 1,
		Name:    "create_sijagur_tables",
		UpFunc: func() error {
			// Register sijagur models with gorp
			db.GetDB().AddTableWithName(DeDetailAnggaran{}, "de_detail_anggaran").SetKeys(true, "id")
			db.GetDB().AddTableWithName(DeDetailBarjas{}, "de_detail_barjas").SetKeys(true, "id")
			db.GetDB().AddTableWithName(DeDetailFisik{}, "de_detail_fisik").SetKeys(true, "id")
			db.GetDB().AddTableWithName(DeDetailKinerja{}, "de_detail_kinerja").SetKeys(true, "id")
			db.GetDB().AddTableWithName(DePetaDetail{}, "de_peta_detail").SetKeys(true, "id_detail")
			db.GetDB().AddTableWithName(DePetaKecamatan{}, "de_peta_kecamatan").SetKeys(true, "id")
			db.GetDB().AddTableWithName(DeRankingOpd{}, "de_ranking_opd").SetKeys(true, "id")
			db.GetDB().AddTableWithName(DeStatusPaket{}, "de_status_paket").SetKeys(true, "id_rup")

			// Create tables using gorp
			err := db.GetDB().CreateTablesIfNotExists()
			if err != nil {
				return fmt.Errorf("failed to create tables: %v", err)
			}
			return nil
		},
		DownFunc: func() error {
			// Drop tables
			tables := []string{"de_detail_anggaran", "de_detail_barjas", "de_detail_fisik", "de_detail_kinerja", "de_peta_detail", "de_peta_kecamatan", "de_ranking_opd", "de_status_paket"}
			for _, table := range tables {
				_, err := db.GetDB().Db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
				if err != nil {
					return fmt.Errorf("failed to drop table %s: %v", table, err)
				}
			}
			return nil
		},
	},
	{
		Version: 2,
		Name:    "add_overdue_perencanaan_column",
		UpFunc: func() error {
			// Add overdue_perencanaan column to de_status_paket
			_, err := db.GetDB().Db.Exec(`ALTER TABLE de_status_paket ADD COLUMN IF NOT EXISTS overdue_perencanaan TEXT`)
			if err != nil {
				return fmt.Errorf("failed to add overdue_perencanaan column: %v", err)
			}
			return nil
		},
		DownFunc: func() error {
			// Remove overdue_perencanaan column
			_, err := db.GetDB().Db.Exec(`ALTER TABLE de_status_paket DROP COLUMN IF EXISTS overdue_perencanaan`)
			if err != nil {
				return fmt.Errorf("failed to drop overdue_perencanaan column: %v", err)
			}
			return nil
		},
	},
}

// RunMigrations runs all pending migrations
func RunMigrations() error {
	log.Println("Running sijagur database migrations...")

	// Create migrations table if not exists
	_, err := db.GetDB().Db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get current migration version
	var currentVersion int
	err = db.GetDB().Db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("failed to get current migration version: %v", err)
	}

	// Run pending migrations
	for _, migration := range migrations {
		if migration.Version > currentVersion {
			log.Printf("Running migration %d: %s", migration.Version, migration.Name)
			err := migration.UpFunc()
			if err != nil {
				return fmt.Errorf("migration %d failed: %v", migration.Version, err)
			}

			// Record migration as executed
			_, err = db.GetDB().Db.Exec("INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", migration.Version, migration.Name)
			if err != nil {
				return fmt.Errorf("failed to record migration %d: %v", migration.Version, err)
			}
			log.Printf("Migration %d completed successfully", migration.Version)
		}
	}

	log.Println("Sijagur database migrations completed successfully")
	return nil
}
