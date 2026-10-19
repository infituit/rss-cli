---
name: rss-cli-how-to
description: Use when planning, implementing, reviewing, or debugging rss-cli changes and deciding which project skills, scoped AGENTS.md files, and language-level guidance to load.
---

# rss-cli How-To

`rss-cli-how-to` is the project-first router for this repository. Identify the
owned surface first, load its nearest `AGENTS.md` and narrow project skill,
then add only the companion guidance needed for the language or tool involved.

## When to Use

Use this skill when a task spans more than one project surface, asks which
skills to load, or could be mistaken for a generic Go, CLI, HTTP, testing, or
documentation task.

Do not use it instead of a narrower skill. A narrow skill and the nearest
scoped `AGENTS.md` still own their contracts.

## Routing Rules

| Task signal                                                                                          | Load first                                                      | Add when needed                                                |
| ---------------------------------------------------------------------------------------------------- | --------------------------------------------------------------- | -------------------------------------------------------------- |
| `cmd/**`, command tree, global flags, stdout/stderr, exit status                                     | `cmd/AGENTS.md`                                                 | `rss-cli-audience-boundary` for user-facing docs               |
| `internal/cli`, context/service assembly, command behavior                                           | `internal/AGENTS.md`                                            | `rss-cli-audience-boundary` for command docs                   |
| config paths, context/service YAML, strict decoding, `current`                                       | `internal/AGENTS.md`                                            | security/testing guidance for boundary changes                 |
| credentials, encryption, permissions, redaction                                                      | `internal/AGENTS.md`                                            | security guidance; stop on credential-boundary changes         |
| `pkg/thirdparty/freshrss`, Fever HTTP/JSON, transport stubs                                          | `pkg/thirdparty/AGENTS.md` then `rss-cli-thirdparty`            | `pkg/AGENTS.md`; current official HTTP/Go docs when unfamiliar |
| `testsuite/support`, Docker Compose FreshRSS fixture, integration tag                                | `rss-cli-testsuite-support`                                     | generic testing guidance only for a concrete gap               |
| README, `docs/**`, Issue, PR, or artifact placement                                                  | `rss-cli-audience-boundary`, `rss-cli-information-architecture` | `rss-cli-create-pr` for PR delivery                            |
| scoped `AGENTS.md` creation or compression                                                           | `rss-cli-scoped-agentsmd`                                       | nearest parent and child rules                                 |
| PR preparation, creation, update, replacement, or split                                              | `rss-cli-create-pr`                                             | affected surface skill and review guidance                     |
| cross-layer placement, dependency direction, shared abstraction, public contract, architecture drift | `rss-cli-architecture-judgment`                                 | affected surface skill; ask before an unapproved departure     |

## Core Pattern

1. Identify changed paths and the public behavior they affect.
2. Read root and nearest scoped `AGENTS.md` files.
3. Inspect representative local implementations and existing skills.
4. Load the smallest project-skill set covering the surface; use generic language guidance only as a companion.
5. Verify the affected contract, adjacent tests, build, and documentation surface.

For a single-file mechanical edit, the target file and nearest rules may be
enough. Do not load every skill by default.

## Project Boundaries

- `internal` owns CLI assembly, configuration, credentials, and output policy.
- `pkg/thirdparty` owns external provider protocol and transport behavior.
- `cmd` owns process entry and exit, not business or protocol logic.
- Real FreshRSS integration uses the project fixture; unit tests use transport stubs.
- Do not create a cross-provider model, generic helper, or future command without demonstrated reuse and approved scope.

## Verification Contract

Choose checks from the changed surface and report only checks actually run:

- Go implementation: `go test ./...` and `go build ./cmd/rss-cli`.
- Formatting: `gofmt` on affected Go paths and `git diff --check`.
- Integration fixture: `go test -tags=integration ./testsuite/... -run '^$' -count=1`; run Docker integration only when available.
- Docs/config/command changes: compare commands, paths, config names, output, and protocol facts with implementation.

## Common Mistakes

- Loading only generic Go advice while ignoring `AGENTS.md` and provider boundaries.
- Loading every project skill instead of the smallest matching set.
- Treating `pkg/thirdparty` as the owner of CLI flags, output formatting, or context policy.
- Starting Docker from the default test suite or putting credentials in fixtures, logs, docs, or examples.
- Updating the router merely because a skill exists; update routing only when ownership or routing semantics change.

## Evals

Pressure scenarios live in `evals/evals.json`. Keep execution output outside the
repository unless explicitly requested.

## Cross-References

- `rss-cli-testsuite-support`: real FreshRSS fixture and integration boundary.
- `rss-cli-architecture-judgment`: pre-flight/post-flight architecture judgment and boundary drift checks.
- `rss-cli-thirdparty`: provider capabilities, raw API contracts, Handler/client separation, and shared client reuse.
- `rss-cli-information-architecture`: durable artifact placement and sync.
- `rss-cli-audience-boundary`: reader and artifact content limits.
- `rss-cli-scoped-agentsmd`: subtree rule generation and compression.
- `rss-cli-create-pr`: delivery and remote-write gates.
