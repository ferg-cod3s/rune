# Contributing to Rune CLI

```
╔═╗ 
╠╦╝ une
╩╚═ 
```

Thank you for your interest in contributing to Rune! This document provides guidelines and information for contributors.

## Code of Conduct

This project adheres to a Code of Conduct that we expect all contributors to follow:

- **Be respectful**: Treat everyone with respect and kindness
- **Be inclusive**: Welcome newcomers and help them get started
- **Be collaborative**: Work together to build something amazing
- **Be constructive**: Provide helpful feedback and suggestions
- **Be patient**: Remember that everyone is learning

## Getting Started

### Prerequisites

- **Go 1.21+**: Primary development language
- **Git**: Version control
- **Make**: Build automation
- **Docker**: For testing cross-platform builds (optional)

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/YOUR_USERNAME/rune.git
   cd rune
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   ```

3. **Build and Test**
   ```bash
   make build
   make test
   ```

4. **Run Development Version**
   ```bash
   make dev
   ```

## Development Guidelines

### Code Style

We follow the guidelines in [AGENTS.md](AGENTS.md). Key points:

- **Test-Driven Development**: Write tests before implementing features
- **Security-First**: Validate all inputs, use secure defaults
- **Code Simplicity**: Readable code with clear, descriptive names
- **Go Conventions**: Follow effective Go guidelines and use gofmt

### Testing Requirements

- **Unit Tests**: All new code must have unit tests
- **Integration Tests**: Test command interactions and workflows
- **Coverage**: Maintain >80% test coverage for critical paths
- **Security Tests**: Test input validation and error handling

### Commit Guidelines

We use conventional commits for clear history:

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```bash
feat(rituals): add conditional execution based on day of week
fix(config): resolve YAML parsing error with nested arrays
docs(readme): update installation instructions for Windows
test(tracking): add integration tests for time tracking
```

## Types of Contributions

### 🐛 Bug Reports

When reporting bugs, please include:

- **Clear description** of the issue
- **Steps to reproduce** the problem
- **Expected vs actual behavior**
- **Environment details** (OS, Go version, etc.)
- **Configuration file** (sanitized of sensitive data)
- **Log output** if available

Use our bug report template:

```markdown
**Bug Description**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Run command '...'
2. See error

**Expected Behavior**
What you expected to happen.

**Environment**
- OS: [e.g. macOS 14.0]
- Go Version: [e.g. 1.21.0]
- Rune Version: [e.g. 0.1.0]

**Additional Context**
Any other context about the problem.
```

### ✨ Feature Requests

For new features, please:

- **Check existing issues** to avoid duplicates
- **Describe the use case** and problem being solved
- **Propose a solution** with examples
- **Consider backwards compatibility**
- **Think about security implications**

### 🔧 Code Contributions

#### Small Changes
- Typo fixes
- Documentation improvements
- Small bug fixes

Submit these as direct pull requests.

#### Large Changes
- New features
- Significant refactoring
- Breaking changes

Please open an issue first to discuss the approach.

## Pull Request Process

### Before Submitting

1. **Run Tests**
   ```bash
   make test
   make lint
   make fmt
   ```

2. **Update Documentation**
   - Update README.md if needed
   - Add/update command documentation
   - Update CHANGELOG.md

3. **Test Cross-Platform**
   - Test on your primary platform
   - Consider impacts on other platforms

### Pull Request Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change)
- [ ] New feature (non-breaking change)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] Tests added/updated
- [ ] CHANGELOG.md updated
```

### Review Process

1. **Automated Checks**: CI/CD pipeline runs tests and linting
2. **Code Review**: Maintainers review for quality and design
3. **Testing**: Manual testing on different platforms if needed
4. **Approval**: At least one maintainer approval required
5. **Merge**: Squash and merge with clean commit message

## Development Workflow

### Branch Naming

- `feature/description` - New features
- `fix/description` - Bug fixes
- `docs/description` - Documentation updates
- `refactor/description` - Code refactoring

### Local Development

```bash
# Create feature branch
git checkout -b feature/amazing-feature

# Make changes and test
make test
make lint

# Commit changes
git add .
git commit -m "feat(scope): add amazing feature"

# Push and create PR
git push origin feature/amazing-feature
```

### Testing

```bash
# Run all tests
make test

# Run specific test
go test ./internal/commands -v

# Run tests with coverage
make test-coverage

# Run integration tests
make test-integration

# Run linting
make lint

# Format code
make fmt
```

## Project Structure

```
rune/
├── cmd/rune/           # Main application entry point
├── internal/           # Private application code
│   ├── commands/       # CLI command implementations
│   ├── config/         # Configuration management
│   ├── rituals/        # Ritual execution engine
│   ├── tracking/       # Time tracking logic
│   └── integrations/   # External service integrations
├── pkg/                # Public API packages
├── docs/               # Documentation
├── scripts/            # Build and utility scripts
├── testdata/           # Test fixtures
└── examples/           # Example configurations
```

## Security Guidelines

### Reporting Security Issues

**DO NOT** open public issues for security vulnerabilities.

Instead, email: security@rune.dev

Include:
- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Security Best Practices

- **Input Validation**: Validate all user inputs
- **Command Execution**: Use secure command execution
- **Credentials**: Store in OS keychain, never in code
- **Dependencies**: Regular security audits
- **Permissions**: Principle of least privilege

## Documentation

### Types of Documentation

- **Code Comments**: For complex algorithms only
- **API Documentation**: Go doc comments for public APIs
- **User Documentation**: README, command help text
- **Developer Documentation**: Architecture decisions, setup guides

### Documentation Standards

- **Clear and Concise**: Easy to understand
- **Examples**: Include practical examples
- **Up-to-Date**: Keep in sync with code changes
- **Accessible**: Consider screen reader compatibility

## Community

### Communication Channels

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: General questions and ideas
- **Discord**: Real-time chat and community support
- **Email**: security@rune.dev for security issues

### Getting Help

- Check existing documentation first
- Search GitHub issues for similar problems
- Ask in GitHub Discussions for general questions
- Join our Discord for real-time help

### Recognition

Contributors are recognized in:
- CHANGELOG.md for each release
- README.md acknowledgments section
- GitHub contributor graphs
- Special badges for significant contributions

## Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):
- **MAJOR**: Breaking changes
- **MINOR**: New features, backwards compatible
- **PATCH**: Bug fixes

### Release Schedule

- **Patch releases**: As needed for critical bugs
- **Minor releases**: Monthly during active development
- **Major releases**: When significant breaking changes accumulate

## Questions?

If you have questions about contributing:

1. Check this document first
2. Search existing GitHub issues
3. Ask in GitHub Discussions
4. Join our Discord community

Thank you for contributing to Rune! 🎉

---

**Remember**: Every contribution, no matter how small, helps make Rune better for everyone.