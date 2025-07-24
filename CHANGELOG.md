# Changelog

## [v0.2.4](https://github.com/takashabe/wkit/compare/v0.2.3...v0.2.4) - 2025-07-24

## [v0.2.1](https://github.com/takashabe/wkit/compare/v0.2.0...v0.2.1) - 2025-07-16
- feat: add GoReleaser configuration for automated releases by @takashabe in https://github.com/takashabe/wkit/pull/47
- feat: add directory copying support to copy_files configuration by @takashabe in https://github.com/takashabe/wkit/pull/48
- feat: add sticky-cwd support for worktree switching by @takashabe in https://github.com/takashabe/wkit/pull/49
- feat: add root command to display repository root directory by @takashabe in https://github.com/takashabe/wkit/pull/50
- feat: add tagpr configuration for automated release management by @takashabe in https://github.com/takashabe/wkit/pull/51

## [v0.2.0](https://github.com/takashabe/wkit/compare/v0.1.1...v0.2.0) - 2025-07-14
- feat: auto-create branch when adding worktree by @takashabe in https://github.com/takashabe/wkit/pull/10
- feat: change default worktree path to .wkit-worktrees by @takashabe in https://github.com/takashabe/wkit/pull/11
- bonsai readme by @takashabe in https://github.com/takashabe/wkit/pull/13
- feat: add auto-switch to new worktree after 'wkit add' by @takashabe in https://github.com/takashabe/wkit/pull/14
- refactor: improve Fish shell alias naming for better compatibility by @takashabe in https://github.com/takashabe/wkit/pull/15
- fix: improve Fisher installation compatibility for auto-switch feature by @takashabe in https://github.com/takashabe/wkit/pull/16
- perf: optimize wkit clean by batching git commands by @takashabe in https://github.com/takashabe/wkit/pull/17
- change: worktreeのデフォルトパスを.git/.wkit-worktreesに変更 by @takashabe in https://github.com/takashabe/wkit/pull/18
- feat: support XDG_CONFIG_HOME for global config file location by @takashabe in https://github.com/takashabe/wkit/pull/20
- feat: add remote branch checkout command for worktree creation by @takashabe in https://github.com/takashabe/wkit/pull/21
- doc: simplify README.md for better clarity and readability by @takashabe in https://github.com/takashabe/wkit/pull/22
- feat: enhance config show command with source file display and copy_files info by @takashabe in https://github.com/takashabe/wkit/pull/23
- feat: enhance copy_files to search for files by name across repository by @takashabe in https://github.com/takashabe/wkit/pull/24
- fix: create new branches from remote default branch instead of current HEAD by @takashabe in https://github.com/takashabe/wkit/pull/25
- fix: unify license to MIT by @takashabe in https://github.com/takashabe/wkit/pull/26
- Display relative paths from repository root in list and status commands by @takashabe in https://github.com/takashabe/wkit/pull/27
- fix: display "(root)" for empty relative paths in list and status commands by @takashabe in https://github.com/takashabe/wkit/pull/28
- fix: resolve worktree paths relative to repository root by @takashabe in https://github.com/takashabe/wkit/pull/29
- refactor: Clean up Go implementation and improve code quality by @takashabe in https://github.com/takashabe/wkit/pull/31
- chore: remove legacy Rust code after Go migration by @takashabe in https://github.com/takashabe/wkit/pull/32
- refactor: move Fish integration to examples pattern by @takashabe in https://github.com/takashabe/wkit/pull/33
- fix: ensure wkit list shows paths relative to main repository root by @takashabe in https://github.com/takashabe/wkit/pull/35
- fix: remove binary from git and add to .gitignore by @takashabe in https://github.com/takashabe/wkit/pull/37
- test: add comprehensive tests for copy files with nested directories by @takashabe in https://github.com/takashabe/wkit/pull/38
- refactor: remove redundant fish integration files from root by @takashabe in https://github.com/takashabe/wkit/pull/40
- Improve wkit list output format to match git worktree list style by @takashabe in https://github.com/takashabe/wkit/pull/41
- refactor: improve wkit list table formatting with text/tabwriter by @takashabe in https://github.com/takashabe/wkit/pull/42
- feat: add --base-branch flag to wkit add command by @takashabe in https://github.com/takashabe/wkit/pull/43
- Remove default_worktree_path backward compatibility by @takashabe in https://github.com/takashabe/wkit/pull/45
- feat: migrate configuration format from TOML to YAML by @takashabe in https://github.com/takashabe/wkit/pull/46

## [v0.1.1](https://github.com/takashabe/wkit/compare/v0.1.0...v0.1.1) - 2025-07-16

## [v0.1.0](https://github.com/takashabe/wkit/commits/v0.1.0) - 2025-06-12
- Implement comprehensive worktree management with Fish integration by @takashabe in https://github.com/takashabe/wkit/pull/5
- Implement z integration for frecency-based worktree jumping by @takashabe in https://github.com/takashabe/wkit/pull/6
- feat: implement development efficiency commands (status, clean, sync) by @takashabe in https://github.com/takashabe/wkit/pull/7
- feat: add automated release workflow and easy installation by @takashabe in https://github.com/takashabe/wkit/pull/8
