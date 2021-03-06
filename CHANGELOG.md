# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog], and this project adheres to
[Semantic Versioning].

<!-- references -->
[Keep a Changelog]: https://keepachangelog.com/en/1.0.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html

## [0.4.1] - 2020-11-21

## Added

- Add `DefaultPrinter`, the printer used by `Write()`, `Format()` and `Print()`

## [0.4.0] - 2020-05-05

## Added

- Add `Config` to enscapsulate the configuration of a `Printer`

## Changed

- **[BC]** Change `DefaultIndent` from `string` constant to `[]byte`
- **[BC]** The `Filter` function signature now accepts a `Config` and `FilterPrinter`
- Zero-value structs are now collapsed to `StructName{<zero>}`

## Removed

- **[BC]** Remove `Printer.Filters`, `Indent` and `RecursionMarker`
- **[BC]** Remove `Value.TypeName()`

## Fixed

- Add mutex lock around writes to `stdout` to prevent fragmented output ([#45], thanks @ilmanzo)

## [0.3.5] - 2019-11-06

## Changed

- `Print()` now accepts multiple arguments ([#23])

## Fixed

- Fix panic when rendering unexported `time.Time` values ([#24])

## [0.3.4] - 2019-11-05

## Changed

- Byte slices and arrays are now rendered in hexdump format ([#15])
- Improved rendering of `sync.Mutex`, `RWMutex` and `Once` ([#14])
- Improved rendering of `time.Time` and `Duration` ([#8])

## [0.3.3] - 2019-04-23

### Fixed

- Render `reflect.Type` names when obtained from unexported struct fields ([#9])

## [0.3.2] - 2019-02-06

### Changed

- Bump Iago to v0.4.0

### Fixed

- `Printer()` now includes the name of user-defined `string` and `bool` types ([#6])
- `iago.Print()` now prints a trailing newline ([#7])

## [0.3.1] - 2019-01-29

### Changed

- Bump Iago to v0.3.0

## [0.3.0] - 2019-01-20

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
[0.1.0]: https://github.com/dogmatiq/dapper/releases/tag/v0.1.0
[0.1.1]: https://github.com/dogmatiq/dapper/releases/tag/v0.1.1
[0.2.0]: https://github.com/dogmatiq/dapper/releases/tag/v0.2.0
[0.3.0]: https://github.com/dogmatiq/dapper/releases/tag/v0.3.0
[0.3.1]: https://github.com/dogmatiq/dapper/releases/tag/v0.3.1
[0.3.2]: https://github.com/dogmatiq/dapper/releases/tag/v0.3.2
[0.3.3]: https://github.com/dogmatiq/dapper/releases/tag/v0.3.3
[0.3.4]: https://github.com/dogmatiq/dapper/releases/tag/v0.3.4
[0.3.5]: https://github.com/dogmatiq/dapper/releases/tag/v0.3.5
[0.4.0]: https://github.com/dogmatiq/dapper/releases/tag/v0.4.0
[0.4.1]: https://github.com/dogmatiq/dapper/releases/tag/v0.4.1

[#6]: https://github.com/dogmatiq/dapper/issues/6
[#7]: https://github.com/dogmatiq/dapper/issues/7
[#8]: https://github.com/dogmatiq/dapper/issues/8
[#9]: https://github.com/dogmatiq/dapper/issues/9
[#14]: https://github.com/dogmatiq/dapper/issues/14
[#15]: https://github.com/dogmatiq/dapper/issues/15
[#23]: https://github.com/dogmatiq/dapper/issues/23
[#24]: https://github.com/dogmatiq/dapper/issues/24
[#45]: https://github.com/dogmatiq/dapper/issues/45

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
