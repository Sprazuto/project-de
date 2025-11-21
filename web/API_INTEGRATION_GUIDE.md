# API Integration Documentation

This document outlines the comprehensive API integration architecture within the SIJAGUR (Sistem Informasi Realisasi Anggaran) Vue 3 frontend application, focusing on seamless communication with the Go/Gin backend API.

## Overview

The SIJAGUR frontend implements a sophisticated API integration system that provides:

- **Automatic JWT authentication** with token management and refresh
- **Centralized error handling** with user-friendly messages
- **Real-time dashboard data management** for procurement performance tracking
- **Reactive state management** with Vue 3 Composition API
- **Type-safe API communication** with consistent request/response patterns

The integration serves a government procurement monitoring system with complex data visualization requirements, including monthly/yearly realization tracking, performance rankings, and geospatial data presentation.

## Architecture

### Core Composables

#### 1. `useAuth.js` - Authentication State Management

- **Location**: `src/composables/useAuth.js`
- **Purpose**: Manages user authentication state and provides token management methods
- **Key Features**:
  - Global reactive authentication state with Vue 3 Composition API
  - JWT token management (access + refresh tokens)
  - User role and permission checking
  - Route protection middleware with Vue Router integration
  - Session persistence using localStorage
  - Automatic token validation and expiration handling
  - Computed properties for user information (userName, userEmail, userRole)

#### 2. `useApi.js` - HTTP Client with Auth Integration

- **Location**: `src/composables/useApi.js`
- **Purpose**: Provides authenticated API requests with seamless token handling using ofetch
- **Key Features**:
  - Automatic JWT token injection from `useAuth` composable
  - Request/response interceptors for authentication and error handling
  - Network status detection (online/offline)
  - Consistent error handling with user-friendly messages
  - Support for all HTTP methods (GET, POST, PUT, PATCH, DELETE)
  - File upload capabilities
  - Request counting for loading states

#### 3. `useDashboard.js` - Dashboard Data Orchestration

- **Location**: `src/composables/useDashboard.js`
- **Purpose**: Orchestrates complex dashboard data fetching and state management
- **Key Features**:
  - Coordinates multiple API endpoints simultaneously
  - Manages loading states for different data sections
  - Implements auto-refresh functionality with configurable intervals
  - Processes and transforms API responses for frontend consumption
  - Handles data relationships between different API endpoints
  - Provides retry mechanisms for failed requests
  - Filters and parameters management (year, month, satker)

#### 4. `useErrorHandler.js` - Centralized Error Management

- **Location**: `src/composables/useErrorHandler.js`
- **Purpose**: Provides consistent error handling and user feedback across the application
- **Key Features**:
  - HTTP status code classification (401, 403, 422, 500, etc.)
  - User-friendly error message generation
  - Error logging and reporting
  - Automatic authentication error handling with logout
  - Network error detection and handling

## Integration Flow

### Authentication Flow

1. User enters credentials in `src/pages/login.vue`
2. `useAuth.login()` calls backend `POST /v1/user/login`
3. JWT tokens (access + refresh) are stored in localStorage
4. `useAuth` updates reactive authentication state
5. `useApi` automatically injects tokens into all subsequent requests
6. Route protection middleware validates authentication status
7. Token expiration triggers automatic logout and redirect

### Dashboard Data Flow

1. User accesses `src/pages/index.vue` (dashboard)
2. `useDashboard` initializes and calls multiple API endpoints simultaneously:
   - `GET /realisasi-bulan` - Monthly realization data
   - `GET /realisasi-tahun` - Yearly realization data
   - `GET /realisasi-perbulan` - Monthly breakdown data
   - `GET /sijagur/peringkat-kinerja` - Performance rankings (OPD & Kecamatan)
3. Loading states are managed per data section
4. Data is processed and transformed for chart components
5. Auto-refresh runs periodically (configurable interval)

### API Request Flow

1. Component calls API method (e.g., `api.get('/realisasi-bulan')`)
2. `useApi` retrieves current access token from `useAuth`
3. Request interceptor adds `Authorization: Bearer <token>` header
4. Request is sent with 30-second timeout using ofetch
5. Response interceptor handles success/error cases
6. Authentication errors (401) trigger automatic logout
7. Other errors are processed through centralized error handling
8. Network errors are detected and user-friendly messages displayed

## Usage Patterns

### Dashboard Data Management

```javascript
// src/pages/index.vue - Main dashboard implementation
import { computed, onMounted } from 'vue'
import { useDashboard } from '@/composables/useDashboard'

export default {
  setup() {
    const dashboard = useDashboard()

    // Processed data for components
    const processedRealisasiBulan = computed(() => {
      return dashboard.realisasiBulan.value || []
    })

    const rankingsOPD = computed(() => {
      return dashboard.rankingsOpd.value || []
    })

    // Auto-initialized on mount
    onMounted(async () => {
      // Dashboard handles all API calls automatically
      await dashboard.refreshData()
    })

    return {
      processedRealisasiBulan,
      rankingsOPD,
      loading: dashboard.loading,
      error: dashboard.error
    }
  }
}
```

### Authentication Integration

```javascript
// src/pages/login.vue - Login implementation
import { useAuth } from '@/composables/useAuth'
import { useRouter } from 'vue-router'

export default {
  setup() {
    const { login, isLoading } = useAuth()
    const router = useRouter()

    const handleLogin = async (credentials) => {
      try {
        await login(credentials.email, credentials.username, credentials.password)
        await router.push('/') // Redirect to dashboard
      } catch (error) {
        // Error handled by useAuth
        console.error('Login failed:', error.message)
      }
    }

    return {
      handleLogin,
      isLoading
    }
  }
}
```

### API Integration in Components

```javascript
// Example component using useApi directly
import { useApi } from '@/composables/useApi'

export default {
  setup() {
    const { get, post } = useApi()

    const fetchArticles = async () => {
      const response = await get('/articles')
      if (response.success) {
        return response.data.results[0]?.data || []
      }
      return []
    }

    const createArticle = async (articleData) => {
      const response = await post('/article', articleData)
      return response.success
    }

    return {
      fetchArticles,
      createArticle
    }
  }
}
```

### Route Protection

```javascript
// src/composables/useAuth.js - Route middleware
const authMiddleware = async (to, from, next) => {
  const token = AuthService.getAccessToken()
  const user = AuthService.getCurrentUser()

  if (token && user) {
    currentUser.value = user
    isAuthenticated.value = true
    next()
  } else {
    next(`/login?redirect=${encodeURIComponent(to.fullPath)}`)
  }
}
```

## Error Handling

### Automatic Error Classification

The `useErrorHandler` composable automatically classifies errors:

- **401 (Authentication)**: Automatic logout and redirect to login
- **403 (Authorization)**: Display permission error message
- **422 (Validation)**: Show validation error details
- **500+ (Server)**: Display server error message
- **Network**: Show connection error message

### Error Handling Example

```javascript
import { useErrorHandler } from '@/composables/useErrorHandler'

export default {
  setup() {
    const { setError, getUserFriendlyMessage } = useErrorHandler()

    const handleApiError = (error) => {
      const userMessage = getUserFriendlyMessage(error)
      setError(error) // Also logs the error

      return userMessage
    }

    return {
      handleApiError
    }
  }
}
```

## Token Management

### Automatic Token Handling

- Tokens are automatically retrieved from `useAuth`
- No manual token management required
- Automatic token injection in API requests
- Proper cleanup on authentication failure

### Token Methods Available

```javascript
const {
  getAccessToken, // Get current access token
  getRefreshToken, // Get refresh token (if available)
  isTokenExpired, // Check if token is expired
  hasValidToken, // Check if user has valid token
  saveAuthState, // Save auth state
  clearAuthState // Clear all auth data
} = useAuth()
```

## Migration Benefits

### Before Integration

- Manual token management in each component
- Inconsistent error handling
- Duplicate authentication logic
- Token storage in localStorage directly
- No centralized error logging

### After Integration

- ✅ Automatic token management
- ✅ Centralized error handling
- ✅ Consistent authentication flow
- ✅ Seamless API integration
- ✅ Automatic logout on token expiration
- ✅ Centralized error logging and reporting

## Best Practices

### 1. Always use the composables together

```javascript
// ✅ Good
const { isAuthenticated } = useAuth()
const { get } = useApi()

// ❌ Bad - mixing direct API calls
const token = localStorage.getItem('token')
```

### 2. Use centralized error handling

```javascript
// ✅ Good
import { useErrorHandler } from '@/composables/useErrorHandler'
const { handleError } = useErrorHandler()

// ❌ Bad - manual error handling
try {
  const response = await fetch('/api/endpoint')
} catch (error) {
  console.error('API error:', error)
}
```

### 3. Leverage authentication state

```javascript
// ✅ Good
if (!isAuthenticated.value) {
  return
}

// ❌ Bad - assuming authentication
const response = await get('/api/data')
```

## Testing Integration

### Test Authentication Flow

```javascript
// In browser console or test environment
import { quickApiCheck } from '@/composables/useApiIntegrationTest.js'

// Run integration check
quickApiCheck().then((result) => {
  console.log('Integration working:', result)
})
```

### Manual Testing Checklist

- [ ] User can log in successfully
- [ ] API requests include authentication headers
- [ ] Invalid tokens trigger automatic logout
- [ ] Error messages are user-friendly
- [ ] Network errors are handled gracefully
- [ ] Protected routes require authentication

## Configuration

### Environment Variables

```env
# API Configuration
VITE_API_BASE_URL=http://localhost:8080/api
VITE_API_URL=http://localhost:8080/api

# Optional: API timeout (default: 30000ms)
VITE_API_TIMEOUT=30000
```

### API Endpoints Integration

The frontend integrates with the following backend API endpoints:

#### Authentication & User Management

- `POST /v1/user/login` - User authentication with JWT tokens
- `POST /v1/user/register` - User registration
- `GET /v1/user/logout` - User logout (token invalidation)
- `GET /v1/user/profile` - Get current user profile
- `POST /v1/user/forgot-password` - Password reset initiation
- `POST /v1/user/assign-role` - Assign user roles (admin only)
- `POST /v1/permission/create` - Create permissions (admin only)
- `POST /v1/token/refresh` - Refresh JWT access token

#### Article Management

- `POST /v1/article` - Create new article (write_article permission)
- `GET /v1/articles` - List all articles (read_article permission)
- `GET /v1/article/{id}` - Get specific article (read_article permission)
- `PUT /v1/article/{id}` - Update article (write_article permission)
- `DELETE /v1/article/{id}` - Delete article (write_article permission)

#### SIJAGUR Data Management (Government Procurement)

- `GET /v1/realisasi-bulan` - Monthly realization data with parameters:
  - `tahun` (year), `bulan` (month), `idsatker` (organization ID)
- `GET /v1/realisasi-tahun` - Yearly realization data
- `GET /v1/realisasi-perbulan` - Monthly breakdown data for charts
- `GET /v1/sijagur/peringkat-kinerja` - Performance rankings with filters:
  - `year`, `month`, `idsatker`, `category`, `dimension`, `scope`

#### Data Categories

The SIJAGUR system tracks four main performance categories:

- **Barjas** (Barang dan Jasa) - Goods and services procurement
- **Fisik** - Physical progress indicators
- **Anggaran** - Budget realization
- **Kinerja** - Overall performance metrics

#### Response Patterns

All API responses follow consistent patterns:

```javascript
// Success response
{
  "success": true,
  "data": { /* actual data */ }
}

// Error response
{
  "success": false,
  "error": "User-friendly error message"
}

// SIJAGUR data responses
{
  "results": [
    {
      "data": [/* processed data array */],
      "meta": { /* pagination/metadata */ }
    }
  ]
}
```

## Troubleshooting

### Common Issues & Solutions

#### 1. **Authentication Token Issues**

**Problem**: API requests failing with 401 errors

```javascript
// Check token status
const { getAccessToken, isAuthenticated } = useAuth()
console.log('Token:', getAccessToken())
console.log('Authenticated:', isAuthenticated.value)
```

**Solutions**:

- Verify user is logged in before making requests
- Check token expiration (30-minute default)
- Clear localStorage and re-login if tokens are corrupted
- Ensure backend `/v1/token/refresh` endpoint is accessible

#### 2. **Dashboard Data Not Loading**

**Problem**: Loading states stuck, no data displayed

```javascript
// Debug dashboard state
const dashboard = useDashboard()
console.log('Loading states:', dashboard.loading.value)
console.log('Errors:', dashboard.error.value)
```

**Solutions**:

- Check network connectivity to backend
- Verify API endpoints are responding (check browser network tab)
- Ensure user has required permissions for data access
- Check for CORS issues in development

#### 3. **Vue Router Navigation Issues**

**Problem**: Redirect loops or authentication middleware failures

```javascript
// Check route guards
const { checkAuthStatus } = useAuth()
checkAuthStatus().then((result) => console.log('Auth status:', result))
```

**Solutions**:

- Verify route meta configuration (`requiresAuth: true`)
- Check Vue Router setup in `vite.config.js`
- Ensure authentication state is properly initialized

#### 4. **Performance Issues**

**Problem**: Slow dashboard loading or excessive API calls

```javascript
// Monitor API requests
const { requestCount, isOnline } = useApi()
console.log('Active requests:', requestCount.value)
console.log('Network status:', isOnline.value)
```

**Solutions**:

- Implement request debouncing for user inputs
- Use computed properties to avoid unnecessary recalculations
- Check for memory leaks in component unmounting
- Optimize chart rendering with virtual scrolling

#### 5. **Data Transformation Errors**

**Problem**: Charts not displaying data correctly

```javascript
// Debug data processing
const dashboard = useDashboard()
console.log('Raw data:', dashboard.realisasiBulan.value)
console.log('Processed data:', processedRealisasiBulan.value)
```

**Solutions**:

- Verify data structure matches expected format
- Check for null/undefined values in API responses
- Ensure data processing functions handle edge cases
- Validate chart component props

#### 6. **Build/Development Issues**

**Problem**: Vite dev server or build failures

```bash
# Check environment variables
echo $VITE_API_BASE_URL

# Verify dependencies
pnpm list --depth=0
```

**Solutions**:

- Ensure all environment variables are set
- Check Node.js and pnpm versions
- Clear node_modules and reinstall if needed
- Verify Vite configuration for auto-imports

### Debug Commands

#### Browser Console Debugging

```javascript
// Check authentication state
import { useAuth } from '@/composables/useAuth'
const auth = useAuth()
console.log('Auth state:', {
  authenticated: auth.isAuthenticated.value,
  user: auth.currentUser.value,
  token: auth.getAccessToken()
})

// Test API connectivity
import { useApi } from '@/composables/useApi'
const api = useApi()
api.get('/user/profile').then((result) => console.log('API test:', result))

// Check dashboard data
import { useDashboard } from '@/composables/useDashboard'
const dashboard = useDashboard()
console.log('Dashboard state:', dashboard.loading.value, dashboard.error.value)
```

#### Network Debugging

- Open browser DevTools → Network tab
- Filter by XHR/Fetch requests
- Check response status codes and payloads
- Verify authentication headers are present
- Monitor request timing and sizes

## Future Enhancements

### Identified Improvement Areas

Based on the codebase analysis, the following enhancements would improve the API integration:

#### 1. **Token Refresh Mechanism**

- Implement automatic token refresh before expiration
- Add refresh token rotation for enhanced security
- Handle concurrent refresh requests gracefully

#### 2. **API Response Caching**

- Implement intelligent caching for frequently accessed data
- Add cache invalidation strategies
- Consider service worker caching for offline support

#### 3. **Request Optimization**

- Implement request deduplication
- Add request batching for multiple similar calls
- Optimize payload sizes and compression

#### 4. **Real-time Data Updates**

- WebSocket integration for live dashboard updates
- Server-sent events for performance metrics
- Real-time notifications for system alerts

#### 5. **Enhanced Error Handling**

- Add retry logic with exponential backoff
- Implement circuit breaker pattern
- Add error reporting and monitoring

#### 6. **Performance Monitoring**

- API response time tracking
- Error rate monitoring
- User experience metrics

### Extension Points

#### Custom Integrations

- Additional authentication providers (OAuth, SAML)
- Third-party API integrations
- Custom data transformers and serializers

#### Developer Experience

- API client code generation
- Enhanced debugging tools
- Development proxy configurations

#### Advanced Features

- GraphQL integration capabilities
- File upload progress tracking
- Advanced caching strategies with Redis integration

## Data Processing & Visualization

### SIJAGUR Data Flow

The dashboard implements a sophisticated data processing pipeline:

#### 1. **API Data Fetching**

```javascript
// useDashboard.js - Parallel API calls
await Promise.allSettled([
  fetchRealisasiBulan(),
  fetchRealisasiTahun(),
  fetchRealisasiPerbulan(),
  fetchRankingsOpd(),
  fetchRankingsKecamatan()
])
```

#### 2. **Data Transformation**

```javascript
// Data processing functions
const processRealisasiBulanData = (response) => {
  // Transform raw API data for chart components
  return response.results?.[0]?.data || []
}

const mapRankingRowToCard = (row) => ({
  name: row.nama_opd,
  total_score: row.score_total,
  categories: [
    { key: 'barjas', percentage: row.score_barjas },
    { key: 'fisik', percentage: row.score_fisik },
    { key: 'anggaran', percentage: row.score_anggaran },
    { key: 'kinerja', percentage: row.score_kinerja }
  ]
})
```

#### 3. **Chart Integration**

```vue
<!-- CardRealisasiBulanSection.vue -->
<CardRealisasiBulanSection
  :realisasi-bulan="processedRealisasiBulan"
  :loading="dashboard.loading.value.bulan"
  :error="dashboard.error.bulan"
/>
```

### Performance Categories

The system tracks four main performance indicators:

- **Barjas (Barang dan Jasa)**: Procurement of goods and services
- **Fisik**: Physical progress indicators
- **Anggaran**: Budget realization percentages
- **Kinerja**: Overall performance metrics

### Real-time Updates

```javascript
// Auto-refresh configuration
const refreshInterval = ref(0) // Configurable interval
const startAutoRefresh = () => {
  if (refreshInterval.value > 0) {
    refreshTimer = setInterval(() => {
      refreshData()
    }, refreshInterval.value)
  }
}
```

### Error Recovery

```javascript
// Retry mechanisms
const retryFetch = async (type) => {
  switch (type) {
    case 'bulan':
      await fetchRealisasiBulan()
      break
    case 'tahun':
      await fetchRealisasiTahun()
      break
    // ... other retry cases
  }
}
```

## Development Guidelines

### Code Organization

1. **Composables**: Business logic in dedicated composable files
2. **Components**: UI components focused on presentation
3. **Utils**: Pure functions for data processing
4. **Types**: TypeScript interfaces for data structures

### Best Practices Implementation

1. **Reactive State**: Vue 3 Composition API for reactive data
2. **Error Boundaries**: Centralized error handling
3. **Loading States**: Granular loading indicators
4. **Type Safety**: Progressive TypeScript adoption
5. **Performance**: Lazy loading and code splitting

### Testing Strategy

```javascript
// Component testing example
import { mount } from '@vue/test-utils'
import { createTestingPinia } from '@pinia/testing'

test('dashboard loads data correctly', async () => {
  const wrapper = mount(Dashboard, {
    global: {
      plugins: [createTestingPinia()]
    }
  })

  // Test loading states, data rendering, error handling
})
```

---

For technical support or questions about the API integration, refer to the code comments in the composable files, the comprehensive documentation in `FRONTEND_DOCUMENTATION.md`, or contact the development team.
