/**
 * Utility functions to process and format realisasi data from API
 * to match the expected format for CardRealisasiBulan and CardRealisasiTahun components
 * Adapted for web2 with Vuetify components
 */

export function getCardColorsByProgress(progress) {
  const progressValue = parseInt(progress)
  if (progressValue >= 75) {
    return {
      bgColor: 'primary',
      textColor: 'text-white',
      chartColors: ['#028C86']
    }
  } else if (progressValue >= 50) {
    return {
      bgColor: 'secondary',
      textColor: 'text-white',
      chartColors: ['#B1D663']
    }
  } else if (progressValue >= 25) {
    return {
      bgColor: 'error',
      textColor: 'text-white',
      chartColors: ['#EF4444']
    }
  } else {
    return {
      bgColor: 'dark',
      textColor: 'text-white',
      chartColors: ['#6B7280']
    }
  }
}

export function processRealisasiBulanData(apiData) {
  if (!apiData || !apiData.results || !apiData.results[0] || !apiData.results[0].data) return []

  const meta = apiData.results[0].meta

  return apiData.results[0].data
    .map((category) => {
      const { category: cat, progress_formatted, items } = category

      switch (cat) {
        case 'barjas':
          return {
            title: 'Persenstase Capaian<br>Realisasi Barjas',
            subtitle: `Januari - ${meta.month_name} ${meta.year}`,
            hintTitle: 'PERSENTASE CAPAIAN REALISASI BARJAS',
            hintDescription:
              'Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.<br>Realisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.<br>Nilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:<br>(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan',
            items: items
              .map((item) => {
                if (item.type === 'perencanaan' && item.detail) {
                  const { selesai, target, terlambat } = item.detail
                  return {
                    label: 'Perencanaan',
                    value: `${selesai}<small>${terlambat > 0 ? `<sup><code class="text-error v-card--variant-elevated">-${terlambat}</code></sup>` : ''}/${target}</small>`,
                    popoverTitle: `${selesai} dari ${target} paket selesai.`,
                    popoverContent:
                      terlambat > 0 ? `<span class="text-error">${terlambat} paket terlambat.</span>` : null
                  }
                }
                if (item.type === 'pemilihan' && item.detail) {
                  const { selesai, target, terlambat } = item.detail
                  return {
                    label: 'Pemilihan',
                    value: `${selesai}<small>${terlambat > 0 ? `<sup><code class="text-error v-card--variant-elevated">-${terlambat}</code></sup>` : ''}/${target}</small>`,
                    popoverTitle: `${selesai} dari ${target} paket selesai.`,
                    popoverContent:
                      terlambat > 0 ? `<span class="text-error">${terlambat} paket terlambat.</span>` : null
                  }
                }
                if (item.type === 'pengadaan' && item.detail) {
                  const { selesai, target, terlambat } = item.detail
                  return {
                    label: 'Pengadaan',
                    value: `${selesai}<small>${terlambat > 0 ? `<sup><code class="text-error v-card--variant-elevated">-${terlambat}</code></sup>` : ''}/${target}</small>`,
                    popoverTitle: `${selesai} dari ${target} paket selesai.`,
                    popoverContent:
                      terlambat > 0 ? `<span class="text-error">${terlambat} paket terlambat.</span>` : null
                  }
                }
                if (item.type === 'penyerahan' && item.detail) {
                  const { selesai, target, terlambat } = item.detail
                  return {
                    label: 'Penyerahan',
                    value: `${selesai}<small>${terlambat > 0 ? `<sup><code class="text-error v-card--variant-elevated">-${terlambat}</code></sup>` : ''}/${target}</small>`,
                    popoverTitle: `${selesai} dari ${target} paket selesai.`,
                    popoverContent:
                      terlambat > 0 ? `<span class="text-error">${terlambat} paket terlambat.</span>` : null
                  }
                }
                return null
              })
              .filter(Boolean),
            progress: progress_formatted
          }

        case 'fisik':
          return {
            title: 'Persenstase Capaian<br>Realisasi Fisik',
            subtitle: `Januari - ${meta.month_name} ${meta.year}`,
            hintTitle: 'PERSENTASE CAPAIAN REALISASI FISIK',
            hintDescription:
              'Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.<br>Realisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.<br>Nilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:<br>(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan',
            items: items
              .map((item) => {
                if (item.type === 'realisasi') {
                  return { label: 'Realisasi', value: item.formatted || item.value }
                }
                if (item.type === 'target') {
                  return { label: 'Target', value: item.formatted || item.value }
                }
                return null
              })
              .filter(Boolean),
            progress: progress_formatted
          }

        case 'anggaran':
          return {
            title: 'Persenstase Capaian<br>Realisasi Anggaran',
            subtitle: `Januari - ${meta.month_name} ${meta.year}`,
            hintTitle: 'PERSENTASE CAPAIAN REALISASI ANGGARAN',
            hintDescription:
              'Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.<br>Realisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.<br>Nilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:<br>(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan',
            items: items
              .map((item) => {
                if (item.type === 'realisasi') {
                  return { label: 'Realisasi', value: item.formatted || item.value }
                }
                if (item.type === 'target') {
                  return { label: 'Target', value: item.formatted || item.value }
                }
                return null
              })
              .filter(Boolean),
            progress: progress_formatted,
            layout: 'rows'
          }

        case 'kinerja':
          return {
            title: 'Persenstase Capaian<br>Realisasi Kinerja',
            subtitle: `Januari - ${meta.month_name} ${meta.year}`,
            hintTitle: 'PERSENTASE CAPAIAN REALISASI KINERJA',
            hintDescription:
              'Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.<br>Realisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.<br>Nilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:<br>(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan',
            items: items
              .map((item) => {
                if (item.type === 'realisasi') {
                  return { label: 'Realisasi', value: item.formatted || item.value }
                }
                if (item.type === 'target') {
                  return { label: 'Target', value: item.formatted || item.value }
                }
                return null
              })
              .filter(Boolean),
            progress: progress_formatted
          }

        default:
          return null
      }
    })
    .filter(Boolean)
}

export function processRealisasiTahunData(apiData) {
  if (!apiData || !apiData.results || !apiData.results[0] || !apiData.results[0].data) return []

  const meta = apiData.results[0].meta

  return apiData.results[0].data
    .map((category) => {
      const { category: cat, progress_formatted, capaian, items } = category

      switch (cat) {
        case 'barjas':
          return {
            title: 'Progres Tahunan<br>Capaian Barjas',
            subtitle: `per-${meta.month_name} ${meta.year}`,
            hintTitle: 'PROGRES TAHUNAN CAPAIAN BARJAS',
            hintDescription:
              'Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.<br>Realisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.<br>Nilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:<br>(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan',
            items: items
              .map((item) => {
                if (item.type === 'realisasi') {
                  return { label: 'Realisasi', value: item.formatted || item.value }
                }
                if (item.type === 'target') {
                  return { label: 'Target', value: item.formatted || item.value }
                }
                return null
              })
              .filter(Boolean),
            progress: progress_formatted,
            color: getCardColorsByProgress(capaian).bgColor
          }

        case 'fisik':
          return {
            title: 'Progres Tahunan<br>Capaian Fisik',
            subtitle: `per-${meta.month_name} ${meta.year}`,
            hintTitle: 'PROGRES TAHUNAN CAPAIAN FISIK',
            hintDescription:
              'Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.<br>Realisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.<br>Nilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:<br>(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan',
            items: items
              .map((item) => {
                if (item.type === 'realisasi') {
                  return { label: 'Realisasi', value: item.formatted || item.value }
                }
                if (item.type === 'target') {
                  return { label: 'Target', value: item.formatted || item.value }
                }
                return null
              })
              .filter(Boolean),
            progress: progress_formatted,
            color: getCardColorsByProgress(capaian).bgColor
          }

        case 'anggaran':
          return {
            title: 'Progres Tahunan<br>Capaian Anggaran',
            subtitle: `per-${meta.month_name} ${meta.year}`,
            hintTitle: 'PROGRES TAHUNAN CAPAIAN ANGGARAN',
            hintDescription:
              'Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.<br>Realisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.<br>Nilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:<br>(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan',
            items: items
              .map((item) => {
                if (item.type === 'realisasi') {
                  return { label: 'Realisasi', value: item.formatted || item.value }
                }
                if (item.type === 'target') {
                  return { label: 'Target', value: item.formatted || item.value }
                }
                return null
              })
              .filter(Boolean),
            progress: progress_formatted,
            color: getCardColorsByProgress(capaian).bgColor,
            layout: 'rows'
          }

        case 'kinerja':
          return {
            title: 'Progres Tahunan<br>Capaian Kinerja',
            subtitle: `per-${meta.month_name} ${meta.year}`,
            hintTitle: 'PROGRES TAHUNAN CAPAIAN KINERJA',
            hintDescription:
              'Adalah nilai persentase yang menunjukkan jumlah capaian yang sudah tercapai dari sejak awal bulan Januari sampai dengan bulan berjalan.<br>Realisasi capaian yang dihitung adalah realisasi paket yang sedang berprogres dan paket yang sudah selesai sampai dengan pembayaran SP2D.<br>Nilai tersebut didapat dari aplikasi (SIRUP, SIBARASAT, SIPDOK & SIPEKAT) yang data tersebut diolah dan dirumuskan sebagai berikut:<br>(Jumlah paket berprogres + Jumlah paket selesai) / Total paket sampai dengan bulan berjalan',
            items: items
              .map((item) => {
                if (item.type === 'realisasi') {
                  return { label: 'Realisasi', value: item.formatted || item.value }
                }
                if (item.type === 'target') {
                  return { label: 'Target', value: item.formatted || item.value }
                }
                return null
              })
              .filter(Boolean),
            progress: progress_formatted,
            color: getCardColorsByProgress(capaian).bgColor
          }

        default:
          return null
      }
    })
    .filter(Boolean)
}
