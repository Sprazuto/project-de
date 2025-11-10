<script setup>
// Vue 3 Composition API with seamless useAuth and useApi integration
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useApi } from '@/composables/useApi'

// Import required Vuetify components
import { VCard, VCardText, VRow, VCol } from 'vuetify/components'

// Import theme components
import { VNodeRenderer } from '@layouts/components/VNodeRenderer'
import { useGenerateImageVariant } from '@core/composable/useGenerateImageVariant'

// Import assets
import authV2LoginIllustrationBorderedDark from '@images/pages/auth-v2-login-illustration-bordered-dark.png'
import authV2LoginIllustrationBorderedLight from '@images/pages/auth-v2-login-illustration-bordered-light.png'
import authV2LoginIllustrationDark from '@images/pages/auth-v2-login-illustration-dark.png'
import authV2LoginIllustrationLight from '@images/pages/auth-v2-login-illustration-light.png'
import authV2MaskDark from '@images/pages/misc-mask-dark.png'
import authV2MaskLight from '@images/pages/misc-mask-light.png'

import { themeConfig } from '@themeConfig'

definePage({
  meta: {
    layout: 'blank',
    unauthenticatedOnly: true
  }
})

// Router
const router = useRouter()
const route = useRoute()

// Composables
const { login, isLoading: authLoading } = useAuth()
const { handleError } = useApi()

// Form data
const form = ref({
  username: '',
  password: '',
  remember: false
})

// Form states
const isPasswordVisible = ref(false)
const formRef = ref(null)
const errorMessage = ref('')

// Validation rules
const rules = {
  required: (value) => !!value || 'This field is required',
  username: (value) => value.length >= 3 || 'Username must be at least 3 characters',
  password: (value) => value.length >= 6 || 'Password must be at least 6 characters'
}

// Get form validation status
const isFormValid = computed(() => {
  return (
    form.value.username &&
    form.value.password &&
    rules.username(form.value.username) === true &&
    rules.password(form.value.password) === true
  )
})

// Computed loading state
const isLoading = computed(() => authLoading.value)

const authThemeImg = useGenerateImageVariant(
  authV2LoginIllustrationLight,
  authV2LoginIllustrationDark,
  authV2LoginIllustrationBorderedLight,
  authV2LoginIllustrationBorderedDark,
  true
)

const authThemeMask = useGenerateImageVariant(authV2MaskLight, authV2MaskDark)

// Authentication handler using the new seamless API pattern
const handleLogin = async () => {
  if (!formRef.value?.validate()) return

  errorMessage.value = ''

  try {
    // Use the integrated useAuth composable for seamless authentication
    const response = await login(null, form.value.username, form.value.password)

    // Success - redirect to intended page or dashboard
    const redirectTo = route.query.redirect || '/'

    await router.push(redirectTo)
  } catch (error) {
    console.error('Login failed:', error)

    // Use centralized error handling from useApi
    errorMessage.value = handleError(error, 'Login failed. Please try again.')
  }
}

// Reset form
const resetForm = () => {
  form.value = {
    username: '',
    password: '',
    remember: false
  }
  errorMessage.value = ''
  formRef.value?.resetValidation()
}
</script>

<template>
  <VRow no-gutters class="auth-wrapper bg-surface">
    <VCol md="8" class="d-none d-md-flex">
      <div class="position-relative bg-background rounded-lg w-100 ma-8 me-0">
        <div class="d-flex align-center justify-center w-100 h-100">
          <VImg max-width="505" :src="authThemeImg" class="auth-illustration mt-16 mb-2" />
        </div>

        <VImg class="auth-footer-mask" :src="authThemeMask" />
      </div>
    </VCol>

    <VCol cols="12" md="4" class="auth-card-v2 d-flex align-center justify-center">
      <VCard flat :max-width="500" class="mt-12 mt-sm-0 pa-4">
        <VCardText>
          <VNodeRenderer :nodes="themeConfig.app.logo" class="mb-6" />
          <h4 class="text-h4 mb-1">
            Welcome to
            <span class="text-capitalize">{{ themeConfig.app.title }}</span
            >! 👋🏻
          </h4>
          <p class="mb-0">Please sign-in to your account and start the adventure</p>
        </VCardText>
        <VCardText>
          <!-- Error Alert -->
          <VAlert
            v-if="errorMessage"
            type="error"
            variant="tonal"
            class="mb-4"
            closable
            @click:close="errorMessage = ''"
          >
            {{ errorMessage }}
          </VAlert>

          <VForm ref="formRef" validate-on="submit" @submit.prevent="handleLogin">
            <VRow>
              <!-- username -->
              <VCol cols="12">
                <AppTextField
                  v-model="form.username"
                  label="Username"
                  placeholder="Enter your username"
                  :rules="[rules.required, rules.username]"
                  :disabled="isLoading"
                  autofocus
                  required
                />
              </VCol>

              <!-- password -->
              <VCol cols="12">
                <AppTextField
                  v-model="form.password"
                  label="Password"
                  placeholder="············"
                  :type="isPasswordVisible ? 'text' : 'password'"
                  :append-inner-icon="isPasswordVisible ? 'tabler-eye-off' : 'tabler-eye'"
                  :rules="[rules.required, rules.password]"
                  :disabled="isLoading"
                  required
                  @click:append-inner="isPasswordVisible = !isPasswordVisible"
                />

                <div class="d-flex align-center flex-wrap justify-space-between mt-2 mb-4">
                  <VCheckbox v-model="form.remember" label="Remember me" :disabled="isLoading" />
                  <a class="text-primary ms-2 mb-1" href="#"> Forgot Password? </a>
                </div>

                <VBtn block type="submit" :disabled="!isFormValid || isLoading" :loading="isLoading">
                  <template #prepend>
                    <VProgressCircular v-if="isLoading" indeterminate size="16" width="2" color="white" />
                  </template>
                  Sign In
                </VBtn>
              </VCol>

              <!-- create account -->
              <VCol cols="12" class="text-center text-base">
                <span>New on our platform?</span>

                <a class="text-primary ms-2" href="#"> Create an account </a>
              </VCol>

              <VCol cols="12" class="d-flex align-center">
                <VDivider />

                <span class="mx-4">or</span>

                <VDivider />
              </VCol>

              <!-- auth providers -->
              <VCol cols="12" class="text-center">
                <AuthProvider />
              </VCol>
            </VRow>
          </VForm>
        </VCardText>
      </VCard>
    </VCol>
  </VRow>
</template>

<style lang="scss">
@use '@core/scss/template/pages/page-auth.scss';
</style>
