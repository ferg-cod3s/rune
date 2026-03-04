---
layout: ../../layouts/DocsLayout.astro
title: "Configuration Examples - Rune CLI | Real-World Workflows"
description: "Real-world configuration examples for Rune CLI including frontend, backend, DevOps, and freelancer workflows. Copy and customize these proven productivity automation setups."
image: "/og-examples.png"
category: "Examples"
lastModified: "2025-03-04T00:00:00Z"
---

## Workflow Examples

<div class="feature-grid">
  <div class="feature-card">
    <span class="emoji">⚛️</span>
    <h4>Frontend Developer</h4>
    <p>Perfect for React, Vue, Svelte, or any frontend development workflow with hot reloading and dependency management.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">🔧</span>
    <h4>Backend Developer</h4>
    <p>Ideal for API development with Docker services, database migrations, and automated testing.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">☁️</span>
    <h4>DevOps Engineer</h4>
    <p>Perfect for infrastructure management, Kubernetes monitoring, and deployment automation.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">💼</span>
    <h4>Freelancer/Consultant</h4>
    <p>Time tracking focused with client project separation and automated reporting.</p>
  </div>
</div>

### Frontend Developer Configuration

Perfect for React, Vue, Svelte, or any frontend development workflow.

```yaml
version: 1
settings:
  work_hours: 8.0
  break_interval: 50m
  idle_threshold: 10m

projects:
  - name: "web-app"
    detect: ["git:web-app", "dir:~/projects/web-app"]

rituals:
  start:
    global:
      - name: "Pull latest changes"
        command: "git pull"
      - name: "Install dependencies"
        command: "pnpm install"
      - name: "Start dev server"
        command: "pnpm dev"
        background: true
      - name: "Open browser"
        command: "open http://localhost:3000"

  stop:
    global:
      - name: "Stop dev server"
        command: "pkill -f 'pnpm dev'"
        optional: true
      - name: "Commit WIP"
        command: "git add -A && git commit -m 'WIP: End of day'"
        optional: true

integrations:
  git:
    enabled: true
    auto_detect_project: true
```

### Backend Developer

Ideal for API development with Docker services.

```yaml
version: 1
settings:
  work_hours: 8.0
  break_interval: 45m
  idle_threshold: 10m

projects:
  - name: "api-service"
    detect: ["git:api-service"]

rituals:
  start:
    global:
      - name: "Start Docker services"
        command: "docker-compose up -d postgres redis"
      - name: "Run migrations"
        command: "go run migrate.go"
        optional: true
      - name: "Start API server"
        command: "air"
        background: true

  stop:
    global:
      - name: "Stop API server"
        command: "pkill -f air"
        optional: true
      - name: "Stop Docker services"
        command: "docker-compose down"
      - name: "Run tests"
        command: "go test ./..."
        optional: true
```

### DevOps Engineer

Perfect for infrastructure and deployment workflows.

```yaml
version: 1
settings:
  work_hours: 8.0
  break_interval: 60m
  idle_threshold: 5m

projects:
  - name: "infrastructure"
    detect: ["git:infra", "dir:~/infra"]

rituals:
  start:
    global:
      - name: "Check cluster status"
        command: "kubectl get nodes"
        optional: true
      - name: "Check monitoring"
        command: "open https://grafana.company.com"
      - name: "Review Terraform state"
        command: "terraform plan"
        optional: true

  stop:
    global:
      - name: "Security scan"
        command: "trivy fs ."
        optional: true
      - name: "Generate report"
        command: "./scripts/daily-report.sh"
        optional: true

integrations:
  git:
    enabled: true
    auto_detect_project: true
  slack:
    workspace: "devops-team"
    dnd_on_start: false
```

### Freelancer/Consultant

Time tracking focused with client project separation.

```yaml
version: 1
settings:
  work_hours: 6.0
  break_interval: 45m
  idle_threshold: 5m

projects:
  - name: "client-acme"
    detect: ["dir:~/clients/acme"]

  - name: "client-beta"
    detect: ["dir:~/clients/beta"]

rituals:
  start:
    global:
      - name: "Time tracking reminder"
        command: "echo 'Starting work session...'"

  stop:
    global:
      - name: "Export time report"
        command: "rune report --today --format csv"
        optional: true
      - name: "Backup work"
        command: "git add -A && git commit -m 'End of day backup'"
        optional: true
```

## Configuration Patterns

### Multi-Project Monorepo

```yaml
projects:
  - name: "frontend"
    detect: ["dir:~/monorepo/apps/web"]

  - name: "backend"
    detect: ["dir:~/monorepo/apps/api"]

  - name: "mobile"
    detect: ["dir:~/monorepo/apps/mobile"]

rituals:
  start:
    per_project:
      frontend:
        - name: "Start web app"
          command: "pnpm dev:web"
          background: true

      backend:
        - name: "Start API"
          command: "pnpm dev:api"
          background: true

      mobile:
        - name: "Start Metro"
          command: "pnpm dev:mobile"
          background: true
```

### Interactive Development Environment

Use tmux templates for multi-pane development setups:

```yaml
rituals:
  start:
    global:
      - name: "Setup dev environment"
        command: "echo 'Starting development environment...'"
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

## Integration Examples

<div class="feature-grid">
  <div class="feature-card">
    <span class="emoji">💬</span>
    <h4>Slack Integration</h4>
    <p>Automatically enable Do Not Disturb on Slack when starting work.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">📅</span>
    <h4>Calendar Integration</h4>
    <p>Block focus time on Google Calendar during work sessions.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">🔀</span>
    <h4>Git Automation</h4>
    <p>Automate commits, pulls, and status checking with Git rituals.</p>
  </div>
</div>

### Slack Integration Configuration

```yaml
integrations:
  slack:
    workspace: "myteam"
    dnd_on_start: true
```

### Calendar Integration

```yaml
integrations:
  calendar:
    provider: "google"
    block_calendar: true
```

### Git Automation via Rituals

```yaml
integrations:
  git:
    enabled: true
    auto_detect_project: true

rituals:
  start:
    global:
      - name: "Pull latest changes"
        command: "git pull"
        optional: true

  stop:
    global:
      - name: "Auto-commit progress"
        command: "git add -A && git commit -m 'Daily progress'"
        optional: true
```

## Platform-Specific Examples

<div class="feature-grid">
  <div class="feature-card">
    <span class="emoji">🍎</span>
    <h4>macOS Integration</h4>
    <p>Leverage Shortcuts and Do Not Disturb for seamless macOS workflow integration.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">🐧</span>
    <h4>Linux Integration</h4>
    <p>GNOME, KDE, and other desktop environment support with notification management.</p>
  </div>

  <div class="feature-card">
    <span class="emoji">🪟</span>
    <h4>Windows Integration</h4>
    <p>Focus Assist and PowerShell automation for Windows productivity workflows.</p>
  </div>
</div>

### macOS Configuration

Rune automatically manages Focus mode on macOS via the Shortcuts app. DND is enabled on `rune start` and disabled on `rune stop` when shortcuts are configured.

### Linux

On Linux, Rune uses `notify-send` for notifications and `dunstctl` or `makoctl` for DND control. Run `rune debug notifications` to check your setup.

### Windows

On Windows, Rune uses PowerShell for notifications and Focus Assist for DND management. Run `rune debug notifications` to verify your configuration.

## End of Day Cleanup

```yaml
rituals:
  stop:
    global:
      - name: "Commit work in progress"
        command: "git add -A && git commit -m 'WIP: End of session'"
        optional: true
      - name: "Push changes"
        command: "git push"
        optional: true
      - name: "Generate daily summary"
        command: "rune report --today --format json"
        optional: true
```

## More Examples

See the `examples/` directory in the repository for additional configurations:

- `config-basic.yaml` - Minimal setup for new users
- `config-developer.yaml` - Advanced developer with tmux templates
- `config-devops.yaml` - DevOps/SRE workflows
- `config-remote-work.yaml` - Remote work boundaries
- `config-freelancer.yaml` - Client billing and time tracking
- `config-academic.yaml` - Research and academic workflows
- `config-creative.yaml` - Creative professional workflows
- `config-content-creator.yaml` - Content creation workflows
- `config-team-lead.yaml` - Engineering management
- `config-telemetry.yaml` - Telemetry integration setup
