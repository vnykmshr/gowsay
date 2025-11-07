# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 2.0.x   | :white_check_mark: |
| < 2.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in gowsay, please report it responsibly.

### How to Report

1. **Do not** open a public GitHub issue for security vulnerabilities
2. Open a private security advisory on GitHub:
   - Go to the repository's Security tab
   - Click "Report a vulnerability"
   - Fill out the advisory form with details

### What to Include

- Description of the vulnerability
- Steps to reproduce the issue
- Potential impact
- Suggested fix (if available)

### Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Timeline**: Depends on severity
  - Critical: 7 days
  - High: 14 days
  - Medium/Low: 30 days

### Security Best Practices

When deploying gowsay:

1. **Authentication**: Always set `GOWSAY_TOKEN` in production
2. **HTTPS**: Deploy behind a reverse proxy with HTTPS
3. **Updates**: Keep dependencies up to date (enable Dependabot)
4. **Monitoring**: Monitor logs for suspicious activity
5. **Input Validation**: The application sanitizes input, but validate at the edge

### Known Security Considerations

- **Input Length**: Large input strings may consume memory. Configure reverse proxy limits.
- **Rate Limiting**: No built-in rate limiting. Implement at reverse proxy level.
- **Token Security**: Store tokens securely (environment variables, secrets manager).

### Dependencies

This project maintains minimal dependencies to reduce attack surface:
- `github.com/mattn/go-runewidth` - Unicode width calculation
- `github.com/mitchellh/go-wordwrap` - Text wrapping

All dependencies are scanned regularly via:
- Dependabot (automated updates)
- `govulncheck` (vulnerability scanning)
- GitHub Actions CI (on every PR)
