# AI Agent Instructions

## Source of truth

Read these files before making changes:

1. AGENTS.md
2. PLAN.md (Product Requirements, Architecture, Implementation Plan)

## Hard constraints

- Do not add AI features.
- Do not add cloud dependencies.
- Do not automatically commit, stash, reset, rebase, merge, or push.
- Do not store environment variable values.
- Do not delete dirty worktrees without explicit force.
- Do not delete Git branches.
- Do not use shell evaluation for configured commands.
- Every lifecycle change requires an integration test.
- Every state schema change requires schema-version handling.
- Linux and macOS are the primary platforms for v0.1.

## Work sequence

1. Read the current milestone.
2. Implement only tasks in that milestone.
3. Add tests.
4. Run format, vet, test, race, and build.
5. Update documentation.
6. Stop after the milestone acceptance criteria pass.
