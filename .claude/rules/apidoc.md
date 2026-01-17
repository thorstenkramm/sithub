## API Documentation

- The API must be documented following OpenAPI Specification (OAS) version 3.1.
- API documentation is stored in `./api-doc` and must be split into multiple files, with `openapi.yaml` as the top-level
  entry point. Each top-level API endpoint (`/api/v1/ping`, `/api/v1/foo`, `/api/v1/bla`, etc.) should have its own file
  referenced by the main `openapi.yaml`.
- After editing the API documentation, lint it with `npx @redocly/cli lint --lint-config off ./api-doc/openapi.yaml`.
- Fix errors reported by the linter.
