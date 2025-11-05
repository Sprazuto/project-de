/**
 * Utility functions to process and format realisasi data from API
 * to match the expected format for CardRealisasiBulan and CardRealisasiTahun components
 */

import { getCardColorsByProgress } from './cardUtils';

export function processRealisasiBulanData(apiData) {
    if (!apiData || !apiData.data) return [];

    return apiData.data.map(category => {
        const { category: cat, progress_formatted, items } = category;

        switch (cat) {
            case 'barjas':
                return {
                    title: "Persenstase Capaian Realisasi Barjas",
                    subtitle: `Januari - ${apiData.meta.month_name} ${apiData.meta.year}`,
                    hintTitle: "PERSENTASE CAPAIAN REALISASI BARJAS",
                    hintDescription: "Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.\nRealisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.\nNilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:\n(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan",
                    items: items.map(item => {
                        if (item.type === 'perencanaan' && item.detail) {
                            const { selesai, target, terlambat } = item.detail;
                            return {
                                label: "Perencanaan",
                                value: `${selesai}<small>${terlambat > 0 ? `<sup><code class="text-danger bg-white p-0">-${terlambat}</code></sup>` : ""}/${target}</small>`,
                                popoverTitle: `${selesai} dari ${target} paket selesai.`,
                                popoverContent: terlambat > 0 ? `<span class="text-danger p-0">${terlambat} paket terlambat.</span>` : null,
                            };
                        }
                        if (item.type === 'pemilihan' && item.detail) {
                            const { selesai, target, terlambat } = item.detail;
                            return {
                                label: "Pemilihan",
                                value: `${selesai}<small>${terlambat > 0 ? `<sup><code class="text-danger bg-white p-0">-${terlambat}</code></sup>` : ""}/${target}</small>`,
                                popoverTitle: `${selesai} dari ${target} paket selesai.`,
                                popoverContent: terlambat > 0 ? `<span class="text-danger p-0">${terlambat} paket terlambat.</span>` : null,
                            };
                        }
                        if (item.type === 'pengadaan' && item.detail) {
                            const { selesai, target, terlambat } = item.detail;
                            return {
                                label: "Pengadaan",
                                value: `${selesai}<small>${terlambat > 0 ? `<sup><code class="text-danger bg-white p-0">-${terlambat}</code></sup>` : ""}/${target}</small>`,
                                popoverTitle: `${selesai} dari ${target} paket selesai.`,
                                popoverContent: terlambat > 0 ? `<span class="text-danger p-0">${terlambat} paket terlambat.</span>` : null,
                            };
                        }
                        if (item.type === 'penyerahan' && item.detail) {
                            const { selesai, target, terlambat } = item.detail;
                            return {
                                label: "Penyerahan",
                                value: `${selesai}<small>${terlambat > 0 ? `<sup><code class="text-danger bg-white p-0">-${terlambat}</code></sup>` : ""}/${target}</small>`,
                                popoverTitle: `${selesai} dari ${target} paket selesai.`,
                                popoverContent: terlambat > 0 ? `<span class="text-danger p-0">${terlambat} paket terlambat.</span>` : null,
                            };
                        }
                        return null;
                    }).filter(Boolean),
                    progress: progress_formatted,
                };

            case 'fisik':
                return {
                    title: "Persenstase Capaian Realisasi Fisik",
                    subtitle: `${apiData.meta.month_name} ${apiData.meta.year}`,
                    hintTitle: "PERSENTASE CAPAIAN REALISASI BARJAS",
                    hintDescription: "Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.\nRealisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.\nNilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:\n(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan",
                    items: items.map(item => {
                        if (item.type === 'realisasi') {
                            return { label: "Realisasi", value: item.formatted || item.value };
                        }
                        if (item.type === 'target') {
                            return { label: "Target", value: item.formatted || item.value };
                        }
                        return null;
                    }).filter(Boolean),
                    progress: progress_formatted,
                };

            case 'anggaran':
                return {
                    title: "Persenstase Capaian Realisasi Anggaran",
                    subtitle: `${apiData.meta.month_name} ${apiData.meta.year}`,
                    hintTitle: "PERSENTASE CAPAIAN REALISASI BARJAS",
                    hintDescription: "Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.\nRealisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.\nNilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:\n(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan",
                    items: items.map(item => {
                        if (item.type === 'realisasi') {
                            return { label: "Realisasi", value: item.formatted || item.value };
                        }
                        if (item.type === 'target') {
                            return { label: "Target", value: item.formatted || item.value };
                        }
                        return null;
                    }).filter(Boolean),
                    progress: progress_formatted,
                    layout: "rows",
                };

            case 'kinerja':
                return {
                    title: "Persenstase Capaian Realisasi Kinerja",
                    subtitle: `${apiData.meta.month_name} ${apiData.meta.year}`,
                    hintTitle: "PERSENTASE CAPAIAN REALISASI BARJAS",
                    hintDescription: "Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.\nRealisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.\nNilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:\n(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan",
                    items: items.map(item => {
                        if (item.type === 'realisasi') {
                            return { label: "Realisasi", value: item.formatted || item.value };
                        }
                        if (item.type === 'target') {
                            return { label: "Target", value: item.formatted || item.value };
                        }
                        return null;
                    }).filter(Boolean),
                    progress: progress_formatted,
                };

            default:
                return null;
        }
    }).filter(Boolean);
}

export function processRealisasiTahunData(apiData) {
    if (!apiData || !apiData.data) return [];

    return apiData.data.map(category => {
        const { category: cat, progress_formatted, capaian, items } = category;

        switch (cat) {
            case 'barjas':
                return {
                    title: "Progres Tahunan Capaian Barjas",
                    subtitle: `per-${apiData.meta.month_name} ${apiData.meta.year}`,
                    hintTitle: "PERSENTASE CAPAIAN REALISASI BARJAS",
                    hintDescription: "Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.\nRealisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.\nNilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:\n(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan",
                    items: items.map(item => {
                        if (item.type === 'realisasi') {
                            return { label: "Realisasi", value: item.formatted || item.value };
                        }
                        if (item.type === 'target') {
                            return { label: "Target", value: item.formatted || item.value };
                        }
                        return null;
                    }).filter(Boolean),
                    progress: progress_formatted,
                    color: getCardColorsByProgress(capaian).bgColor,
                };

            case 'fisik':
                return {
                    title: "Progres Tahunan Capaian Fisik",
                    subtitle: `per-${apiData.meta.month_name} ${apiData.meta.year}`,
                    hintTitle: "PERSENTASE CAPAIAN REALISASI BARJAS",
                    hintDescription: "Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.\nRealisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.\nNilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:\n(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan",
                    items: items.map(item => {
                        if (item.type === 'realisasi') {
                            return { label: "Realisasi", value: item.formatted || item.value };
                        }
                        if (item.type === 'target') {
                            return { label: "Target", value: item.formatted || item.value };
                        }
                        return null;
                    }).filter(Boolean),
                    progress: progress_formatted,
                    color: getCardColorsByProgress(capaian).bgColor,
                };

            case 'anggaran':
                return {
                    title: "Progres Tahunan Capaian Anggaran",
                    subtitle: `per-${apiData.meta.month_name} ${apiData.meta.year}`,
                    hintTitle: "PERSENTASE CAPAIAN REALISASI BARJAS",
                    hintDescription: "Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.\nRealisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.\nNilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:\n(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan",
                    items: items.map(item => {
                        if (item.type === 'realisasi') {
                            return { label: "Realisasi", value: item.formatted || item.value };
                        }
                        if (item.type === 'target') {
                            return { label: "Target", value: item.formatted || item.value };
                        }
                        return null;
                    }).filter(Boolean),
                    progress: progress_formatted,
                    layout: "rows",
                    color: getCardColorsByProgress(capaian).bgColor,
                };

            case 'kinerja':
                return {
                    title: "Progres Tahunan Capaian Kinerja",
                    subtitle: `per-${apiData.meta.month_name} ${apiData.meta.year}`,
                    hintTitle: "PERSENTASE CAPAIAN REALISASI BARJAS",
                    hintDescription: "Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.\nRealisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.\nNilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:\n(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan",
                    items: items.map(item => {
                        if (item.type === 'realisasi') {
                            return { label: "Realisasi", value: item.formatted || item.value };
                        }
                        if (item.type === 'target') {
                            return { label: "Target", value: item.formatted || item.value };
                        }
                        return null;
                    }).filter(Boolean),
                    progress: progress_formatted,
                    color: getCardColorsByProgress(capaian).bgColor,
                };

            default:
                return null;
        }
    }).filter(Boolean);
}