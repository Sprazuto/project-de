# API Integration Documentation

This document outlines the seamless integration between the `useAuth` and `useApi` composables across the SIJAGUR application.

## Overview

The integration provides a unified approach to API communication with automatic authentication token management, consistent error handling, and streamlined request/response patterns throughout the application.

## Architecture

### Core Composables

#### 1. `useAuth.js` - Authentication State Management

- **Location**: `src/composables/useAuth.js`
- **Purpose**: Manages user authentication state and provides token management methods
- **Key Features**:
  - Global authentication state
  - Automatic token management
  - User role and permission checking
  - Route protection middleware
  - Session persistence

#### 2. `useApi.js` - API Client with Auth Integration

- **Location**: `src/composables/useApi.js`
- **Purpose**: Provides authenticated API requests with seamless token handling
- **Key Features**:
  - Automatic token injection from `useAuth`
  - Centralized error handling
  - Request/response interceptors
  - Authentication failure handling
  - Network status detection

#### 3. `useErrorHandler.js` - Centralized Error Management

- **Location**: `src/composables/useErrorHandler.js`
- **Purpose**: Provides consistent error handling and user feedback
- **Key Features**:
  - Error type classification
  - User-friendly error messages
  - Error reporting and logging
  - Automatic auth error handling

## Integration Flow

### Authentication Flow

1. User logs in through `login.vue`
2. `useAuth` updates authentication state
3. `useApi` automatically retrieves tokens from `useAuth`
4. All subsequent API requests include proper authentication headers
5. Token expiration or invalid responses trigger automatic logout

### API Request Flow

1. Component calls API method (e.g., `api.get('/endpoint')`)
2. `useApi` automatically adds authentication headers from `useAuth`
3. Request is sent with proper error handling
4. Response is processed with consistent error handling
5. Authentication failures redirect to login automatically

## Usage Patterns

### Basic API Integration

```javascript
import { useAuth } from '@/composables/useAuth'
import { useApi } from '@/composables/useApi'

export default {
  setup() {
    const { isAuthenticated, getAccessToken } = useAuth()
    const { get, post, handleError } = useApi()

    const fetchData = async () => {
      try {
        // Authentication is handled automatically
        const response = await get('/api/data')

        if (response.success) {
          return response.data
        } else {
          throw new Error(response.error)
        }
      } catch (error) {
        // Use centralized error handling
        const userMessage = handleError(error)
        console.error('Error:', userMessage)
      }
    }

    return {
      fetchData,
      isAuthenticated,
    }
  },
}
```

### Integration with useAuth

```javascript
import { useAuth } from '@/composables/useAuth'
import { useApi } from '@/composables/useApi'

export default {
  setup() {
    const { isAuthenticated, currentUser, login, logout, getAccessToken, hasValidToken } = useAuth()

    const { get, post } = useApi()

    // All API calls automatically use authentication
    const authenticatedRequest = async () => {
      const response = await get('/protected-endpoint')
      return response
    }

    return {
      isAuthenticated,
      currentUser,
      authenticatedRequest,
      login,
      logout,
    }
  },
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
      handleApiError,
    }
  },
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
  clearAuthState, // Clear all auth data
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

### API Endpoints Expected

The integration assumes the following API structure:

- `POST /v1/user/login` - User authentication
- `GET /v1/user/profile` - User profile (protected)
- `GET /v1/articles` - Articles (protected)
- `GET /sijagur/realisasi-bulan` - Monthly data (protected)
- `GET /sijagur/realisasi-tahun` - Yearly data (protected)

## Troubleshooting

### Common Issues

#### 1. Authentication Token Not Found

- Ensure `login()` is called before API requests
- Check if user is properly authenticated
- Verify token storage in localStorage

#### 2. API Requests Failing

- Check network connectivity
- Verify API endpoint URLs
- Ensure proper error handling

#### 3. Automatic Logout Not Working

- Check if 401 responses are properly handled
- Verify `useAuth.logout()` method availability
- Check router configuration

## Future Enhancements

### Planned Features

- [ ] Token refresh mechanism
- [ ] Request retry logic with exponential backoff
- [ ] API response caching
- [ ] WebSocket integration for real-time updates
- [ ] Offline support with service workers

### Extension Points

- Custom error handlers
- Additional authentication providers
- API request/response transformers
- Custom logging implementations

---

For technical support or questions about the API integration, refer to the code comments in the composable files or contact the development team.
