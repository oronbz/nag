# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.3.2] - 2026-03-04

### Fixed
- Due-date placeholder text cut off in create dialog

## [0.3.0] - 2026-03-04

### Added
- Sorting toggle (`s`) with default, created, due date, and title modes

### Fixed
- Due dates displaying in UTC instead of local timezone

## [0.2.0] - 2026-03-04

### Added
- Context-sensitive `n`/`e`/`d` keys for both reminders and lists
- List CRUD operations (create, rename, delete)

## [0.1.0] - 2026-03-04

### Added
- Initial release
- Two-panel layout with lists sidebar and reminders
- Smart lists: Today (includes overdue) and Scheduled
- Vim-style navigation (`j`/`k`, `g`/`G`, `Ctrl-d`/`Ctrl-u`)
- Create reminders with title, due date, time, and priority
- Complete/uncomplete toggle with 2-second grace period
- Delete with confirmation
- Open in Apple Reminders
- Show/hide completed reminders
- Auto-refresh (10-second polling)
- Fuzzy filter/search
- Mouse support
- Help overlay (`?`)
