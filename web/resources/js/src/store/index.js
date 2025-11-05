import Vue from 'vue'
import Vuex from 'vuex'

// Modules
import app from './app'
import appConfig from './app-config'
import verticalMenu from './vertical-menu'
import rankings from './rankings'

Vue.use(Vuex)

export default new Vuex.Store({
    modules: {
        app,
        appConfig,
        verticalMenu,
        rankings,
    },
    strict: process.env.DEV,
})
