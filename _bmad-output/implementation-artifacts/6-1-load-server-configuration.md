# Story 6.1: Load Server Configuration

Status: complete

## Story

As an operator,
I want the server to load settings from a config file,
So that I can control listen address, port, and data directory.

## Acceptance Criteria

1. **Given** a valid configuration file  
   **When** the server starts  
   **Then** the server loads the settings  
   **And** invalid settings prevent startup with a clear error

## Tasks / Subtasks

- [x] Server loads TOML configuration file (AC: 1)
- [x] Listen address configurable via `main.listen` (AC: 1)
- [x] Port configurable via `main.port` (AC: 1)
- [x] Data directory configurable via `main.data_dir` (AC: 1)
- [x] Invalid/missing config prevents startup with clear error (AC: 1)
- [x] Environment variables can override config (SITHUB_ prefix) (AC: 1)

## Dev Notes

This functionality was implemented as foundational infrastructure work.

### Implementation Details

- Configuration loaded via `config.Load()` function using Viper
- TOML format for human readability
- Environment variable support with `SITHUB_` prefix
- Defaults: listen=127.0.0.1, port=9900, data_dir=.

### File List

**Existing files (implemented previously):**
- `internal/config/config.go` - Configuration loading and validation
- `internal/config/config_test.go` - Tests
- `cmd/sithub/main.go` - CLI that loads config
- `config.toml.example` - Example configuration

### Change Log

- Pre-sprint: Implemented as foundational work
- 2026-01-19: Story documented as complete (functionality already exists)
