import Vue from 'vue'

// axios
import axios from 'axios'

const axiosIns = axios.create({
    baseURL: process.env.MIX_GIN_API_URL,
    timeout: 1000,
    headers: {
        'Content-Type': 'application/json',
    },
})

Vue.prototype.$http = axiosIns

export default axiosIns
