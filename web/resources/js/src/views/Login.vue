<template>
  <div class="auth-wrapper auth-v1 px-2">
    <div class="auth-inner py-2">

      <!-- Login v1 -->
      <b-card class="mb-0">
        <b-link class="brand-logo">
          <vuexy-logo />

          <h2 class="brand-text text-primary ml-1">
            Vuexy
          </h2>
        </b-link>

        <b-card-title class="mb-1">
          Welcome to Vuexy! 👋
        </b-card-title>
        <b-card-text class="mb-2">
          Please sign-in to your account and start the adventure
        </b-card-text>

        <!-- form -->
        <validation-observer
          ref="loginForm"
          #default="{invalid}"
        >
          <b-form
            class="auth-login-form mt-2"
            @submit.prevent="login"
          >

            <!-- username or email -->
            <b-form-group
              label-for="login-field"
              label="Username or Email"
            >
              <validation-provider
                #default="{ errors }"
                name="Login Field"
                rules="required"
              >
                <b-form-input
                  id="login-field"
                  v-model="loginField"
                  name="login-field"
                  :state="errors.length > 0 ? false:null"
                  placeholder="Username or email"
                  autofocus
                />
                <small class="text-danger">{{ errors[0] }}</small>
              </validation-provider>
            </b-form-group>

            <!-- password -->
            <b-form-group>
              <div class="d-flex justify-content-between">
                <label for="password">Password</label>
                <b-link :to="{name:'auth-forgot-password-v1'}">
                  <small>Forgot Password?</small>
                </b-link>
              </div>
              <validation-provider
                #default="{ errors }"
                name="Password"
                rules="required"
              >
                <b-input-group
                  class="input-group-merge"
                  :class="errors.length > 0 ? 'is-invalid':null"
                >
                  <b-form-input
                    id="password"
                    v-model="password"
                    :type="passwordFieldType"
                    class="form-control-merge"
                    :state="errors.length > 0 ? false:null"
                    name="login-password"
                    placeholder="Password"
                  />

                  <b-input-group-append is-text>
                    <feather-icon
                      class="cursor-pointer"
                      :icon="passwordToggleIcon"
                      @click="togglePasswordVisibility"
                    />
                  </b-input-group-append>
                </b-input-group>
                <small class="text-danger">{{ errors[0] }}</small>
              </validation-provider>
            </b-form-group>

            <!-- checkbox -->
            <b-form-group>
              <b-form-checkbox
                id="remember-me"
                v-model="status"
                name="checkbox-1"
              >
                Remember Me
              </b-form-checkbox>
            </b-form-group>

            <!-- submit button -->
            <b-button
              variant="primary"
              type="submit"
              block
              :disabled="invalid"
            >
              Sign in
            </b-button>
          </b-form>
        </validation-observer>

        <b-card-text class="text-center mt-2">
          <span>New on our platform? </span>
          <b-link :to="{name:'auth-register-v1'}">
            <span>Create an account</span>
          </b-link>
        </b-card-text>

        <div class="divider my-2">
          <div class="divider-text">
            or
          </div>
        </div>

        <!-- social button -->
        <div class="auth-footer-btn d-flex justify-content-center">
          <b-button
            href="javascript:void(0)"
            variant="facebook"
          >
            <feather-icon icon="FacebookIcon" />
          </b-button>
          <b-button
            href="javascript:void(0)"
            variant="twitter"
          >
            <feather-icon icon="TwitterIcon" />
          </b-button>
          <b-button
            href="javascript:void(0)"
            variant="google"
          >
            <feather-icon icon="MailIcon" />
          </b-button>
          <b-button
            href="javascript:void(0)"
            variant="github"
          >
            <feather-icon icon="GithubIcon" />
          </b-button>
        </div>
      </b-card>
      <!-- /Login v1 -->
    </div>
  </div>
</template>

<script>
import { ValidationProvider, ValidationObserver } from 'vee-validate'
import {
  BButton,
  BForm,
  BFormInput,
  BFormGroup,
  BCard,
  BLink,
  BCardTitle,
  BCardText,
  BInputGroup,
  BInputGroupAppend,
  BFormCheckbox,
} from 'bootstrap-vue'
import VuexyLogo from '@core/layouts/components/Logo.vue'
import { togglePasswordVisibility } from '@core/mixins/ui/forms'
import AuthService from '@/services/auth.js'

export default {
  components: {
    // BSV
    BButton,
    BForm,
    BFormInput,
    BFormGroup,
    BCard,
    BCardTitle,
    BLink,
    VuexyLogo,
    BCardText,
    BInputGroup,
    BInputGroupAppend,
    BFormCheckbox,
    ValidationProvider,
    ValidationObserver,
  },
  mixins: [togglePasswordVisibility],
  data() {
    return {
      loginField: '',
      password: '',
      status: '',
    }
  },
  computed: {
    passwordToggleIcon() {
      return this.passwordFieldType === 'password' ? 'EyeIcon' : 'EyeOffIcon'
    },
  },
  methods: {
    async login() {
      const payload = {
        email: this.loginField.includes('@') ? this.loginField : '',
        username: !this.loginField.includes('@') ? this.loginField : '',
        password: this.password,
      }

      try {
        await AuthService.login(payload.email, payload.username, payload.password)
        this.$router.push({ name: 'home' })
      } catch (error) {
        let errorText = String(error.response?.data?.message || error.message || (error.response?.status ? `Server error (status ${error.response.status})` : 'Network error. Please check your connection.'));
        if (!errorText) errorText = 'An error occurred during login. Please try again.';
        this.$toast({
          component: 'ToastificationContent',
          props: {
            title: 'Login Failed',
            icon: 'CoffeeIcon',
            variant: 'danger',
            text: errorText,
          },
        })
      }
    },
  },
}
</script>

<style lang="scss">
@import '~@resources/scss/vue/pages/page-auth.scss';
</style>
