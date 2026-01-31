# Contributing to paqet

Thank you for considering contributing to paqet! This document provides guidelines and instructions for contributing.

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [How to Contribute](#how-to-contribute)
4. [Development Process](#development-process)
5. [Coding Standards](#coding-standards)
6. [Submitting Changes](#submitting-changes)
7. [Reporting Bugs](#reporting-bugs)
8. [Suggesting Features](#suggesting-features)
9. [Security Issues](#security-issues)

## Code of Conduct

This project adheres to a code of conduct that all contributors are expected to follow:

- Be respectful and inclusive
- Focus on constructive criticism
- Accept responsibility for mistakes
- Prioritize the community's best interest
- Show empathy toward other community members

## Getting Started

### Prerequisites

Before contributing, ensure you have:

1. **Go 1.25+** installed
2. **libpcap** development libraries installed
3. **Git** configured with your name and email
4. **GitHub account** for submitting pull requests

### Setting Up Development Environment

1. **Fork the repository** on GitHub

2. **Clone your fork:**
   ```bash
   git clone https://github.com/YOUR_USERNAME/paqet.git
   cd paqet
   ```

3. **Add upstream remote:**
   ```bash
   git remote add upstream https://github.com/hanselime/paqet.git
   ```

4. **Install dependencies:**
   ```bash
   go mod download
   ```

5. **Build the project:**
   ```bash
   go build -o paqet cmd/main.go
   ```

6. **Verify it works:**
   ```bash
   ./paqet --help
   ./paqet iface
   ```

## How to Contribute

### Types of Contributions We Welcome

- **Bug fixes** - Fix issues or unexpected behavior
- **New features** - Add new functionality
- **Documentation** - Improve or add documentation
- **Tests** - Add or improve test coverage
- **Performance** - Optimize existing code
- **Refactoring** - Improve code quality without changing behavior
- **Examples** - Add usage examples or tutorials

### Good First Issues

Look for issues labeled:
- `good first issue` - Good for newcomers
- `help wanted` - We need community help
- `documentation` - Documentation improvements

## Development Process

### 1. Create a Feature Branch

Always create a new branch for your work:

```bash
git checkout -b feature/descriptive-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions or changes

Examples:
- `feature/add-http-proxy`
- `fix/connection-timeout`
- `docs/improve-readme`

### 2. Make Your Changes

- Keep changes focused and atomic
- Write clear, descriptive commit messages
- Follow the coding standards (see below)
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Build the project
go build -o paqet cmd/main.go

# Run tests (if available)
go test ./...

# Test manually with test configuration
sudo ./paqet run -c test-config.yaml

# Check for formatting issues
go fmt ./...

# Run linter (if configured)
golangci-lint run
```

### 4. Commit Your Changes

Write clear commit messages following this format:

```
Short summary (50 chars or less)

More detailed explanation if needed. Wrap at 72 characters.
Explain what changed and why, not how.

- Bullet points are okay
- Use present tense ("Add feature" not "Added feature")
- Reference issue numbers: Fixes #123
```

Examples:
```
Add HTTP proxy support

Implements HTTP CONNECT proxy alongside existing SOCKS5 proxy.
This allows paqet to be used with applications that only support
HTTP proxies.

Fixes #45
```

### 5. Keep Your Branch Updated

```bash
# Fetch latest changes from upstream
git fetch upstream

# Rebase your branch on latest master
git rebase upstream/master

# If conflicts occur, resolve them and continue
git rebase --continue
```

### 6. Push to Your Fork

```bash
git push origin feature/descriptive-name
```

## Coding Standards

### Go Code Style

Follow standard Go conventions:

1. **Use `gofmt`** to format all code:
   ```bash
   go fmt ./...
   ```

2. **Follow Go naming conventions:**
   - Exported names start with uppercase: `PublicFunction`
   - Unexported names start with lowercase: `privateFunction`
   - Use camelCase, not snake_case
   - Acronyms should be uppercase: `ParseHTTP`, `WriteJSON`

3. **Write idiomatic Go:**
   - Use short variable names in short scopes
   - Return errors, don't panic
   - Accept interfaces, return structs
   - Keep functions small and focused

4. **Comment exported items:**
   ```go
   // Client manages connections to the paqet server.
   // It handles connection pooling and automatic reconnection.
   type Client struct {
       // ...
   }

   // Connect establishes a connection to the server.
   // It returns an error if the connection fails.
   func (c *Client) Connect() error {
       // ...
   }
   ```

5. **Handle errors properly:**
   ```go
   // Good
   if err != nil {
       return fmt.Errorf("failed to connect: %w", err)
   }

   // Avoid
   if err != nil {
       panic(err)  // Don't panic in library code
   }
   ```

### Project-Specific Guidelines

1. **Use the logging package:**
   ```go
   import "paqet/internal/flog"

   flog.Infof("Starting server on %s", addr)
   flog.Errorf("Failed to connect: %v", err)
   flog.Debugf("Received packet: %v", pkt)
   ```

2. **Configuration validation:**
   - All config structs must have `validate()` method
   - All config structs must have `setDefaults()` method
   - Use meaningful error messages

3. **Error handling:**
   - Use `fmt.Errorf` with `%w` to wrap errors
   - Provide context in error messages
   - Check errors immediately after operations

4. **Resource cleanup:**
   - Always defer `Close()` calls
   - Use context for cancellation
   - Clean up goroutines properly

### Documentation Standards

1. **Code comments:**
   - Comment all exported types, functions, and constants
   - Explain "why" not "what" in complex code
   - Keep comments up to date

2. **README updates:**
   - Update README.md if you add features
   - Add examples for new functionality
   - Keep documentation accurate

3. **Markdown formatting:**
   - Use proper heading hierarchy
   - Include code blocks with language tags
   - Add tables of contents for long docs

## Submitting Changes

### Pull Request Process

1. **Ensure your code:**
   - Builds successfully
   - Passes all tests
   - Follows coding standards
   - Includes appropriate documentation

2. **Create a pull request:**
   - Go to your fork on GitHub
   - Click "New Pull Request"
   - Select your feature branch
   - Fill in the PR template

3. **PR title format:**
   ```
   [Type] Short description

   Examples:
   [Feature] Add HTTP proxy support
   [Fix] Resolve connection timeout issue
   [Docs] Improve installation guide
   [Refactor] Simplify packet handling
   ```

4. **PR description should include:**
   - What changes were made
   - Why the changes were necessary
   - How to test the changes
   - Related issue numbers

### PR Template

```markdown
## Description
Brief description of changes

## Motivation
Why are these changes needed?

## Changes Made
- Change 1
- Change 2
- Change 3

## Testing
How were these changes tested?

## Checklist
- [ ] Code builds successfully
- [ ] Tests pass (if applicable)
- [ ] Documentation updated
- [ ] Follows coding standards
- [ ] Commit messages are clear

## Related Issues
Fixes #123
Related to #456
```

### Review Process

1. **Maintainers will review** your PR
2. **Address feedback** promptly
3. **Make requested changes** in new commits
4. **Once approved**, a maintainer will merge

## Reporting Bugs

### Before Reporting

1. **Search existing issues** to avoid duplicates
2. **Try the latest version** to see if it's already fixed
3. **Gather information** about your environment

### Bug Report Template

```markdown
## Bug Description
Clear description of the bug

## Steps to Reproduce
1. Step one
2. Step two
3. Step three

## Expected Behavior
What should happen

## Actual Behavior
What actually happens

## Environment
- OS: [e.g., Ubuntu 22.04]
- Go Version: [e.g., 1.25]
- paqet Version: [e.g., v1.0.0]
- Configuration: [relevant config snippet]

## Logs
```
Paste relevant logs here
```

## Additional Context
Any other relevant information
```

### Label Your Issue

Add appropriate labels:
- `bug` - Something isn't working
- `crash` - Application crashes
- `performance` - Performance issues
- `security` - Security vulnerabilities

## Suggesting Features

### Feature Request Template

```markdown
## Feature Description
Clear description of the proposed feature

## Motivation
Why is this feature needed?
What problem does it solve?

## Proposed Solution
How should this feature work?

## Alternatives Considered
What other approaches did you consider?

## Additional Context
Examples, mockups, or related features
```

### Discussion

- Be open to feedback
- Engage in constructive discussion
- Be willing to compromise
- Consider implementation complexity

## Security Issues

### Reporting Security Vulnerabilities

**DO NOT** open public issues for security vulnerabilities.

Instead:
1. Email the maintainers directly (see SECURITY.md if available)
2. Provide detailed information about the vulnerability
3. Allow time for a fix before public disclosure

### Security Best Practices

When contributing:
- Never commit secrets or credentials
- Validate all user input
- Use secure defaults
- Follow cryptography best practices
- Be aware of timing attacks
- Consider DoS scenarios

## Additional Guidelines

### Testing

While comprehensive tests may not yet exist, consider:
- Testing your changes manually
- Documenting test procedures
- Proposing test infrastructure improvements

### Performance

- Profile code for performance issues
- Avoid premature optimization
- Document performance characteristics
- Benchmark critical paths

### Dependencies

- Minimize new dependencies
- Use well-maintained libraries
- Keep dependencies updated
- Document why dependencies are needed

### Versioning

This project follows semantic versioning:
- MAJOR version for incompatible API changes
- MINOR version for new functionality
- PATCH version for bug fixes

## Getting Help

If you need help:
- Read the [Developer Guide](DEVELOPER_GUIDE.md)
- Check existing documentation
- Ask in GitHub Discussions
- Open an issue with the `question` label

## Recognition

Contributors will be:
- Listed in release notes
- Credited in the repository
- Acknowledged in documentation

Thank you for contributing to paqet! 🎉
