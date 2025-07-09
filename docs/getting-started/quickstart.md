# Quick Start Guide

Get up and running with Rune in 5 minutes!

## 1. Initialize Configuration

```bash
rune init --guided
```

This will walk you through:
- Setting your work hours
- Configuring break intervals
- Setting up your first project
- Creating basic start/stop rituals

## 2. Start Your First Session

```bash
rune start
```

This will:
- Execute your start rituals
- Begin time tracking
- Enable focus mode (if configured)

## 3. Check Your Status

```bash
rune status
```

You'll see:
```
██▀▀█ █   █ █▀▀█ ██▀▀
██▀▀  █   █ █  █ █▀▀▀
██ ██ █   █ █  █ █   
██  █ █████ █  █ ████

📊 Current Session
   Project: my-project
   Started: 9:00 AM (2h 30m ago)
   Status:  Working

🎯 Today's Progress
   Work Time: 2h 30m / 8h
   Breaks:    15m
   Projects:  my-project (2h 30m)

⚡ Active Rituals
   ✓ Git status check
   ✓ Docker containers started
   ✓ Development server running
```

## 4. Take a Break

```bash
rune pause
```

When ready to resume:
```bash
rune resume
```

## 5. End Your Day

```bash
rune stop
```

This will:
- Execute your stop rituals
- Generate a daily summary
- Clean up background processes

## 6. View Your Report

```bash
rune report --today
```

## Example Workflow

Here's a typical developer workflow:

```yaml
# ~/.rune/config.yaml
version: 1

settings:
  work_hours: 8.0
  break_interval: 25m
  idle_threshold: 5m

projects:
  - name: "web-app"
    detect: ["*/web-app/*", "git:web-app"]

rituals:
  start:
    global:
      - name: "Check Git Status"
        command: "git status --porcelain"
        optional: true
      - name: "Start Docker"
        command: "docker-compose up -d"
        background: true
        
  stop:
    global:
      - name: "Commit WIP"
        command: "git add -A && git commit -m 'WIP: End of day'"
        optional: true
      - name: "Stop Docker"
        command: "docker-compose down"
```

## Next Steps

- [Configuration Guide](../configuration/setup.md)
- [Command Reference](../commands/reference.md)
- [Integration Setup](../integrations/)
- [Example Workflows](../examples/)

## Need Help?

- Run `rune --help` for command help
- Check [Troubleshooting](./troubleshooting.md)
- Join our [Community Discussions](https://github.com/johnferguson/rune/discussions)