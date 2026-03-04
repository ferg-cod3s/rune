---
layout: ../../layouts/DocsLayout.astro
title: "Configuration Guide - Rune CLI | YAML Setup & Customization"
description: "Complete configuration guide for Rune CLI. Learn how to set up YAML config files, configure projects, rituals, integrations, focus mode, and customize your developer workflow automation."
image: "/og-configuration.png"
category: "Configuration"
lastModified: "2025-03-04T00:00:00Z"
---

# Configuration Guide

Rune uses a YAML configuration file located at `~/.rune/config.yaml`. This file controls all aspects of Rune's behavior.

## Configuration Structure

```yaml
version: 1
settings:
  work_hours: 8.0
  break_interval: 50m
  idle_threshold: 10m
  notifications:
    enabled: true
    break_reminders: true
    end_of_day_reminders: true
    session_complete: true
    idle_detection: true
    sound: true

projects:
  - name: "main-app"
    detect: ["git:main-app", "dir:~/projects/main-app"]

rituals:
  start:
    global: []
    per_project: {}
  stop:
    global: []
    per_project: {}

integrations:
  git:
    enabled: true
    auto_detect_project: true
  slack:
    workspace: "myteam"
    dnd_on_start: true
  calendar:
    provider: "google"
    block_calendar: true
  telemetry:
    enabled: false
    sentry_dsn: ""

logging:
  level: "info"
  format: "json"
  output: "stdout"
  error_file: ""
```

## Settings Section

### Basic Settings

```yaml
settings:
  work_hours: 8.0        # Target work hours per day (0-24)
  break_interval: 50m    # Suggested break interval
  idle_threshold: 10m    # Idle time before auto-pause
```

### Notification Settings

```yaml
settings:
  notifications:
    enabled: true              # Enable OS-level notifications
    break_reminders: true      # Remind to take breaks
    end_of_day_reminders: true # Remind when approaching target hours
    session_complete: true     # Notify when session ends
    idle_detection: true       # Notify when idle detected
    sound: true                # Enable notification sounds
```

## Projects Section

Projects define how Rune detects and categorizes your work.

### Project Detection

```yaml
projects:
  - name: "web-app"
    detect:
      - "git:web-app"              # Git repository name
      - "dir:~/projects/web-app"   # Directory path

  - name: "api-service"
    detect:
      - "git:api-service"
      - "dir:~/work/api"
```

### Detection Patterns

- `git:name` - Match Git repository name
- `dir:path` - Match directory path (supports `~` expansion)

## Rituals Section

Rituals are automated command sequences that run during start/stop events.

### Basic Ritual Structure

```yaml
rituals:
  start:
    global:
      - name: "Update repositories"
        command: "git -C ~/projects pull --all"
        optional: false

    per_project:
      web-app:
        - name: "Start dev server"
          command: "bun run dev"
          background: true

  stop:
    global:
      - name: "Commit work in progress"
        command: "git add -A && git commit -m 'WIP: End of day'"
        optional: true
```

### Ritual Command Options

```yaml
- name: "Command description"    # Display name for the command
  command: "shell command"        # Shell command to execute
  optional: true                  # Don't fail ritual if command fails (default: false)
  background: false               # Run in background (default: false)
  interactive: false              # Enable interactive mode with TTY (default: false)
  tmux_session: "session-name"   # Custom tmux session name (optional)
  tmux_template: "template-name" # Reference to a tmux template (optional)
```

### Interactive Commands & tmux Templates

Commands with `interactive: true` can create tmux sessions for development environments:

```yaml
rituals:
  start:
    global:
      - name: "Setup Development Environment"
        command: "echo 'Starting dev environment...'"
        interactive: true
        tmux_template: "fullstack-dev"

  templates:
    fullstack-dev:
      session_name: "dev-{{.Project}}"
      windows:
        - name: "editor"
          layout: "main-horizontal"
          panes:
            - "vim ."
            - "git status"
        - name: "servers"
          layout: "tiled"
          panes:
            - "npm run dev"
            - "go run main.go"
```

### Ritual Types

Rune supports two ritual types:

```yaml
rituals:
  start:  # Run when starting work (via `rune start`)
  stop:   # Run when stopping work (via `rune stop`)
```

## Integrations Section

### Git Integration

```yaml
integrations:
  git:
    enabled: true               # Enable Git integration
    auto_detect_project: true   # Auto-detect project from Git repo
```

### Slack Integration

```yaml
integrations:
  slack:
    workspace: "myteam"       # Slack workspace name
    dnd_on_start: true         # Enable DND on Slack when starting work
```

### Calendar Integration

```yaml
integrations:
  calendar:
    provider: "google"         # Calendar provider (google)
    block_calendar: true       # Block time on calendar during work
```

### Telemetry Settings

```yaml
integrations:
  telemetry:
    enabled: false             # Enable telemetry (opt-in)
    sentry_dsn: ""             # Sentry DSN (prefer RUNE_SENTRY_DSN env var)
```

## Logging Section

```yaml
logging:
  level: "info"       # Log level: debug, info, warn, error
  format: "json"      # Log format: text, json
  output: "stdout"    # Output: stdout, stderr, or file path
  error_file: ""      # JSON file for structured error logging
```

## Environment Variables

Sensitive values should be set via environment variables:

```bash
# Telemetry
export RUNE_SENTRY_DSN="https://key@sentry.io/project_id"
export RUNE_OTLP_ENDPOINT="http://localhost:4318/v1/logs"

# General
export RUNE_TELEMETRY_DISABLED="true"  # Disable all telemetry
export RUNE_DEBUG="true"               # Enable debug logging
```

Environment variables override config file values.

## Validation

Validate your configuration:

```bash
rune config validate
```

Common validation errors:

- Invalid YAML syntax
- Unsupported config version (must be `1`)
- `work_hours` outside 0-24 range
- Negative `break_interval` or `idle_threshold`
- Empty project names or detection patterns
- References to undefined tmux templates

## Migration

### From Watson

```bash
rune migrate watson ~/.config/watson/frames
rune migrate watson ~/.config/watson/frames --dry-run
```

### From Timewarrior

```bash
rune migrate timewarrior ~/.timewarrior
rune migrate timewarrior ~/.timewarrior --dry-run
```

## Best Practices

1. **Version Control**: Keep your config in a dotfiles repository
2. **Environment Variables**: Use env vars for sensitive data (API keys, DSNs)
3. **Optional Commands**: Mark commands that may fail as `optional: true`
4. **Background Processes**: Use `background: true` for long-running dev servers
5. **Testing**: Test rituals with `rune ritual test start` before using `rune start`
