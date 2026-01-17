## API Proxy (Development)

During development, the backend API runs on its own web server at
<http://localhost:9900>. The Vite dev server proxies all `/api`
requests to that backend, so the frontend can call `/api/v1/*`
without changing hosts. If your backend uses a different port,
update `vite.config.ts`.

Run both services in separate terminals:

1. Backend:

   ```bash
   cd backend
   go run cmd/sithub/main.go run \
     --config ./sithub.toml
   ```

2. Frontend:

   ```bash
   npm run dev
   ```

Once both are running, verify the proxy with:

```bash
curl http://localhost:5173/api/v1/ping
```

## Available Scripts

- `npm run dev`: Start the Vite development server.
- `npm run build`: Build for production.
- `npm run preview`: Preview the production build.
- `npm run test:unit`: Run unit tests with Vitest.
- `npm run test:e2e`: Run Cypress tests headless.
- `npm run test:e2e:open`: Open the Cypress UI runner.
- `npm run lint`: Lint and fix files with ESLint.
- `npm run type-check`: Run TypeScript type checking.