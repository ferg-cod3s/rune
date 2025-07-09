# Rune Configuration Examples

This directory contains example configurations for different use cases and workflows. Choose the one that best matches your work style and customize it to fit your needs.

## Available Examples

### 1. Basic Configuration (`config-basic.yaml`)
**Best for:** New users getting started with Rune

**Features:**
- Minimal setup with essential features
- Simple Git integration
- Basic start/stop rituals
- 5-minute idle detection
- Pomodoro-style 25-minute breaks

**Use this if:** You want to start simple and gradually add more features.

### 2. Developer Configuration (`config-developer.yaml`)
**Best for:** Software developers working on multiple projects

**Features:**
- Project-specific rituals for frontend, backend, mobile, and documentation
- Automated dependency installation and testing
- Background development servers
- Code formatting and security checks
- Slack integration with DND mode

**Use this if:** You're a developer juggling multiple codebases and want automated project setup.

### 3. Remote Work Configuration (`config-remote-work.yaml`)
**Best for:** Remote workers who need strong communication boundaries

**Features:**
- Slack status automation
- Calendar integration
- Communication app management
- Meeting preparation rituals
- End-of-day summaries
- Longer breaks optimized for remote work

**Use this if:** You work remotely and need help managing communication and boundaries.

### 4. Freelancer Configuration (`config-freelancer.yaml`)
**Best for:** Freelancers managing multiple client projects

**Features:**
- Strict time tracking for billing
- Client-specific project detection
- Automated billing logs
- Time report generation
- Business metrics tracking
- Flexible work hours

**Use this if:** You're a freelancer who needs precise time tracking for client billing.

## How to Use These Examples

1. **Choose an example** that matches your workflow
2. **Copy the file** to your Rune config directory:
   ```bash
   cp examples/config-developer.yaml ~/.rune/config.yaml
   ```
3. **Customize the configuration** to match your specific needs:
   - Update project detection patterns
   - Modify ritual commands
   - Adjust timing settings
   - Configure integrations

4. **Test your configuration**:
   ```bash
   rune init --guided  # Validate your config
   rune ritual test start  # Test start rituals
   ```

## Configuration Structure

All configurations follow this structure:

```yaml
version: 1

settings:
  work_hours: 8.0          # Daily work hour target
  break_interval: 30m      # Time between break reminders
  idle_threshold: 5m       # Auto-pause threshold

projects:
  - name: "project-name"   # Project identifier
    detect: ["pattern"]    # File/directory patterns for detection

rituals:
  start:
    global: []             # Commands run for all projects
    per_project: {}        # Project-specific commands
  stop:
    global: []
    per_project: {}

integrations:
  git:
    enabled: true
    auto_detect_project: true
  slack:
    workspace: "team-name"
    dnd_on_start: true
  calendar:
    provider: "google"
    block_calendar: true
```

## Customization Tips

### Project Detection
Projects are automatically detected based on file patterns:
- `package.json` → Node.js project
- `go.mod` → Go project
- `Cargo.toml` → Rust project
- `requirements.txt` → Python project
- `.git` → Git repository

### Ritual Commands
- Use `optional: true` for commands that might fail
- Use `background: true` for long-running processes
- Environment variables are available in commands
- Commands run in the project directory

### Time Settings
- `work_hours`: Daily target (used for reporting)
- `break_interval`: How often to remind about breaks
- `idle_threshold`: Auto-pause after inactivity

### Integration Setup
- **Git**: Automatic project detection from repositories
- **Slack**: Requires `SLACK_TOKEN` environment variable
- **Calendar**: Requires calendar CLI tools (gcalcli, etc.)

## Environment Variables

Some examples use environment variables:
- `SLACK_TOKEN`: For Slack API integration
- `GITHUB_TOKEN`: For GitHub API calls
- `GOOGLE_CALENDAR_TOKEN`: For calendar integration

Set these in your shell profile:
```bash
export SLACK_TOKEN="xoxp-your-token-here"
export GITHUB_TOKEN="ghp_your-token-here"
```

## Troubleshooting

### Common Issues

1. **Commands not found**: Ensure all tools used in rituals are installed
2. **Permission denied**: Check file permissions for scripts
3. **Slow startup**: Remove or make optional any slow commands
4. **Git errors**: Ensure you're in a Git repository

### Testing Your Configuration

```bash
# Validate configuration syntax
rune config validate

# Test specific rituals
rune ritual test start --project myproject
rune ritual test stop --project myproject

# Check project detection
rune status  # Shows detected project
```

### Getting Help

- Check the main documentation: `rune --help`
- Validate your config: `rune config validate`
- Test rituals: `rune ritual test [start|stop]`
- Report issues: [GitHub Issues](https://github.com/johnferguson/rune/issues)

## Contributing Examples

Have a great configuration for a specific workflow? We'd love to include it! Please:

1. Create a new example file following the naming pattern
2. Add documentation to this README
3. Test thoroughly with different projects
4. Submit a pull request

Examples we'd love to see:
- Academic/research workflows
- Design/creative workflows
- DevOps/infrastructure workflows
- Content creation workflows
- Team lead/management workflows