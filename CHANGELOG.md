# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

Nothing yet.

## [2.0.0] - 2025-11-07

Version 2.0 represents a complete refactoring focused on simplicity, maintainability, and modern Go practices.

### Added
- CLI tool with full command-line interface support
  - Pipe support for stdin input
  - Multiple flags: `-c` (cow), `-m` (mood), `-t` (think), `-l` (list), `-r` (random), `-w` (columns)
  - Version flag `-v`
- Modern web UI with polished design
  - Dark mode toggle
  - Dropdown selectors for cows and moods
  - Copy to clipboard functionality
  - Mobile responsive design
- JSON API endpoints
  - `POST /api/moo` - Main rendering endpoint (JSON and query params)
  - `GET /api/cows` - List all available cows
  - `GET /api/moods` - List all available moods
  - `GET /health` - Health check endpoint
- CORS middleware for cross-origin requests
- Comprehensive test suite with 69% overall coverage (97.6% for core rendering)
- Single binary deployment with embedded web assets
- 8 moods: borg, dead, greedy, paranoid, stoned, tired, wired, young
- 41 different cow variations

### Changed
- **Breaking**: Removed legacy Tokopedia dependencies (grace, logging, gcfg)
- **Breaking**: Switched to stdlib `net/http` for server
- **Breaking**: Configuration now via environment variables only
- Migrated to `log/slog` for structured logging
- Extracted core rendering logic to dedicated `cow/` package
- Refactored HTTP handlers to `api/` package
- Modernized codebase to Go 1.21+ standards
- Improved code readability and organization
- Optimized text wrapping and balloon building algorithms

### Removed
- Dependency on `gopkg.in/tokopedia/grace.v1`
- Dependency on `gopkg.in/tokopedia/logging.v1`
- Dependency on `gopkg.in/gcfg.v1`
- Config file support (replaced with environment variables)

### Fixed
- Improved Unicode character width calculation
- Better word wrapping for multi-line messages
- Various test improvements and bug fixes

### Dependencies
- `github.com/mattn/go-runewidth` v0.0.16 - Unicode width calculation
- `github.com/mitchellh/go-wordwrap` v1.0.1 - Text word wrapping

## [1.0.0] - Initial Release

### Added
- Basic HTTP server for Slack `/moo` command integration
- 41 cow templates
- 7 mood variations
- Text rendering with word wrap
- Token authentication
- Graceful shutdown support
