import { deepMerge } from '@antfu/utils'
import { createVuetify } from 'vuetify'
import { VBtn } from 'vuetify/components/VBtn'
import { VSkeletonLoader } from 'vuetify/lib/labs/VSkeletonLoader/index.mjs'
import defaults from './defaults'
import { icons } from './icons'
import { staticPrimaryColor, themes } from './theme'

// Styles
import { cookieRef } from '@/@layouts/stores/config'
import '@core/scss/template/libs/vuetify/index.scss'
import 'vuetify/styles'

export default function (app) {
  const cookieThemeValues = {
    defaultTheme: resolveVuetifyTheme(),
    themes: {
      light: {
        colors: {
          primary: cookieRef('lightThemePrimaryColor', staticPrimaryColor).value
        }
      },
      dark: {
        colors: {
          primary: cookieRef('darkThemePrimaryColor', staticPrimaryColor).value
        }
      }
    }
  }

  const optionTheme = deepMerge({ themes }, cookieThemeValues)

  const vuetify = createVuetify({
    aliases: {
      IconBtn: VBtn
    },
    components: {
      VSkeletonLoader
    },
    defaults,
    icons,
    theme: optionTheme
  })

  app.use(vuetify)
}
