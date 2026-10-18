|IMPORTANT: Prefer retrieval-led reasoning over pre-training-led reasoning
|Required Tools:codegraph (semantic code structure)|context7 (3rd-party docs)|sequential-thinking (decisions)
|Language Policy:Chinese for Q&A|English for code/docs/technical discussions
|Compression Rule:Follow `.agents/skills/rss-cli-scoped-agentsmd` (pipe-index format, concise, no prose/code blocks)
|Scope:.
|Overview:repository-wide implementation, documentation, testing, and agent-workflow rules; deeper AGENTS.md adds local constraints without weakening this file
|Precedence:user instructions→repository rules→verified public contracts→local conventions→general best practice
|Retrieval:before edits read nearest AGENTS.md, inspect target and representative local pattern; use codegraph when indexed, otherwise focused Glob/Grep/Read; unfamiliar third-party API requires current official Context7 docs
|Memory:session start recall relevant project decisions/file history; task skill load→same-scope lesson recall; save durable decisions, conventions, security findings, and lessons; never save credentials or transient logs
|Skill Routing:for multi-surface or skill-selection work read `.agents/skills/rss-cli-how-to/SKILL.md`; load the narrowest project skill and nearest scoped AGENTS.md before generic guidance
|Planning:state material assumptions; high-cost changes require blind-spot pass over user choices, public contract, data model, security boundary, irreversible operations, and user-visible behavior; define verifiable success criteria
|Design:reduce high-cost work to invariants/constraints/success criteria; prefer smallest baseline-preserving change; new dependency/shared abstraction/cross-layer organization/public contract needs evidence and approval when not already agreed
|Deviation:record pain, measurable improvement, maintenance/migration cost, and why baseline is insufficient; ask before changing public contract, data model, security boundary, long-lived architecture, irreversible operation, or ambiguous user-visible behavior
|Change discipline:use apply_patch for manual edits; preserve unexpected worktree/staging changes; never revert unrelated work; no destructive Git commands; no commit/amend/push/PR/remote/repository/folder rename without exact user authorization
|Validation failure:inspect exact formatter/parser/schema/validator error and current state; retry only with a changed hypothesis; never bypass validation or silently broaden scope
|Product:rss-cli project/binary/root command; module `github.com/infituit/rss-cli`; FreshRSS is first service command, not product identity; no forced cross-service business model without demonstrated reuse
|Context:environment may contain multiple services and only a subset; one service has at most one credential set per context; service-specific logic stays outside generic context management
|Configuration:precedence `RSS_CLI_CONFIG_DIR`→`XDG_CONFIG_HOME/rss-cli`→`~/.config/rss-cli`; `XDG_DATA_HOME`/`XDG_CACHE_HOME` reserved; context/service YAML separate and strict; unknown fields fail
|Credentials:only application-encrypted `credentials.enc` with restrictive permissions; never expose secrets in YAML, output, logs, tests, docs, or examples; examples use placeholders
|CLI contract:global `--context`/`--output` after root command; context overrides `current`; default JSON, optional table; success stdout; machine-readable JSON errors stderr; non-zero failure status
|Documentation:README stable entry point; `docs/concepts` goals/boundaries; `docs/operation` user workflows; `docs/developer` implementation/contribution; `docs/reference` external protocol facts; `.agents/skills` repeatable agent workflows
|Placement:before long-lived content record target, responsibility, excluded neighbors, evidence, synchronization, and risk; never use docs/.agents as temporary notes
|Reader boundary:reader-facing artifacts contain verified current behavior, limits, links, risks, and next actions; exclude grilling, investigation logs, temporary plans, unfinished proposals, and author process
|Consistency:check commands, paths, configuration names, output examples, protocol facts, naming, errors, tests, docs, and scoped AGENTS.md against implementation before delivery; report departures
|Review:review findings first by severity with exact file/line references; delegation does not transfer responsibility; verify subagent output against current tree
|Commands:gofmt -w ./cmd ./internal ./pkg ./testsuite|go test ./...|make build|make lint|git diff --check
|Integration:for integration changes compile `go test -tags=integration ./testsuite/... -run '^$' -count=1`; run real opt-in Docker test only when available; report unavailable checks truthfully
|Child AGENTS:cmd/AGENTS.md|internal/AGENTS.md|pkg/AGENTS.md|docs/AGENTS.md
