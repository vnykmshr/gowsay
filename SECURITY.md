# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in gowsay, please report it responsibly.

**How to report:**
1. Open a private security advisory on GitHub (preferred)
2. Go to the repository's Security tab → "Report a vulnerability"
3. Include: description, steps to reproduce, and potential impact

I'll respond and address issues as quickly as I can.

## Security Considerations

When deploying gowsay:

- Set `GOWSAY_TOKEN` in production environments
- Deploy behind HTTPS (reverse proxy recommended)
- Consider rate limiting at the edge (application has basic input limits)
- Monitor logs for unusual patterns

## Dependencies

This project maintains minimal dependencies to reduce attack surface. Dependencies are monitored via:
- Dependabot (automated updates)
- `govulncheck` (vulnerability scanning in CI)
