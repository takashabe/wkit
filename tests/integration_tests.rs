use std::process::Command;
use assert_cmd::prelude::*;
use predicates::prelude::*;

#[test]
fn test_cli_help() {
    let mut cmd = Command::cargo_bin("wkit").unwrap();
    cmd.arg("--help");
    cmd.assert()
        .success()
        .stdout(predicates::str::contains("Convenient git worktree management toolkit"));
}

#[test]
fn test_config_show() {
    let mut cmd = Command::cargo_bin("wkit").unwrap();
    cmd.args(["config", "show"]);
    cmd.assert()
        .success()
        .stdout(predicates::str::contains("Current configuration:"))
        .stdout(predicates::str::contains("default_worktree_path:"))
        .stdout(predicates::str::contains("auto_cleanup:"))
        .stdout(predicates::str::contains("z_integration:"));
}

#[test]
fn test_config_invalid_key() {
    let mut cmd = Command::cargo_bin("wkit").unwrap();
    cmd.args(["config", "set", "invalid_key", "value"]);
    cmd.assert()
        .failure()
        .stderr(predicates::str::contains("Unknown configuration key: invalid_key"));
}

#[test]
fn test_config_invalid_boolean() {
    let mut cmd = Command::cargo_bin("wkit").unwrap();
    cmd.args(["config", "set", "auto_cleanup", "not_a_boolean"]);
    cmd.assert()
        .failure()
        .stderr(predicates::str::contains("Invalid boolean value"));
}

#[test]
fn test_switch_nonexistent_worktree() {
    let mut cmd = Command::cargo_bin("wkit").unwrap();
    cmd.args(["switch", "nonexistent"]);
    cmd.assert()
        .failure()
        .stderr(predicates::str::contains("Worktree 'nonexistent' not found"));
}

#[test]
fn test_z_help() {
    let mut cmd = Command::cargo_bin("wkit").unwrap();
    cmd.args(["z", "--help"]);
    cmd.assert()
        .success()
        .stdout(predicates::str::contains("Z-style frecency-based worktree jumping"))
        .stdout(predicates::str::contains("--list"))
        .stdout(predicates::str::contains("--clean"))
        .stdout(predicates::str::contains("--add"));
}

#[test]
fn test_z_add_current_directory() {
    let mut cmd = Command::cargo_bin("wkit").unwrap();
    cmd.args(["z", "--add"]);
    cmd.assert()
        .success()
        .stdout(predicates::str::contains("Added current directory to z database"));
}

#[test]
fn test_z_list_without_query() {
    let mut cmd = Command::cargo_bin("wkit").unwrap();
    cmd.args(["z"]);
    cmd.assert()
        .success()
        .stdout(predicates::str::contains("Z database entries"));
}

#[test]
fn test_add_help_contains_no_switch_flag() {
    let mut cmd = Command::cargo_bin("wkit").unwrap();
    cmd.args(["add", "--help"]);
    cmd.assert()
        .success()
        .stdout(predicates::str::contains("--no-switch"))
        .stdout(predicates::str::contains("Skip automatic switching to new worktree"));
}

#[test]
fn test_checkout_help() {
    let mut cmd = Command::cargo_bin("wkit").unwrap();
    cmd.args(["checkout", "--help"]);
    cmd.assert()
        .success()
        .stdout(predicates::str::contains("Checkout a remote branch and create worktree"))
        .stdout(predicates::str::contains("Remote branch name (e.g., origin/feature-branch)"))
        .stdout(predicates::str::contains("--no-switch"));
}

#[test]
fn test_checkout_invalid_remote_branch_format() {
    let mut cmd = Command::cargo_bin("wkit").unwrap();
    cmd.args(["checkout", "invalid-format"]);
    cmd.assert()
        .failure()
        .stderr(predicates::str::contains("Invalid remote branch format"));
}