import { setupLayouts } from 'virtual:generated-layouts'
import { createRouter, createWebHistory } from 'vue-router/auto'

function recursiveLayouts(route) {
  if (route.children) {
    for (let i = 0; i < route.children.length; i++) {
      route.children[i] = recursiveLayouts(route.children[i])
    }

    return route
  }

  return setupLayouts([route])[0]
}

// Define route protection rules BEFORE router creation
function getRouteProtectionConfig(routeName, routePath) {
  const name = routeName?.toString().toLowerCase() || ''
  const path = routePath?.toLowerCase() || ''

  // Login page - only accessible when unauthenticated
  if (name.includes('login') || name === 'auth-login' || path === '/login') {
    return {
      requiresAuth: false,
      unauthenticatedOnly: true
    }
  }

  // Dashboard/Index - requires authentication
  // Check both route name and path for index/dashboard
  if (name === '' || name === 'dashboard' || name.includes('index') || path === '/' || path === '/index' || path === '/dashboard') {
    return {
      requiresAuth: true,
      unauthenticatedOnly: false
    }
  }

  // All other routes default to protected
  else {
    return {
      requiresAuth: true,
      unauthenticatedOnly: false
    }
  }
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to) {
    if (to.hash) return { el: to.hash, behavior: 'smooth', top: 60 }

    return { top: 0 }
  },
  extendRoutes: (pages) => {
    const routes = [...pages].map((route) => {
      // Set meta properties BEFORE route is used
      const protection = getRouteProtectionConfig(route.name, route.path)

      route.meta = {
        ...route.meta,
        ...protection
      }

      return recursiveLayouts(route)
    })

    return routes
  }
})

// Pre-import auth composable to avoid timing issues
let auth = null

const initAuth = async () => {
  if (!auth) {
    const { useAuth } = await import('@/composables/useAuth')

    auth = useAuth()
  }

  return auth
}

// Authentication Middleware - now with immediate access to auth state
router.beforeEach(async (to, from, next) => {
  // Initialize auth composable once
  const authInstance = await initAuth()

  // Check if route requires authentication
  if (to.meta.requiresAuth) {
    // Protected route - requires authentication
    const isAuthValid = await authInstance.checkAuthStatus()

    if (!isAuthValid || !authInstance.isAuthenticated.value) {
      // User is not authenticated, redirect to login
      const redirectQuery = to.fullPath !== '/' ? { redirect: to.fullPath } : {}

      next({
        path: '/login',
        query: redirectQuery
      })

      return
    }
  } else if (to.meta.unauthenticatedOnly) {
    // Public route that should only be accessible when not authenticated
    const isAuthValid = await authInstance.checkAuthStatus()

    if (isAuthValid && authInstance.isAuthenticated.value) {
      // User is authenticated but trying to access login/register pages
      // Redirect to dashboard or specified redirect
      const redirectTo = to.query.redirect || from.query.redirect || '/'

      next(redirectTo)

      return
    }
  }

  // Allow navigation
  next()
})

export { router }
export default function (app) {
  app.use(router)
}
