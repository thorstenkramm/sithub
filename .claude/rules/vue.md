# Frontend Tech Stack and Rules

- [Vue 3](https://vuejs.org/)
- [TypeScript](https://vuejs.org/guide/typescript/overview.html)
- [Create-Vue](https://github.com/vuejs/create-vue)
- Node.js v24 and npm
- [Pinia](https://pinia.vuejs.org/)
- [Vue Router](https://router.vuejs.org/guide/advanced/navigation-guards.html)
- [Vuetify](https://vuetifyjs.com/) for ready-to-use components
- ESLint and [eslint-plugin-vue](https://github.com/vuejs/eslint-plugin-vue/tree/master)
- [Vue Test Utils](https://test-utils.vuejs.org/) and [Vitest](https://vuejs.org/guide/scaling-up/testing.html) as the
  test runner
- [Cypress](https://www.cypress.io/) for browser testing

The development follows the [official Vue style guide](https://vuejs.org/style-guide/).

## Preface

This project is developed by experienced senior software developers with a strong focus on future maintenance.
Code is written for humans. Understandability and a comprehensive structure are crucial.

## UI Components
Use Vuetify for ready-to-use components. Use the installed [Vuetify MCP](https://github.com/vuetifyjs/mcp).
This integration enables seamless access to Vuetify's extensive component ecosystem directly within your development
workflow.

## Testing and Linting

After each significant change and after completing a task, run the tests and linters as follows.
Fix the discovered issues and run the tests or linters again until they pass without errors.

Chrome DevTools MCP is installed and configured. Use it for advanced browser debugging: Analyze network requests,
take screenshots, and check the browser console.

Run the dev server separately with `npm run dev` to execute tests.
Before starting the dev server, stop any other dev servers that might be running, for example: `pkill -f 'npm run dev'`.

### Security and vulnerability tests

```bash
# Base vulnerability scanning
npm audit
# Trivy scan, The All-in-One Security Scanner, https://trivy.dev/
trivy repository ./ --skip-version-check
```

Try to fix findings with `npm audit fix`. Never run `npm audit fix --force` as this can lead to unpredictable results.

### Linting

```bash
# ESLint
npm run lint
# Code duplication check
npx jscpd --pattern "**/*.ts" --ignore "**/node_modules/**" --threshold 0 --exitCode 1
```

### Unit Tests

```bash
npm run type-check
npx vitest run
```

Test coverage must always be 75% or higher. To get test coverage:

```bash
npm i -D @vitest/coverage-v8@3.2.4
npx vitest run --coverage
```

### Building

Test and make sure the project builds without errors.

```bash
npm run build
```

### Cypress E2E Tests

Running tests with the built-in electron browser is sufficient.

```bash
npm run test:e2e -- --browser electron
```

See `rules/cypress.md` for more details about E2E tests.

## Top 5 Dos

### 1. Use Composition API with `<script setup>`

Prefer the Composition API with `<script setup>` syntax for all new components. It provides better TypeScript
inference, less boilerplate, and improved code organization.

### 2. Use TypeScript with Proper Type Annotations

Define explicit types for props, emits, and reactive state. Use `defineProps<T>()` and `defineEmits<T>()` for
type-safe component interfaces.

### 3. Use Composables for Reusable Logic

Extract shared logic into composable functions (e.g., `useAuth`, `useForm`). Follow the `useX` naming convention
and return reactive refs or computed values.

### 4. Use Pinia Stores for Shared State

Organize application state in Pinia stores following the `useXStore` convention. Keep stores focused and avoid
storing derived data that could be computed.

### 5. Use Reactive Refs and Computed Properties Correctly

Use `ref()` for primitives, `reactive()` for objects, and `computed()` for derived values. Never mutate computed
properties directly.

## Top 5 Don'ts

### 1. Don't Use the Options API

Avoid `data()`, `methods`, `computed`, and `watch` options syntax. The Composition API is the standard for this
project and provides better TypeScript integration.

### 2. Don't Use `any` Type

Avoid TypeScript `any` type. Use `unknown` for truly unknown types and narrow them with type guards. Properly
type all variables, parameters, and return values.

### 3. Don't Mutate Props Directly

Never modify prop values in child components. Emit events to request changes from the parent or use
`v-model` with `defineModel()` for two-way binding.

### 4. Don't Use Mixins

Mixins cause naming conflicts and unclear data origins. Use composables instead for shared logic between
components.

### 5. Don't Access Reactive State Without `.value`

Remember that `ref()` values require `.value` in script code. Forgetting this leads to bugs where reactivity
is lost or values are not updated correctly.



## API Communication

This project uses a layered approach for API communication with the SitHub backend:

```text
┌─────────────────────────────────────────────────────────┐
│  Vue Components                                         │
│  └── useApi() composable (reactive loading/error state) │
│       └── api service (fetch wrapper, interceptors)     │
│            └── TypeScript types (JSON:API responses)    │
└─────────────────────────────────────────────────────────┘
```

### Architecture Overview

| Layer      | Location                    | Purpose                                          |
|------------|-----------------------------|--------------------------------------------------|
| Types      | `src/types/api.ts`          | TypeScript interfaces for API responses          |
| Service    | `src/lib/api.ts`            | Core fetch wrapper, interceptors, error handling |
| Composable | `src/composables/useApi.ts` | Reactive state wrapper for components            |

**When to use which:**

- **In components:** Use `useApi()` composable for reactive loading/error state
- **In Pinia stores:** Use the `api` service directly
- **In utility functions:** Use the `api` service directly

### TypeScript Types for API Responses

The backend uses JSON:API format. Define types in `src/types/api.ts`:

```typescript
// Generic JSON:API types
export interface JsonApiResource<T extends string, A> {
    type: T
    id: string
    attributes: A
}

export interface JsonApiResponse<T extends string, A> {
    data: JsonApiResource<T, A>
}

// Domain-specific types
export interface UserAttributes {
    name: string
    email: string
    createdAt: string
}

export type UserResponse = JsonApiResponse<'user', UserAttributes>
export type UsersResponse = JsonApiCollectionResponse<'user', UserAttributes>
```

### Using the API Service

The `api` service provides typed HTTP methods:

```typescript
import {api} from '@/lib/api'
import type {UserResponse, UsersResponse} from '@/types/api'

// GET request
const user = await api.get<UserResponse>('/users/123')
console.log(user.data.attributes.name)

// POST request
const newUser = await api.post<UserResponse>('/users', {
    data: {
        type: 'user',
        attributes: {name: 'Alice', email: 'alice@example.com'},
    },
})

// Other methods: put, patch, del
```

### Using the useApi Composable

For components that need reactive loading and error state:

```vue

<script setup lang="ts">
  import {useApi, api} from '@/composables/useApi'
  import type {UsersResponse} from '@/types/api'

  const {
    data: users,
    error,
    isLoading,
    execute: fetchUsers,
  } = useApi<UsersResponse>(() => api.get('/users'))

  // Fetch on mount
  fetchUsers()
</script>

<template>
  <div v-if="isLoading">Loading...</div>
  <div v-else-if="error" class="text-destructive">{{ error }}</div>
  <ul v-else-if="users">
    <li v-for="user in users.data" :key="user.id">
      {{ user.attributes.name }}
    </li>
  </ul>
</template>
```

For immediate execution on component mount:

```typescript
const {data, isLoading} = useApi<UsersResponse>(
    () => api.get('/users'),
    {immediate: true}
)
```

### Error Handling

The API service throws `ApiError` for failed requests:

```typescript
import {ApiError} from '@/types/api'

try {
    await api.get('/users/123')
} catch (err) {
    if (err instanceof ApiError) {
        console.log(err.message)  // Error message
        console.log(err.status)   // HTTP status code (404, 500, etc.)
        console.log(err.errors)   // JSON:API errors array (if present)
    }
}
```

The `useApi` composable handles errors automatically and exposes them via the `error` ref.

### Request/Response Interceptors

Register interceptors for cross-cutting concerns like authentication:

```typescript
import {api} from '@/lib/api'

// Add auth header to all requests
const removeInterceptor = api.onRequest((url, options) => {
    const token = localStorage.getItem('auth_token')
    if (token) {
        options.headers = {
            ...options.headers,
            Authorization: `Bearer ${token}`,
        }
    }
    return options
})

// Handle 401 responses globally
api.onResponse((response) => {
    if (response.status === 401) {
        // Redirect to login or refresh token
        window.location.href = '/login'
    }
    return response
})

// Remove interceptor when no longer needed
removeInterceptor()
```

### Using API in Pinia Stores

For shared state, use the api service directly in stores:

```typescript
// src/stores/useUsersStore.ts
import {defineStore} from 'pinia'
import {ref} from 'vue'
import {api} from '@/lib/api'
import {ApiError} from '@/types/api'
import type {UsersResponse, UserResponse} from '@/types/api'

export const useUsersStore = defineStore('users', () => {
    const users = ref<UsersResponse['data']>([])
    const isLoading = ref(false)
    const error = ref<string | null>(null)

    async function fetchUsers() {
        isLoading.value = true
        error.value = null

        try {
            const response = await api.get<UsersResponse>('/users')
            users.value = response.data
        } catch (err) {
            error.value = err instanceof ApiError ? err.message : 'Failed to fetch users'
        } finally {
            isLoading.value = false
        }
    }

    return {users, isLoading, error, fetchUsers}
})
```

### Best Practices

1. **Always define response types** in `src/types/api.ts` before implementing API calls. This catches
   errors at compile time.

2. **Use useApi in components** for simple data fetching with loading states. It reduces boilerplate
   and ensures consistent error handling.

3. **Use api service in stores** when multiple components share the same data or when you need to
   cache responses.

4. **Handle errors explicitly** - don't let API errors crash your application. The `useApi` composable
   handles this automatically; in stores, wrap calls in try/catch.

5. **Use interceptors sparingly** - they're powerful but can make debugging harder. Good use cases:
   authentication headers, logging, global error handling.

6. **Keep API types in sync** with the backend. When the backend changes, update `src/types/api.ts`
   first, then fix any TypeScript errors that arise.
