---
layout: ../../layouts/DocsLayout.astro
title: "CLI Commands Reference - Rune CLI | Complete Command Guide"
description: "Complete reference guide for all Rune CLI commands including start, stop, pause, resume, report, config, ritual, and utility commands with examples and flags."
image: "/og-commands.png"
category: "Commands Reference"
lastModified: "2025-03-04T00:00:00Z"
---

# Commands Reference

## Core Commands

### `rune start`

Start your workday and execute start rituals.

```bash
rune start [project]
```

**Arguments:**

- `project` - Optional project name. If omitted, project is auto-detected from the current directory.

**Behavior:**

- Begins time tracking for your work session
- Executes global start rituals
- Executes project-specific start rituals (if detected)
- Launches interactive development environments (tmux sessions)
- Enables focus mode (Do Not Disturb) if configured

**Examples:**

```bash
rune start
rune start my-app
```

### `rune stop`

End your workday and execute stop rituals.

```bash
rune stop
```

**Behavior:**

- Stops time tracking for your work session
- Executes global stop rituals
- Executes project-specific stop rituals
- Disables focus mode (Do Not Disturb)
- Generates a session summary

### `rune pause`

Pause the current work session.

```bash
rune pause
```

**Behavior:**

- Pauses the active work timer
- Saves the current session state

### `rune resume`

Resume a paused work session.

```bash
rune resume
```

### `rune status`

Show current session status and statistics.

```bash
rune status
```

**Displays:**

- Current timer state (running, paused, stopped)
- Active project (if detected)
- Session duration
- Today's total work time
- Idle status
- Focus mode status

## Time Tracking Commands

### `rune report`

Generate time tracking reports.

```bash
rune report [flags]
```

**Flags:**

- `--today` - Show today's report (default if no period specified)
- `--week` - Show this week's report
- `--month` - Show this month's report
- `--project <name>` - Filter by project name
- `--format <format>` - Output format: `text` (default), `csv`, `json`
- `--output <file>` - Output file (default: stdout)

**Examples:**

```bash
rune report --today
rune report --week
rune report --month --project my-app
rune report --today --format csv --output report.csv
rune report --week --format json
```

## Configuration Commands

### `rune config`

Manage configuration settings.

#### `rune config edit`

Open configuration file in your default editor (`$EDITOR`, falls back to `vi`).

```bash
rune config edit
```

#### `rune config validate`

Validate configuration file syntax and content.

```bash
rune config validate
```

#### `rune config show`

Display the current configuration file contents.

```bash
rune config show
```

#### `rune config setup-telemetry`

Configure telemetry integration (Sentry).

```bash
rune config setup-telemetry [flags]
```

**Flags:**

- `--sentry-dsn <dsn>` - Set Sentry DSN for error tracking
- `--enable` - Enable telemetry after setting keys (default: true)
- `--disable` - Disable telemetry
- `--interactive` - Interactive setup mode
- `--show-examples` - Show example values and setup instructions

**Examples:**

```bash
rune config setup-telemetry --interactive
rune config setup-telemetry --sentry-dsn https://key@sentry.io/id
rune config setup-telemetry --disable
rune config setup-telemetry --show-examples
```

**Quick shortcut:**

```bash
rune config --add-sentry-dsn https://key@sentry.io/id
```

## Data Migration Commands

### `rune migrate`

Migrate data from other time tracking tools.

#### `rune migrate watson`

Import from Watson time tracker.

```bash
rune migrate watson <frames-file> [flags]
```

**Arguments:**

- `frames-file` - Path to Watson frames file (typically `~/.config/watson/frames`)

**Flags:**

- `--dry-run` - Preview import without actually importing
- `--project-map <old=new>` - Map old project names to new ones

**Examples:**

```bash
rune migrate watson ~/.config/watson/frames
rune migrate watson ~/.config/watson/frames --dry-run
rune migrate watson frames.json --project-map "old-project=new-project"
```

#### `rune migrate timewarrior`

Import from Timewarrior.

```bash
rune migrate timewarrior <data-dir> [flags]
```

**Arguments:**

- `data-dir` - Path to Timewarrior data directory (typically `~/.timewarrior`)

**Flags:**

- `--dry-run` - Preview import without actually importing
- `--project-map <old=new>` - Map old project names to new ones

**Examples:**

```bash
rune migrate timewarrior ~/.timewarrior
rune migrate timewarrior ~/.timewarrior --dry-run
```

## Ritual Commands

### `rune ritual`

Manage and test rituals.

#### `rune ritual list`

List all configured rituals.

```bash
rune ritual list
```

#### `rune ritual run`

Run a specific ritual without affecting time tracking.

```bash
rune ritual run <start|stop> [project]
```

**Arguments:**

- `start|stop` - The ritual type to run
- `project` - Optional project name (auto-detected if omitted)

**Examples:**

```bash
rune ritual run start
rune ritual run stop
rune ritual run start my-app
```

#### `rune ritual test`

Test a ritual configuration without actually executing commands.

```bash
rune ritual test <start|stop> [project]
```

**Arguments:**

- `start|stop` - The ritual type to test
- `project` - Optional project name (auto-detected if omitted)

**Examples:**

```bash
rune ritual test start
rune ritual test stop my-app
```

## Utility Commands

### `rune init`

Initialize Rune configuration.

```bash
rune init [flags]
```

**Flags:**

- `--guided` - Interactive setup wizard with telemetry opt-in

**Examples:**

```bash
rune init
rune init --guided
```

### `rune update`

Update Rune to the latest version from GitHub releases.

```bash
rune update [flags]
```

**Flags:**

- `--check` - Check for updates without installing

**Examples:**

```bash
rune update
rune update --check
```

### `rune logs`

Display recent application logs.

```bash
rune logs [flags]
```

**Flags:**

- `-n, --lines <count>` - Number of log lines to display (default: 50)
- `-l, --level <level>` - Filter by log level (debug, info, warn, error)
- `-c, --component <name>` - Filter by component
- `-f, --file <path>` - Specific log file to read
- `--json` - Output logs in JSON format
- `--follow` - Follow log output (like `tail -f`)
- `--list` - List available log files

**Examples:**

```bash
rune logs
rune logs -n 20 -l error
rune logs --json
rune logs --list
rune logs -f ~/.rune/logs/rune-session-2025.log
```

### `rune completion`

Generate shell completion scripts.

```bash
rune completion <bash|zsh|fish|powershell>
```

**Supported shells:**

- `bash`
- `zsh`
- `fish`
- `powershell`

**Examples:**

```bash
# Bash
rune completion bash > /etc/bash_completion.d/rune

# Zsh
rune completion zsh > "${fpath[1]}/_rune"

# Fish
rune completion fish > ~/.config/fish/completions/rune.fish
```

## Debugging & Testing Commands

### `rune debug`

Debug and diagnostic commands for troubleshooting.

#### `rune debug telemetry`

Debug telemetry configuration and connectivity. Shows system info, environment variables, configuration status, and sends a test event.

```bash
rune debug telemetry
```

#### `rune debug keys`

Show telemetry key configuration (masked for security).

```bash
rune debug keys
```

#### `rune debug notifications`

Diagnose OS notification tooling and provide setup guidance.

```bash
rune debug notifications
```

### `rune test`

Test various Rune functionality.

#### `rune test notifications`

Send test notifications to verify the OS notification system is working.

```bash
rune test notifications
```

#### `rune test dnd`

Test Do Not Disturb functionality (enable, disable, check shortcuts).

```bash
rune test dnd
```

#### `rune test logging`

Test the structured logging system by generating test log entries at various levels.

```bash
rune test logging
```

## Global Flags

These flags are available for all commands:

- `--config <file>` - Specify config file (default: `~/.rune/config.yaml`)
- `-v, --verbose` - Enable verbose output
- `--no-color` - Disable colored output
- `--log` - Print recent logs
- `-h, --help` - Show help for command
- `--version` - Show version information
