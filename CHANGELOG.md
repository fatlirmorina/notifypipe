# Changelog

All notable changes to NotifyPipe will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.1] - 2025-11-05

### Fixed

- Fixed PocketBase API compatibility issues with `FindRecordsByExpr` (replaced with `FindRecordsByFilter`)
- Fixed 500 errors when accessing containers API
- Improved error handling - APIs now gracefully handle missing database collections
- Fixed SQLite compilation issues in Docker builds with musl libc
- Added proper CGO flags for Alpine Linux builds

### Changed

- Updated all API endpoints to use correct PocketBase DAO methods
- Improved database query error handling across all endpoints
- Docker build now uses Go 1.23

## [1.0.0] - 2025-11-05

### Added

- Initial release of NotifyPipe
- Docker event monitoring and detection
- Support for multiple notification channels (Slack, Telegram, Discord, Email)
- Web-based dashboard with dark theme
- Container-specific notification settings
- Event logging and history
- REST API for programmatic access
- Docker Compose deployment
- PocketBase integration for data persistence
- Real-time event monitoring
- Notification testing functionality
- Health check endpoints

### Features

- 🐳 Real-time Docker container monitoring
- 📬 Multi-channel notifications via Shoutrrr
- 🎨 Beautiful dark-themed UI
- 💾 Lightweight SQLite-based storage
- 🔧 Easy setup wizard
- 📊 Event history and analytics
- 🚀 One-command Docker Compose deployment

---

[1.0.1]: https://github.com/fatlirmorina/notifypipe/releases/tag/v1.0.1
[1.0.0]: https://github.com/fatlirmorina/notifypipe/releases/tag/v1.0.0
