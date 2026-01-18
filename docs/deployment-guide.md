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

## CI/CD
- README mentions GitHub Actions, but no workflow files were found in this repo.

## Missing Details
- Environment variables / secrets management
- Health checks / monitoring
- Backup/restore strategy for SQLite
- Deployment scripts or containerization details
