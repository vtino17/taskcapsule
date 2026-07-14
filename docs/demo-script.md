<<<<<<< HEAD
# TaskCapsule 20-Second Demo

Use the released binary and a real temporary Git repository.

## Terminal size

- 100–120 columns
- large readable font
- dark background
- no personal paths, tokens, or private repository names

## Recording sequence

```bash
# Task 1: normal feature work
taskcapsule start checkout-retry
taskcapsule note checkout-retry "Fix the retry assertion next"
taskcapsule pause checkout-retry

# An interruption arrives
taskcapsule start urgent-hotfix
taskcapsule pause urgent-hotfix

# Return to the original task
taskcapsule resume checkout-retry
taskcapsule where checkout-retry
```

## On-screen overlays

Opening:

```text
A hotfix interrupted your feature.
```

Middle:

```text
Pause the whole task, not just the branch.
```

Ending:

```text
TaskCapsule
Resume without losing your place.
```

## What must be visible

- a capsule entering `running`
- services stopping during `pause`
- the second capsule starting
- the first capsule resuming
- the saved note appearing in `where`

## Export

Create both:

- MP4 for X, LinkedIn, Reddit, and Product Hunt
- GIF under 8 MB for the README

Place the final README asset at:

```text
assets/taskcapsule-demo.gif
```

Then add it below the badges in `README.md`:

```markdown
![TaskCapsule demo](assets/taskcapsule-demo.gif)
```
=======
# Demo Script

## Workflow

This demo shows the core TaskCapsule lifecycle: interrupt a feature with a hotfix, then resume the original task.

```bash
# Step 1: Start working on a feature
taskcapsule start feature-checkout
taskcapsule note feature-checkout "Continue retry test next"

# Step 2: Pause when interrupted by a hotfix
taskcapsule pause feature-checkout

# Step 3: Handle the urgent interruption
taskcapsule start urgent-hotfix
taskcapsule pause urgent-hotfix

# Step 4: Resume the original task
taskcapsule resume feature-checkout
taskcapsule where feature-checkout
```

## What to show visually

1. Terminal window, ~80x24 columns
2. Run each command with a ~1s pause between steps
3. Show that `where` reconstructs the task context
4. No private paths, emails, tokens, or credentials visible

## Duration

15–25 seconds total.

## Format

Terminal GIF (preferred), MP4, or asciinema.
>>>>>>> 0fbb5fc (docs: add campaign docs, update README with demo and badges)
