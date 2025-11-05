export default {
    namespaced: true,
    state: {
        rankings: [],
        loading: false,
        error: null,
    },
    getters: {
        getRankings: (state) => state.rankings,
        isLoading: (state) => state.loading,
        getError: (state) => state.error,
    },
    mutations: {
        SET_RANKINGS(state, rankings) {
            state.rankings = rankings;
        },
        SET_LOADING(state, loading) {
            state.loading = loading;
        },
        SET_ERROR(state, error) {
            state.error = error;
        },
    },
    actions: {
        async fetchRankings({ commit }) {
            commit('SET_LOADING', true);
            commit('SET_ERROR', null);
            try {
                // Mock data - replace with actual API call
                const mockRankings = [
                    {
                        name: "Kabupaten A",
                        total_score: 85.75,
                        categories: [
                            {
                                title: "Realisasi Barjas",
                                subtitle: "Capaian Tahunan",
                                percentage: 85,
                                icon: "LayersIcon",
                            },
                            {
                                title: "Realisasi Fisik",
                                subtitle: "Capaian Tahunan",
                                percentage: 78,
                                icon: "MapPinIcon",
                            },
                            {
                                title: "Realisasi Anggaran",
                                subtitle: "Capaian Tahunan",
                                percentage: 92,
                                icon: "ShoppingBagIcon",
                            },
                            {
                                title: "Realisasi Kinerja",
                                subtitle: "Capaian Tahunan",
                                percentage: 88,
                                icon: "SlidersIcon",
                            },
                        ],
                    },
                    {
                        name: "Kabupaten B",
                        total_score: 82.75,
                        categories: [
                            {
                                title: "Realisasi Barjas",
                                subtitle: "Capaian Tahunan",
                                percentage: 82,
                                icon: "LayersIcon",
                            },
                            {
                                title: "Realisasi Fisik",
                                subtitle: "Capaian Tahunan",
                                percentage: 75,
                                icon: "MapPinIcon",
                            },
                            {
                                title: "Realisasi Anggaran",
                                subtitle: "Capaian Tahunan",
                                percentage: 89,
                                icon: "ShoppingBagIcon",
                            },
                            {
                                title: "Realisasi Kinerja",
                                subtitle: "Capaian Tahunan",
                                percentage: 85,
                                icon: "SlidersIcon",
                            },
                        ],
                    },
                    {
                        name: "Kabupaten C",
                        total_score: 79.75,
                        categories: [
                            {
                                title: "Realisasi Barjas",
                                subtitle: "Capaian Tahunan",
                                percentage: 79,
                                icon: "LayersIcon",
                            },
                            {
                                title: "Realisasi Fisik",
                                subtitle: "Capaian Tahunan",
                                percentage: 72,
                                icon: "MapPinIcon",
                            },
                            {
                                title: "Realisasi Anggaran",
                                subtitle: "Capaian Tahunan",
                                percentage: 86,
                                icon: "ShoppingBagIcon",
                            },
                            {
                                title: "Realisasi Kinerja",
                                subtitle: "Capaian Tahunan",
                                percentage: 82,
                                icon: "SlidersIcon",
                            },
                        ],
                    },
                    {
                        name: "Kabupaten D",
                        total_score: 76.75,
                        categories: [
                            {
                                title: "Realisasi Barjas",
                                subtitle: "Capaian Tahunan",
                                percentage: 76,
                                icon: "LayersIcon",
                            },
                            {
                                title: "Realisasi Fisik",
                                subtitle: "Capaian Tahunan",
                                percentage: 69,
                                icon: "MapPinIcon",
                            },
                            {
                                title: "Realisasi Anggaran",
                                subtitle: "Capaian Tahunan",
                                percentage: 83,
                                icon: "ShoppingBagIcon",
                            },
                            {
                                title: "Realisasi Kinerja",
                                subtitle: "Capaian Tahunan",
                                percentage: 79,
                                icon: "SlidersIcon",
                            },
                        ],
                    },
                    {
                        name: "Kabupaten E",
                        total_score: 73.75,
                        categories: [
                            {
                                title: "Realisasi Barjas",
                                subtitle: "Capaian Tahunan",
                                percentage: 73,
                                icon: "LayersIcon",
                            },
                            {
                                title: "Realisasi Fisik",
                                subtitle: "Capaian Tahunan",
                                percentage: 66,
                                icon: "MapPinIcon",
                            },
                            {
                                title: "Realisasi Anggaran",
                                subtitle: "Capaian Tahunan",
                                percentage: 80,
                                icon: "ShoppingBagIcon",
                            },
                            {
                                title: "Realisasi Kinerja",
                                subtitle: "Capaian Tahunan",
                                percentage: 76,
                                icon: "SlidersIcon",
                            },
                        ],
                    },
                ];
                commit('SET_RANKINGS', mockRankings);
            } catch (error) {
                commit('SET_ERROR', error.message || 'Failed to fetch rankings');
                console.error('Error fetching rankings:', error);
            } finally {
                commit('SET_LOADING', false);
            }
        },
    },
};