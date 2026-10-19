---
name: rss-cli-architecture-judgment
description: Use when making or reviewing rss-cli architecture decisions, especially cmd/internal/pkg placement, new pkg helpers or abstractions, provider boundary changes, dependency direction, public Handler changes, or architecture drift self-review.
---

# rss-cli Architecture Judgment

## Overview

Use this skill as an advisory guardrail before or after risky architecture work. It makes the decision and evidence explicit; it does not authorize code changes or replace the narrow skill for the affected surface.

## When To Use

- `pre-flight`: before cross-layer work, new shared helpers, interfaces, adapters, or public Handler changes.
- `post-flight`: after implementation or when reviewing dependency and boundary drift.
- Use for dependency direction, shared contracts, trust boundaries, or explicit architecture scans.

Do not use it for ordinary single-package implementation, Go formatting, naming, lint, or error handling without an architecture signal.

## Evidence First

1. Read the root and nearest scoped `AGENTS.md` files.
2. Inspect 2-3 representative implementations in the affected `cmd`, `internal`, or `pkg` surface.
3. Check the actual import direction and current caller flow.
4. Report missing evidence instead of inventing a project rule.

## Verdict

Emit one verdict. `Stop` takes precedence over `Warn`, which takes precedence over `Continue`.

| Verdict  | Use when                                                                                                    | Required action                                                         |
| -------- | ----------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| Stop     | A scoped rule, public contract, trust boundary, or dependency direction is violated or unresolved           | Do not implement; state the blocker and one decision needed             |
| Warn     | The change is reversible but introduces an under-proven abstraction, seam, duplication, or terminology risk | Keep the smallest approved path and record the verification requirement |
| Continue | Local evidence and scoped rules support the existing path with no unresolved architecture risk              | Proceed without broadening the scope                                    |

## Output Contract

```text
Mode: pre-flight | post-flight
Verdict: Stop | Warn | Continue
Decision: [architecture decision]
Evidence: [rules, imports, and local references]
Watch points: [concrete risks]
Required next action: [none, one question, or smallest adjustment]
Side effects: none
```

For `post-flight`, lead with findings and classify drift as contract, layer, dependency, scope, terminology, side-effect, or verification drift.

## Boundary Rules

- `cmd` owns process entry and exit.
- `internal` owns command assembly, configuration, credentials, output, and private orchestration.
- `pkg` owns reusable libraries and provider protocol behavior; `pkg` must not depend on `internal`.
- `pkg/thirdparty` owns external-service adapters; do not move command policy into it.
- A new shared helper, interface, or abstraction needs a second real use or explicit approval.

## Common Mistakes

- Treating this skill as permission to refactor.
- Hard-coding dependencies between second-level packages in the top-level architecture contract.
- Adding a `pkg` abstraction for one provider operation.
- Loading this skill instead of the narrower provider or CLI skill.

## See Also

- `rss-cli-how-to`
- `rss-cli-thirdparty`
- `rss-cli-scoped-agentsmd`
