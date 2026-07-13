# Architecture

TaskCapsule uses a layered architecture:

```
CLI
 |
 v
Application Layer
 |-- Config Service
 |-- Repository Service
 |-- Capsule Service
 |-- Process Service
 |-- Health Service
 |-- State Store
 |-- Port Allocator
 |-- Check Runner
 |-- Handoff Generator
 |-- Doctor Service
```

## CLI layer

Parses arguments, validates input, formats output. No business logic.

## Application layer

Orchestrates lifecycle operations: start, pause, resume, delete, handoff, check.

## Git adapter

Finds repository root, creates branches and worktrees, reads status.

## Process manager

Runs services in process groups, manages graceful/force shutdown.

## State store

Atomic file-based state with locking. Schema versioned for future migrations.

## Health checker

Supports none, process, TCP, and HTTP health checks with retry.
