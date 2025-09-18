import Vue from 'vue'
import { ToastPlugin, ModalPlugin } from 'bootstrap-vue'
import VueCompositionAPI from '@vue/composition-api'

import router from './router'
import store from './store'
import App from './App.vue'
import axios from 'axios'

import { default as useJwt } from './@core/auth/jwt/useJwt.js'
import AuthService from './services/auth.js'

// Global Components
import './global-components'

// 3rd party plugins
import '@/libs/portal-vue'
import '@/libs/toastification'

// BSV Plugin Registration
Vue.use(ToastPlugin)
Vue.use(ModalPlugin)

// Composition API
Vue.use(VueCompositionAPI)

// import core styles
require('@resources/scss/core.scss')

// import assets styles
require('@resources/assets/scss/style.scss')

Vue.config.productionTip = false

axios.defaults.baseURL = process.env.MIX_GIN_API_URL

// Initialize JWT service with global axios for automatic token refresh
const { jwt: jwtInstance } = useJwt(axios)
Vue.prototype.$jwt = jwtInstance

Vue.prototype.$http = axios

// Auto authenticate before starting the app
AuthService.autoAuth().then(() => {
    new Vue({
        router,
        store,
        render: h => h(App),
    }).$mount('#app')
}).catch((error) => {
    console.error('main.js: AutoAuth failed:', error)
    // Still mount app, but without auth
    new Vue({
        router,
        store,
        render: h => h(App),
    }).$mount('#app')
})
