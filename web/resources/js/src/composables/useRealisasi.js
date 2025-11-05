import Vue from 'vue';

const { set } = Vue;
import { processRealisasiBulanData, processRealisasiTahunData } from '@/utils/realisasiDataProcessor';

export function useRealisasi() {
    const realisasiBulan = Vue.observable([]);
    const realisasiTahun = Vue.observable([]);
    const loading = Vue.observable({
        bulan: false,
        tahun: false,
    });
    const error = Vue.observable({
        bulan: null,
        tahun: null,
    });

    const fetchRealisasiBulan = async (params = {}) => {
        set(loading, 'bulan', true);
        set(error, 'bulan', null);
        try {
            const data = await realisasiService.getRealisasiBulan(params);
            const processed = processRealisasiBulanData(data);
            realisasiBulan.splice(0, realisasiBulan.length, ...processed);
        } catch (err) {
            set(error, 'bulan', err.message || 'Failed to fetch realisasi bulan data');
            console.error('Error fetching realisasi bulan:', err);
        } finally {
            set(loading, 'bulan', false);
        }
    };

    const fetchRealisasiTahun = async (params = {}) => {
        set(loading, 'tahun', true);
        set(error, 'tahun', null);
        try {
            const data = await realisasiService.getRealisasiTahun(params);
            const processed = processRealisasiTahunData(data);
            realisasiTahun.splice(0, realisasiTahun.length, ...processed);
        } catch (err) {
            set(error, 'tahun', err.message || 'Failed to fetch realisasi tahun data');
            console.error('Error fetching realisasi tahun:', err);
        } finally {
            set(loading, 'tahun', false);
        }
    };

    const refreshData = async (params = {}) => {
        await Promise.all([
            fetchRealisasiBulan(params),
            fetchRealisasiTahun(params),
        ]);
    };

    // Initialize data on mount
    refreshData();

    return {
        realisasiBulan,
        realisasiTahun,
        loading,
        error,
        fetchRealisasiBulan,
        fetchRealisasiTahun,
        refreshData,
    };
}
