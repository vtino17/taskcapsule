<<<<<<< HEAD
# Early Adopter Feedback

TaskCapsule is looking for feedback from developers who regularly switch between features, hotfixes, reviews, and experiments.

When opening an issue, include:

1. your operating system
2. the project stack
3. how many services the task starts
4. what interrupted the original task
5. which context was difficult to recover
6. the TaskCapsule command and output
7. whether the issue is reproducible in a small repository

Do not include:

- tokens
- API keys
- passwords
- private repository URLs
- confidential logs
- `.env` contents

Useful feedback categories:

- installation friction
- configuration confusion
- Git worktree behavior
- process shutdown or leaked children
- health-check behavior
- task-switching workflow gaps
- unclear CLI output
- Windows compatibility

A good report helps answer this question:

> What did you expect TaskCapsule to preserve or restore, and what happened instead?
=======
# Early Adopter Feedback Guide

## What we are looking for

Real-world feedback from developers who switch between features, hotfixes, reviews, and experiments.

## Questions for early adopters

1. What did your most recent task switch look like?
2. Which context took the longest to reconstruct?
3. What commands or tools do you currently use when switching tasks?
4. Did TaskCapsule help reduce the time to resume?

## Known limitations

- Windows process-tree management is experimental (no Job Objects yet)
- Port allocation has a small listen-close-bind race (no daemon in v0.1)
- Stale locks are detected but not automatically repaired
- HTTP health checks use fixed per-request timeout
- No daemon mode

## How to submit feedback

Open an issue at:
https://github.com/vtino17/taskcapsule/issues/new
>>>>>>> 0fbb5fc (docs: add campaign docs, update README with demo and badges)
