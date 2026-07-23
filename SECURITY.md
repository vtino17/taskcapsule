# Security Policy

## Reporting a Vulnerability

To report a security vulnerability in TaskCapsule, please open a private security advisory on GitHub:

https://github.com/vtino17/taskcapsule/security/advisories/new

Do not open a public issue for security vulnerabilities.

## What to Include

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if known)

## Supported Versions

| Version | Supported |
|---------|-----------|
| latest  | Yes |
| older   | No |

## Known Security Properties

- TaskCapsule never persists environment variable values
- Handoff reports redact likely secrets (API keys, tokens, passwords)
- State files use restrictive permissions (0600)
- No network services listen by default
- No external API calls occur
