module.exports = {
  env: {
    browser: true,
    es2021: true,
    node: true
  },
  extends: ['eslint:recommended', 'plugin:vue/vue3-recommended', 'plugin:prettier/recommended'],
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 13,
    sourceType: 'module'
  },
  plugins: ['vue'],
  ignorePatterns: ['src/plugins/iconify/*.js', 'src/assets/**', 'node_modules', 'dist', 'public', '*.d.ts'],
  rules: {
    // Disable all rules for maximum performance
    'no-console': 'off',
    'no-debugger': 'off',
    'no-unused-vars': 'off',
    'no-undef': 'off',
    'no-shadow': 'off',
    'max-len': 'off',
    indent: 'off',
    'vue/html-indent': 'off',
    camelcase: 'off',
    'arrow-parens': 'off',
    semi: 'off',
    'comma-dangle': 'off',
    'object-curly-spacing': 'off',

    // Vue specific - disable most rules
    'vue/multi-word-component-names': 'off',
    'vue/require-default-prop': 'off',
    'vue/no-restricted-class': 'off',
    'vue/valid-v-slot': 'off',
    'vue/no-v-html': 'off',

    // Disable import rules
    'import/no-unresolved': 'off',
    'import/extensions': 'off',
    'import/newline-after-import': 'off',

    // Disable all slow plugins
    'sonarjs/cognitive-complexity': 'off',
    'sonarjs/no-duplicate-string': 'off',
    'sonarjs/no-nested-template-literals': 'off'
  },
  globals: {
    // Vue 3 Composition API globals
    ref: 'readonly',
    computed: 'readonly',
    watch: 'readonly',
    watchEffect: 'readonly',
    reactive: 'readonly',
    toRefs: 'readonly',
    toRef: 'readonly',
    unref: 'readonly',
    isRef: 'readonly',
    h: 'readonly',
    Fragment: 'readonly',
    isProxy: 'readonly',
    isReactive: 'readonly',
    isReadonly: 'readonly',
    markRaw: 'readonly',
    shallowReactive: 'readonly',
    shallowRef: 'readonly',
    shallowReadonly: 'readonly',
    triggerRef: 'readonly',
    customRef: 'readonly',
    effect: 'readonly',
    effectScope: 'readonly',
    getCurrentScope: 'readonly',
    onScopeDispose: 'readonly',
    defineComponent: 'readonly',
    defineAsyncComponent: 'readonly',
    useAttrs: 'readonly',
    useSlots: 'readonly',
    useRoute: 'readonly',
    useRouter: 'readonly',
    useI18n: 'readonly',
    useCurrentInstance: 'readonly',
    getCurrentInstance: 'readonly',
    useMagicKeys: 'readonly',
    onMounted: 'readonly',
    onUpdated: 'readonly',
    onUnmounted: 'readonly',
    onBeforeMount: 'readonly',
    onBeforeUpdate: 'readonly',
    onBeforeUnmount: 'readonly',
    onErrorCaptured: 'readonly',
    onActivated: 'readonly',
    onDeactivated: 'readonly',
    onRenderTriggered: 'readonly',
    onRenderTracked: 'readonly',
    provide: 'readonly',
    inject: 'readonly',
    nextTick: 'readonly',
    useHead: 'readonly',
    useMouse: 'readonly',
    useWindowSize: 'readonly',
    useWindowScroll: 'readonly',
    useEventListener: 'readonly',
    useElementHover: 'readonly',
    useMediaQuery: 'readonly',
    useToggle: 'readonly',
    useSkins: 'readonly',
    useThemeConfig: 'readonly',
    useTheme: 'readonly',
    useGenerateImageVariant: 'readonly',
    useDebounceFn: 'readonly',
    useThrottleFn: 'readonly',
    useCookie: 'readonly',
    useStorage: 'readonly',
    useLocalStorage: 'readonly',
    useSessionStorage: 'readonly',
    useFetch: 'readonly',
    useApi: 'readonly',
    useAuth: 'readonly',
    useErrorHandler: 'readonly',
    handleError: 'readonly',
    definePage: 'readonly',
    createUrl: 'readonly',
    useResponsiveSidebar: 'readonly',
    syncRef: 'readonly',
    until: 'readonly',
    useMounted: 'readonly',
    useWindow: 'readonly'
  },
  settings: {
    'import/resolver': {
      node: {
        extensions: ['.js', '.jsx', '.ts', '.tsx', '.vue']
      },
      'eslint-import-resolver-custom-alias': {
        alias: {
          '@': './src',
          '@core': './src/@core',
          '@layouts': './src/@layouts'
        },
        extensions: ['.ts', '.js', '.tsx', '.jsx', '.vue']
      }
    }
  }
}
