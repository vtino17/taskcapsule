<<<<<<< HEAD
# TaskCapsule Launch Kit

This file contains ready-to-publish launch copy. Replace `[REPOSITORY LINK]` and `[DEMO VIDEO]` before posting.

## Core positioning

**One sentence**

> TaskCapsule lets developers pause and resume complete coding-task contexts using isolated Git worktrees and managed local services.

**Tagline**

> Pause one coding task. Handle the interruption. Resume without losing your place.

**Three-line pitch**

Task switching costs more than a branch checkout. Developers also need to reconstruct services, ports, logs, failing checks, and the thought they were holding before the interruption. TaskCapsule packages those pieces into a local, resumable capsule.

## Show HN

**Title**

```text
Show HN: TaskCapsule – pause and resume coding tasks with Git worktrees
```

**Body**

```text
I built TaskCapsule because switching from a feature to an urgent hotfix kept costing me more context than the Git branch itself.

I had to remember which development services were running, which ports they used, which check failed, and what I planned to do next.

TaskCapsule groups a Git worktree, local processes, logs, notes, and check results into one resumable task:

    taskcapsule start payment-timeout
    taskcapsule note payment-timeout "Investigate duplicate retry"
    taskcapsule pause payment-timeout

    taskcapsule start urgent-hotfix

    taskcapsule resume payment-timeout
    taskcapsule where payment-timeout

It is written in Go and works locally without a daemon, cloud account, API key, or AI model. It also never automatically commits, stashes, resets, or pushes code.

Linux and macOS are fully supported in v0.1. Windows process-tree management is still experimental.

I would especially appreciate feedback about the lifecycle model, configuration format, and whether this fits a real workflow you have.

[REPOSITORY LINK]
```

## Reddit

**Suggested title**

```text
I built a Go CLI to pause and resume complete coding-task contexts
```

**Body**

```text
I kept losing context when an urgent task interrupted a feature.

Checking out the old branch was not the hard part. The hard part was reconstructing the worktree, development processes, ports, logs, latest failed check, and the thought I was holding before I switched.

I built TaskCapsule to manage those pieces as one local task:

    taskcapsule start feature-checkout
    taskcapsule note feature-checkout "Fix retry assertion next"
    taskcapsule pause feature-checkout

    # handle another task

    taskcapsule resume feature-checkout
    taskcapsule where feature-checkout

It uses Git worktrees and local process management. There is no cloud account, daemon, API key, or AI dependency. It deliberately does not auto-commit, stash, reset, merge, or push.

Current limitations are documented: Windows process trees are experimental, port reservation has a small race, and stale locks are detected but not auto-repaired.

I am looking for feedback from people who frequently switch between features, hotfixes, and reviews. What part of your context-switching workflow is still manual?

[DEMO VIDEO]
[REPOSITORY LINK]
```

## X / Twitter

```text
Switching branches is easy.

Reconstructing the task is not:
- dev servers
- ports
- logs
- failing checks
- the thought you were holding

I built TaskCapsule to pause and resume that entire local context.

No cloud. No daemon. No API keys. No AI.

[DEMO VIDEO]
[REPOSITORY LINK]
```

**Follow-up post**

```text
taskcapsule pause feature-a
taskcapsule start urgent-hotfix

# later
taskcapsule resume feature-a
taskcapsule where feature-a

It uses isolated Git worktrees and managed local processes, while refusing to auto-commit, stash, reset, or push.
```

## LinkedIn

```text
An urgent hotfix interrupted a feature I was building.

Returning to the branch later was easy. Recovering the real working context was not.

Which services were running? Which ports did they use? Which test had failed? What was the exact next step I had in mind?

That problem led me to build TaskCapsule, an open-source Go CLI that turns a coding task into a resumable local capsule:

- isolated Git worktree
- development processes
- ports and health checks
- logs and latest check result
- a note about where you left off
- a secret-safe handoff report

The basic workflow is:

    taskcapsule start feature-a
    taskcapsule pause feature-a
    taskcapsule start urgent-hotfix
    taskcapsule resume feature-a
    taskcapsule where feature-a

TaskCapsule runs locally without a cloud account, daemon, API key, or AI model. It also never automatically commits, stashes, resets, or pushes source code.

The first public release is available now. Feedback from developers who regularly juggle features, hotfixes, and reviews would be extremely useful.

[DEMO VIDEO]
[REPOSITORY LINK]
```

## Product Hunt

**Name**

```text
TaskCapsule
```

**Tagline**

```text
Hibernate coding tasks, not your computer
```

**Short description**

```text
Pause and resume Git worktrees, development services, logs, notes, and checks as one local task.
```

**Maker comment**

```text
I built TaskCapsule after repeatedly losing working context when a feature was interrupted by a hotfix or review.

The Git branch was rarely the real problem. I also needed to recover services, ports, logs, the latest failed check, and my next intended action.

TaskCapsule manages that lifecycle locally and conservatively. It has no cloud or AI dependency and never automatically commits, stashes, resets, or pushes code.

I am launching the first version to learn which parts of developer task switching are still painful in real projects.
```

## Demo script

Target duration: 15–20 seconds.

```bash
taskcapsule start feature-checkout
taskcapsule note feature-checkout "Fix retry assertion next"
=======
# Launch Kit

## One-liner

Pause one coding task. Handle the interruption. Resume without losing your place.

## Positioning

TaskCapsule is a local-first CLI that groups Git worktree, development processes, ports, logs, notes, checks, and handoff into one resumable coding task.

## Key facts

- No cloud account required
- No daemon process
- No API keys
- No AI dependency
- No automatic Git mutations (commit, stash, reset, merge, push)
- Written in Go
- Linux and macOS fully supported
- Windows process-tree management is experimental
- Apache 2.0 license

## Target audience

Software developers who frequently switch between features, hotfixes, reviews, and experiments.

## Problem statement

Switching Git branches is easy. Reconstructing the complete coding task is not:
- Which dev servers were running?
- Which ports?
- What was the last test result?
- What was I about to do next?

## Solution

```
start → work → pause → switch → resume → handoff → delete safely
```

## Example

```bash
taskcapsule start feature-checkout
taskcapsule note feature-checkout "Continue retry test next"
>>>>>>> 0fbb5fc (docs: add campaign docs, update README with demo and badges)
taskcapsule pause feature-checkout

taskcapsule start urgent-hotfix
taskcapsule pause urgent-hotfix

taskcapsule resume feature-checkout
taskcapsule where feature-checkout
```

<<<<<<< HEAD
Suggested overlay:

```text
2 coding tasks
0 manual branch switching
0 forgotten context
```

The recording must show real command output from the released binary. Do not fake process restoration or hide errors with editing.

## Launch-day response principles

- Answer technical questions directly.
- Do not ask people to star the repository.
- Ask commenters about their real workflow.
- Acknowledge known limitations without becoming defensive.
- Explain that TaskCapsule coordinates Git worktrees and processes rather than replacing Git, tmux, or Docker Compose.
- Turn recurring questions into README improvements or focused issues.

## Metrics that matter

For the first launch, prioritize:

- successful installations
- capsules created by external users
- reproducible bug reports
- workflow feedback
- first external contributor

Stars are useful distribution signals, but they are not proof that the product is solving the problem.
=======
## Commands

| Command     | Description                          |
|-------------|--------------------------------------|
| `init`      | Create `.taskcapsule.json`           |
| `start`     | Create capsule, worktree, services   |
| `pause`     | Stop services, release resources     |
| `resume`    | Restart services, restore context    |
| `list`      | List capsules                        |
| `status`    | Show detailed capsule state          |
| `note`      | Save context note                    |
| `where`     | Show where you left off              |
| `check`     | Run and record a validation command  |
| `logs`      | View service logs                    |
| `handoff`   | Generate Markdown handoff report     |
| `delete`    | Remove capsule and worktree safely   |
| `doctor`    | Diagnose local state                 |
| `version`   | Show build information               |

## Channels

- GitHub: https://github.com/vtino17/taskcapsule
- X (primary): https://x.com/cyberfix17/status/2076837216350048637
- Feedback issue: https://github.com/vtino17/taskcapsule/issues/1
>>>>>>> 0fbb5fc (docs: add campaign docs, update README with demo and badges)
