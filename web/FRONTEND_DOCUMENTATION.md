# Frontend Documentation

## Overview

This document provides comprehensive documentation for the frontend application built with Vue 3, Vuetify 3, and modern development tooling. The application serves as a dashboard for the SIJAGUR (Sistem Informasi Realisasi Anggaran) system, providing data visualization and management capabilities for government procurement performance tracking.

## Technology Stack

### Core Framework

- **Vue 3**: Progressive JavaScript framework with Composition API
- **Vite**: Fast build tool and development server
- **Vuetify 3**: Material Design component library
- **Vue Router**: Official routing library with auto-imports
- **Pinia**: State management solution

### Development Tools

- **ESLint**: Code linting with Vue-specific rules
- **Prettier**: Code formatting
- **TypeScript**: Type safety (partial adoption)
- **Auto-imports**: Automatic import management
- **Vue DevTools**: Development debugging tools

### UI & Visualization

- **ApexCharts**: Advanced charting library
- **Chart.js**: Simple charting solutions
- **Mapbox GL**: Interactive maps
- **TipTap**: Rich text editor
- **Swiper**: Touch-enabled sliders
- **Perfect Scrollbar**: Custom scrollbars

### HTTP & API

- **ofetch**: Lightweight HTTP client
- **JWT Decode**: Token parsing and validation
- **Cookie-es**: Cookie management

### Utilities

- **VueUse**: Collection of Vue composition utilities
- **Vue I18n**: Internationalization support
- **Day.js/Flatpickr**: Date/time handling
- **Prism.js**: Syntax highlighting
- **WebFont Loader**: Font loading optimization

## Project Structure

```
web/
├── src/
│   ├── @core/                    # Core framework components
│   │   ├── components/           # Reusable UI components
│   │   ├── composable/           # Core composables
│   │   ├── layouts/              # Layout components
│   │   ├── scss/                 # Core styling
│   │   ├── stores/               # Pinia stores
│   │   └── utils/                # Core utilities
│   ├── @layouts/                 # Layout system
│   ├── assets/                   # Static assets
│   ├── components/               # Application components
│   ├── composables/              # Business logic composables
│   ├── layouts/                  # Page layouts
│   ├── pages/                    # Route components
│   ├── plugins/                  # Vue plugins
│   ├── utils/                    # Utility functions
│   └── views/                    # Complex view components
├── public/                       # Public assets
├── tests/                        # Test files
├── vite.config.js               # Build configuration
├── themeConfig.js               # Theme configuration
└── package.json                 # Dependencies
```

## Architecture Patterns

### Composition API Usage

The application extensively uses Vue 3's Composition API with custom composables for:

```javascript
// Composable pattern
export function useDashboard() {
  const state = ref(initialValue)
  const computedValue = computed(() => state.value * 2)

  const action = async () => {
    // Business logic
  }

  return {
    state,
    computedValue,
    action
  }
}
```

### Component Organization

Components follow a hierarchical structure:

- **Pages**: Route-level components in `src/pages/`
- **Views**: Complex feature components in `src/views/`
- **Components**: Reusable UI components in `src/components/`
- **Core Components**: Framework-level components in `src/@core/components/`

### State Management

**Pinia Stores** are used for global state:

```javascript
// Store definition
export const useConfigStore = defineStore('config', {
  state: () => ({
    theme: 'light',
    language: 'en'
  }),

  actions: {
    setTheme(theme) {
      this.theme = theme
    }
  }
})
```

### API Integration

**Custom Composables** handle API communication:

```javascript
// useApi composable
export function useApi() {
  const $api = ofetch.create({
    baseURL: import.meta.env.VITE_API_BASE_URL
    // Interceptors for auth, error handling
  })

  return {
    get: (url) => $api(url),
    post: (url, data) => $api(url, { method: 'POST', body: data })
  }
}
```

## Key Features

### Authentication System

**useAuth Composable** manages authentication state:

- JWT token handling
- Automatic logout on token expiration
- Route protection middleware
- User role management

```javascript
const { isAuthenticated, login, logout, userRole } = useAuth()
```

### Dashboard Data Management

**useDashboard Composable** orchestrates dashboard data:

- Multiple API endpoints coordination
- Loading states management
- Error handling
- Auto-refresh functionality
- Data transformation and caching

### API Integration

**Seamless API Integration** with:

- Automatic token injection
- Centralized error handling
- Request/response interceptors
- Network status detection
- Authentication failure handling

## Component Architecture

### Layout System

The application uses a flexible layout system:

```vue
<template>
  <VApp>
    <VNavigationDrawer />
    <VAppBar />
    <VMain>
      <RouterView />
    </VMain>
  </VApp>
</template>
```

### Component Patterns

**Composition Components**:

```vue
<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  title: String,
  data: Array
})

const processedData = computed(() => {
  return props.data.map((item) => ({
    ...item,
    formatted: formatValue(item.value)
  }))
})
</script>

<template>
  <div class="component">
    <h3>{{ title }}</h3>
    <div v-for="item in processedData" :key="item.id">
      {{ item.formatted }}
    </div>
  </div>
</template>
```

### Styling Approach

**SCSS with Vuetify**:

- Custom theme variables
- Component-scoped styles
- Responsive design utilities
- Dark/light theme support

## Configuration

### Environment Variables

```env
# API Configuration
VITE_API_BASE_URL=http://localhost:8080/api
VITE_API_TIMEOUT=30000

# Application Settings
VITE_APP_TITLE=SIJAGUR Dashboard
VITE_APP_VERSION=1.0.0
```

### Vite Configuration

Advanced build configuration with:

- Vue Router auto-imports
- Component auto-registration
- Path aliases
- Plugin integrations

### Theme Configuration

Customizable theme system:

```javascript
export const themeConfig = {
  app: {
    title: 'SIJAGUR',
    logo: logoComponent,
    theme: 'system',
    skin: 'default'
  },
  navbar: {
    type: 'sticky',
    blur: true
  }
}
```

## Development Workflow

### Setup and Installation

```bash
# Install dependencies
pnpm install

# Start development server
pnpm run dev

# Build for production
pnpm run build

# Preview production build
pnpm run preview

# Format code
pnpm run format

# Build icons
pnpm run build:icons
```

### Code Quality

**ESLint Configuration**:

- Vue 3 specific rules
- TypeScript support
- Import sorting
- Code consistency enforcement

**Prettier Configuration**:

- Consistent code formatting
- Vue file support
- SCSS formatting

### Auto-imports

Automatic import management for:

- Vue composables
- Router functions
- UI components
- Utility functions

## Key Components

### Dashboard Components

**CardRealisasiBulanSection**: Monthly realization data display
**CardRealisasiTahunSection**: Yearly realization data display
**CardRealisasiPerbulanSection**: Monthly breakdown charts
**CardRankingSection**: Performance rankings display

### Core Components

**AppBarSearch**: Global search functionality
**Notifications**: Notification management
**ThemeSwitcher**: Theme switching
**ScrollToTop**: Navigation utility

### Form Components

**AppCombobox**: Enhanced select dropdowns
**AppDateTimePicker**: Date/time selection
**AppTextField**: Styled text inputs
**CustomCheckboxes**: Themed checkboxes

## State Management

### Pinia Stores

**Config Store**: Application configuration
**User Store**: User preferences and settings
**Notification Store**: Notification state

### Reactive State

**Composables** manage component-level state:

- `useDashboard`: Dashboard data and loading states
- `useAuth`: Authentication state
- `useApi`: API client state

## Routing

### Vue Router Configuration

Auto-generated routes with:

- File-based routing
- Route meta fields
- Authentication guards
- Lazy loading

### Route Protection

```javascript
// Route definition
definePage({
  meta: {
    requiresAuth: true,
    title: 'Dashboard',
    roles: ['admin', 'user']
  }
})
```

## API Integration

### Request/Response Handling

**Consistent API Patterns**:

```javascript
// API call with error handling
const { get } = useApi()

const fetchData = async () => {
  const response = await get('/api/data')

  if (response.success) {
    return response.data
  } else {
    throw new Error(response.error)
  }
}
```

### Authentication Flow

1. Login → Token storage
2. API requests → Automatic token injection
3. Token expiration → Automatic refresh/logout
4. Unauthorized → Redirect to login

## Styling & Theming

### SCSS Architecture

**Modular SCSS**:

- Base styles
- Component styles
- Layout styles
- Theme variables
- Utility classes

### Vuetify Theme

**Custom Theme Configuration**:

```scss
// Custom theme variables
$primary: #7367F0
$secondary: #A8AAAE
$success: #28C76F

// Apply theme
.v-application {
  background-color: var(--v-theme-background) !important;
}
```

### Responsive Design

**Breakpoint System**:

- Mobile-first approach
- Vuetify breakpoint utilities
- Custom responsive mixins

## Performance Optimization

### Build Optimizations

**Vite Features**:

- Fast HMR (Hot Module Replacement)
- Tree shaking
- Code splitting
- Asset optimization

### Runtime Optimizations

**Vue Optimizations**:

- Lazy loading of routes
- Component async loading
- Virtual scrolling for large lists
- Image lazy loading

## Testing

### Test Setup

**Vitest Configuration**:

- Vue 3 component testing
- Composition API testing
- Mock utilities
- Coverage reporting

### Test Patterns

```javascript
// Component test example
import { mount } from '@vue/test-utils'

test('renders component', () => {
  const wrapper = mount(Component, {
    props: { title: 'Test' }
  })

  expect(wrapper.text()).toContain('Test')
})
```

## Deployment

### Build Process

```bash
# Production build
pnpm run build

# Output directory: dist/
# - Optimized bundles
# - Minified assets
# - Service worker (if configured)
```

### Environment Configuration

**Multi-environment Support**:

- Development
- Staging
- Production
- Environment-specific variables

## Best Practices

### Code Organization

1. **Component Structure**: Clear separation of concerns
2. **Composable Logic**: Extract reusable logic
3. **Type Safety**: Use TypeScript where beneficial
4. **Consistent Naming**: Follow Vue.js conventions

### Performance

1. **Lazy Loading**: Route and component lazy loading
2. **Memoization**: Computed properties for expensive operations
3. **Debouncing**: Input handling optimization
4. **Bundle Splitting**: Strategic code splitting

### Accessibility

1. **Semantic HTML**: Proper element usage
2. **ARIA Labels**: Screen reader support
3. **Keyboard Navigation**: Full keyboard support
4. **Color Contrast**: WCAG compliance

## Areas for Improvement

### Code Quality

1. **TypeScript Adoption**: Expand TypeScript usage for better type safety
2. **Test Coverage**: Increase unit and integration test coverage
3. **Error Boundaries**: Implement Vue error boundaries
4. **Performance Monitoring**: Add performance tracking

### Architecture

1. **State Management**: Consider more centralized state for complex features
2. **Caching Strategy**: Implement intelligent API response caching
3. **Offline Support**: Add service worker for offline functionality
4. **Real-time Updates**: WebSocket integration for live data

### Developer Experience

1. **Storybook**: Component documentation and testing
2. **Automated Testing**: CI/CD pipeline with automated tests
3. **Code Generation**: CLI tools for component scaffolding
4. **Documentation**: API documentation integration

### Security

1. **CSP Headers**: Content Security Policy implementation
2. **Input Validation**: Client-side validation enhancement
3. **XSS Prevention**: Sanitize user inputs
4. **Secure Storage**: Review token storage security

## Conclusion

This Vue 3 application demonstrates modern frontend development practices with excellent tooling, architecture, and user experience. The combination of Vue 3, Vuetify 3, Vite, and comprehensive tooling provides a solid foundation for scalable dashboard applications.

The architecture emphasizes:

- **Composability**: Reusable logic through composables
- **Type Safety**: Progressive TypeScript adoption
- **Developer Experience**: Excellent tooling and DX
- **Performance**: Optimized build and runtime performance
- **Maintainability**: Clear code organization and patterns

The application successfully integrates complex data visualization, authentication, and API management while maintaining clean, maintainable code.
