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

Only the newest published release receives security fixes. The default branch and pull-request builds are development candidates, not supported releases.

## Known Security Properties

- Values named in `inheritEnvironment` are read at process start and are not persisted to capsule state
- Static values placed directly in `.taskcapsule.json` are part of the repository configuration and must not contain secrets
- Handoff reports redact likely secrets (API keys, tokens, passwords)
- On Unix, state directories use mode `0700` and state, log, check, and handoff files use mode `0600`; Windows storage inherits the current account's ACLs
- Capsule identifiers and loaded state are validated before filesystem operations
- Recursive worktree cleanup is restricted to the managed TaskCapsule worktree root
- No network services listen by default
- HTTP health checks can make requests to URLs explicitly provided by trusted repository configuration

## Threat Model

TaskCapsule assumes the current operating-system account and repository configuration are trusted. It executes configured commands and is not a sandbox for hostile code. The project protects against accidental path escape, malformed or inconsistent state, secret persistence through inherited environment values, concurrent lifecycle operations, and unsafe deletion boundaries.

An attacker who already controls the user's account, TaskCapsule binary, Git executable, or repository configuration can execute code with that user's permissions and is outside this threat model.
