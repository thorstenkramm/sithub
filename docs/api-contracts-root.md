# API Contracts (root)

## Scan Method
- Quick scan (pattern-based). No source files were read.
- Searched for route/controller/handler directories and API client patterns.

## Detected Evidence
- No route/controller/handler directories found.
- No API-related source files detected via patterns.

## Known API Notes (from README)
- Backend exposes a REST API that follows the JSON:API specification.
- Authentication uses Entra ID (SSO), with group-based access control.

## Documented Endpoints
- No concrete endpoints detected in this scan.

## Next Steps to Complete
- Provide the backend source tree (or enable deep scan) to extract routes.
- Identify base URL, versioning, and error format if applicable.
- Capture authentication flow and required headers (e.g., JSON:API media type).
