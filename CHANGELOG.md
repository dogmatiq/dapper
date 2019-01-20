# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->
[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [Unreleased]

### Added

- Added support for filters, allowing custom rendering on a per-value basis
- Added a built-in filter for `reflect.Type`

### Fixed

- Fixed a bug whereby IO errors were not propagated to the caller

## [0.2.0] - 2019-01-19

### Added

- Added `Print()`, which prints directly to `os.Stdout`

### Fixed

- Fixed support for formatting of structs with unexported fields

## [0.1.1] - 2019-01-16

### Changed

- Use [Iago] for indentation, etc

## [0.1.0] - 2019-01-16

- Initial release

<!-- references -->
[Unreleased]: https://github.com/dogmatiq/dapper
[0.1.1]: https://github.com/dogmatiq/dapper/releases/tag/v0.1.1
[0.1.0]: https://github.com/dogmatiq/dapper/releases/tag/v0.1.0
[0.2.0]: https://github.com/dogmatiq/dapper/releases/tag/v0.2.0

[Iago]: https://github.com/dogmatiq/iago

<!-- version template
## [0.0.1] - YYYY-MM-DD

### Added
### Changed
### Deprecated
### Removed
### Fixed
### Security
-->
