/*
Navicat Premium Data Transfer

Source Server         : sijagur.sumedangkab.go.id - VPC
Source Server Type    : MySQL
Source Server Version : 50743
Source Host           : localhost:3306
Source Schema         : db_sijagur

Target Server Type    : MySQL
Target Server Version : 50743
File Encoding         : 65001

Date: 03/11/2025 06:21:20
*/

SET NAMES utf8mb4;

SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for de_detail_anggaran
-- ----------------------------
DROP TABLE IF EXISTS `de_detail_anggaran`;

CREATE TABLE `de_detail_anggaran` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `id_ranking_opd` int(11) DEFAULT NULL,
  `c_anggaran_target` double DEFAULT NULL,
  `c_anggaran_realisasi` double DEFAULT NULL,
  `k_anggaran_target` double DEFAULT NULL,
  `k_anggaran_realisasi` double DEFAULT NULL,
  `p_anggaran_target` double DEFAULT NULL,
  `p_anggaran_realisasi` double DEFAULT NULL,
  `last_update` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=617 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for de_detail_barjas
-- ----------------------------
DROP TABLE IF EXISTS `de_detail_barjas`;

CREATE TABLE `de_detail_barjas` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `id_ranking_opd` int(11) DEFAULT NULL,
  `c_barjas_target` double DEFAULT NULL,
  `c_barjas_realisasi` double DEFAULT NULL,
  `k_barjas_target` double DEFAULT NULL,
  `k_barjas_realisasi` double DEFAULT NULL,
  `p_barjas_target` double DEFAULT NULL,
  `p_barjas_realisasi` double DEFAULT NULL,
  `c_perencanaan_selesai` int(11) DEFAULT NULL,
  `c_perencanaan_terlambat` int(11) DEFAULT NULL,
  `c_perencanaan_target` int(11) DEFAULT NULL,
  `c_pemilihan_selesai` int(11) DEFAULT NULL,
  `c_pemilihan_terlambat` int(11) DEFAULT NULL,
  `c_pemilihan_target` int(11) DEFAULT NULL,
  `c_pengadaan_selesai` int(11) DEFAULT NULL,
  `c_pengadaan_terlambat` int(11) DEFAULT NULL,
  `c_pengadaan_target` int(11) DEFAULT NULL,
  `c_penyerahan_selesai` int(11) DEFAULT NULL,
  `c_penyerahan_terlambat` int(11) DEFAULT NULL,
  `c_penyerahan_target` int(11) DEFAULT NULL,
  `k_perencanaan_selesai` int(11) DEFAULT NULL,
  `k_perencanaan_terlambat` int(11) DEFAULT NULL,
  `k_perencanaan_target` int(11) DEFAULT NULL,
  `k_pemilihan_selesai` int(11) DEFAULT NULL,
  `k_pemilihan_terlambat` int(11) DEFAULT NULL,
  `k_pemilihan_target` int(11) DEFAULT NULL,
  `k_pengadaan_selesai` int(11) DEFAULT NULL,
  `k_pengadaan_terlambat` int(11) DEFAULT NULL,
  `k_pengadaan_target` int(11) DEFAULT NULL,
  `k_penyerahan_selesai` int(11) DEFAULT NULL,
  `k_penyerahan_terlambat` int(11) DEFAULT NULL,
  `k_penyerahan_target` int(11) DEFAULT NULL,
  `p_perencanaan_selesai` int(11) DEFAULT NULL,
  `p_perencanaan_terlambat` int(11) DEFAULT NULL,
  `p_perencanaan_target` int(11) DEFAULT NULL,
  `p_pemilihan_selesai` int(11) DEFAULT NULL,
  `p_pemilihan_terlambat` int(11) DEFAULT NULL,
  `p_pemilihan_target` int(11) DEFAULT NULL,
  `p_pengadaan_selesai` int(11) DEFAULT NULL,
  `p_pengadaan_terlambat` int(11) DEFAULT NULL,
  `p_pengadaan_target` int(11) DEFAULT NULL,
  `p_penyerahan_selesai` int(11) DEFAULT NULL,
  `p_penyerahan_terlambat` int(11) DEFAULT NULL,
  `p_penyerahan_target` int(11) DEFAULT NULL,
  `last_update` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=617 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for de_detail_fisik
-- ----------------------------
DROP TABLE IF EXISTS `de_detail_fisik`;

CREATE TABLE `de_detail_fisik` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `id_ranking_opd` int(11) DEFAULT NULL,
  `c_fisik_target` double DEFAULT NULL,
  `c_fisik_realisasi` double DEFAULT NULL,
  `k_fisik_target` double DEFAULT NULL,
  `k_fisik_realisasi` double DEFAULT NULL,
  `p_fisik_target` double DEFAULT NULL,
  `p_fisik_realisasi` double DEFAULT NULL,
  `last_update` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=617 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for de_detail_kinerja
-- ----------------------------
DROP TABLE IF EXISTS `de_detail_kinerja`;

CREATE TABLE `de_detail_kinerja` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `id_ranking_opd` int(11) DEFAULT NULL,
  `c_kinerja_target` double DEFAULT NULL,
  `c_kinerja_realisasi` double DEFAULT NULL,
  `k_kinerja_target` double DEFAULT NULL,
  `k_kinerja_realisasi` double DEFAULT NULL,
  `p_kinerja_target` double DEFAULT NULL,
  `p_kinerja_realisasi` double DEFAULT NULL,
  `last_update` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=617 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for de_peta_detail
-- ----------------------------
DROP TABLE IF EXISTS `de_peta_detail`;

CREATE TABLE `de_peta_detail` (
  `id_detail` int(11) NOT NULL AUTO_INCREMENT,
  `kode_gadm` text,
  `koordinat` text,
  `tahun` year(4) DEFAULT NULL,
  `id_skpd` int(11) DEFAULT NULL,
  `idsatker` int(11) DEFAULT NULL,
  `nama_opd` varchar(255) DEFAULT NULL,
  `id_kontrak` int(11) DEFAULT NULL,
  `id_rup` int(11) DEFAULT NULL,
  `nama_paket` varchar(255) DEFAULT NULL,
  `status` enum('perencanaan','pelaksanaan','selesai') DEFAULT 'perencanaan',
  `progres` double DEFAULT NULL,
  `tanggal_mulai` date DEFAULT NULL,
  `tanggal_selesai` date DEFAULT NULL,
  `is_removed` int(11) DEFAULT NULL,
  `last_update` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id_detail`)
) ENGINE=InnoDB AUTO_INCREMENT=311 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for de_peta_kecamatan
-- ----------------------------
DROP TABLE IF EXISTS `de_peta_kecamatan`;

CREATE TABLE `de_peta_kecamatan` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `kode_gadm` varchar(20) NOT NULL,
  `nama_kecamatan` varchar(100) NOT NULL,
  `nama_kabupaten` varchar(100) DEFAULT 'Sumedang',
  `latitude` decimal(10,6) DEFAULT NULL,
  `longitude` decimal(10,6) DEFAULT NULL,
  `total_paket` int(11) DEFAULT NULL,
  `persentase_progres` double DEFAULT NULL,
  `jumlah_opd` int(11) DEFAULT NULL,
  `label_status` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `kode_gadm` (`kode_gadm`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for de_ranking_opd
-- ----------------------------
DROP TABLE IF EXISTS `de_ranking_opd`;

CREATE TABLE `de_ranking_opd` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `id_skpd` int(11) NOT NULL,
  `idsatker` int(11) DEFAULT NULL,
  `nama_opd` varchar(255) DEFAULT NULL,
  `jenis_opd` varchar(255) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `capaian_opd` double DEFAULT NULL,
  `capaian_barjas` double DEFAULT NULL,
  `capaian_fisik` double DEFAULT NULL,
  `capaian_anggaran` double DEFAULT NULL,
  `capaian_kinerja` double DEFAULT NULL,
  `kumulatif_opd` double DEFAULT NULL,
  `kumulatif_barjas` double DEFAULT NULL,
  `kumulatif_fisik` double DEFAULT NULL,
  `kumulatif_anggaran` double DEFAULT NULL,
  `kumulatif_kinerja` double DEFAULT NULL,
  `periodik_opd` double DEFAULT NULL,
  `periodik_barjas` double DEFAULT NULL,
  `periodik_fisik` double DEFAULT NULL,
  `periodik_anggaran` double DEFAULT NULL,
  `periodik_kinerja` double DEFAULT NULL,
  `peringkat_opd` int(11) DEFAULT NULL,
  `peringkat_barjas` int(11) DEFAULT NULL,
  `peringkat_fisik` int(11) DEFAULT NULL,
  `peringkat_anggaran` int(11) DEFAULT NULL,
  `peringkat_kinerja` int(11) DEFAULT NULL,
  `tahun` int(4) DEFAULT NULL,
  `bulan` int(2) DEFAULT NULL,
  `last_update` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=617 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for de_status_paket
-- ----------------------------
DROP TABLE IF EXISTS `de_status_paket`;

CREATE TABLE `de_status_paket` (
  `id_rup` int(11) NOT NULL,
  `idsatker` int(11) DEFAULT NULL,
  `tahun` int(4) DEFAULT NULL,
  `bulan` int(2) DEFAULT NULL,
  `skor` double DEFAULT NULL,
  `status_perencanaan` varchar(255) DEFAULT NULL COMMENT 'Unset; No; Start; Process; Yes; Overdue',
  `status_pemilihan` varchar(255) DEFAULT NULL,
  `status_pengadaan` varchar(255) DEFAULT NULL,
  `status_penyerahan` varchar(255) DEFAULT NULL,
  `overdue_pemilihan` date DEFAULT NULL,
  `overdue_pengadaan` date DEFAULT NULL,
  `overdue_penyerahan` date DEFAULT NULL,
  `last_update` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `is_removed` int(11) DEFAULT '0',
  PRIMARY KEY (`id_rup`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;