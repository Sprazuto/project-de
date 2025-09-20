import Vue from 'vue'
import { ToastPlugin, ModalPlugin } from 'bootstrap-vue'
import VueCompositionAPI from '@vue/composition-api'

import router from './router'
import store from './store'
import App from './App.vue'
import axios from '@axios'

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

// VeeValidate
import { extend } from 'vee-validate'
import { required } from 'vee-validate/dist/rules'
import { ValidationObserver, ValidationProvider } from 'vee-validate'

extend('required', required)

// Install components globally
Vue.component('ValidationObserver', ValidationObserver)
Vue.component('ValidationProvider', ValidationProvider)

// import core styles
/* require('@/resources/scss/core.scss') */

// import assets styles
/* require('@/resources/assets/scss/style.scss') */

Vue.config.productionTip = false

// ToastificationContent
import ToastificationContent from '@core/components/toastification/ToastificationContent.vue'
Vue.component('ToastificationContent', ToastificationContent)

// Initialize JWT service with global axios for automatic token refresh
const { jwt: jwtInstance } = useJwt(axios)

Vue.prototype.$jwt = jwtInstance

Vue.prototype.$http = axios

// Mount Vue to #app
new Vue({
    router,
    store,
    render: h => h(App),
}).$mount('#app')
