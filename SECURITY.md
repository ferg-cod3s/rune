# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | ✅ Yes (current)   |
| < 0.1.0 | ❌ No              |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability in Rune, please report it responsibly.

### How to Report

**Email:** security@v1truv1us.dev (or your preferred security contact)

Please include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact assessment
- Suggested fix (if any)

### Response Timeline

- **Acknowledgment:** Within 48 hours
- **Initial assessment:** Within 7 days
- **Fix timeline:** Target 90 days (or as negotiated for complex issues)
- **Disclosure:** Coordinated disclosure 14 days after fix release

### Security Considerations for Rune

Rune executes user-defined commands from YAML configuration. Key security concerns:

1. **Command Execution:** Rune runs shell commands defined in rituals. Ensure:
   - Only trusted configuration sources
   - Sandboxing where possible
   - Audit logging enabled

2. **Environment Variables:** Rune may handle sensitive tokens. Ensure:
   - No logging of env vars with secrets
   - OS keychain integration for credentials

3. **Network Access:** Git integrations and updates require network. Ensure:
   - HTTPS for all external calls
   - Certificate validation

### Disclosure Policy

We follow coordinated disclosure:
1. Reporter submits vulnerability
2. We acknowledge and assess
3. We develop and test fix
4. We release fix and credit reporter (with permission)
5. Public disclosure after grace period

### Bug Bounty

Rune does not currently offer a bug bounty program. We appreciate responsible disclosure and will credit researchers in our CHANGELOG.
