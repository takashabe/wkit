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