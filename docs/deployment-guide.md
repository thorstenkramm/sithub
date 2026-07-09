# Deployment Guide

## Scan Method

- Quick scan (pattern-based). No source files were read.
- No Docker/Kubernetes/IaC or CI/CD pipeline files detected.

## Deployment Model (from README)

- Single-binary distribution
- Frontend assets embedded into the binary
- Built-in SQLite3 database (no external DB dependency)

## Suggested Deployment Flow (from README)

1. Download the binary.
2. Create a configuration file (see `sithub.example.toml`).
3. Define rooms/desks in `sithub_areas.example.yaml`.
4. Start the server.
5. Configure a reverse proxy for SSL termination.

## Authentication & Access

- Entra ID SSO integration
- Group-based access and admin roles managed via Entra ID groups

### Session cookie keys (persistent sessions)

SitHub signs and encrypts its session cookie (`sithub_user`) and OAuth-state cookie with keys
stored in `cookie.key`, located in the configured `data_dir` (next to `sithub.db`). The file is
created automatically with random keys on first start (mode `0600`) and reused on every later
start, so **users stay logged in across server restarts and deployments**.

> [!IMPORTANT]
> Preserve `cookie.key` across upgrades and restarts. Deleting or replacing it generates a new
> key pair on the next start, which invalidates every active session — all users are redirected
> to login on their next request. A corrupt/truncated `cookie.key` is treated as a hard startup
> error rather than being silently regenerated, so an operator never logs out the whole user base
> by accident. Back it up together with `sithub.db`.

## CI/CD

- README mentions GitHub Actions, but no workflow files were found in this repo.

## Missing Details

- Environment variables / secrets management
- Health checks / monitoring
- Backup/restore strategy for SQLite
- Deployment scripts or containerization details
