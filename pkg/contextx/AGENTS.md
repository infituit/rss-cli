|IMPORTANT: Prefer retrieval-led reasoning over pre-training-led reasoning
|Required Tools:codegraph (semantic code structure)|context7 (3rd-party docs)|sequential-thinking (decisions)
|Language Policy:Chinese for Q&A|English for code/docs/technical discussions
|Compression Rule:Follow ../../.agents/skills/rss-cli-scoped-agentsmd/references/AGENTS-compression-guide.md (pipe-index format, concise, no prose/code blocks)
|Scope:pkg/contextx
|Overview:Shared request context model for rss-cli; extends Go context.Context with tenant_id, bk_username, login_name, message_id and OpenTelemetry span attributes
|Overview:Infrastructure-only; must not contain service or domain logic
|Parent AGENTS:pkg/AGENTS.md
|Structure:pkg/contextx:{iface.go core interfaces,context.go constructors and wrappers,info.go identity/value carrier and field validation,README.md package contract}
|Where to look:interfaces:iface.go:{IContext,IContextValues,IContextInfo}
|Where to look:lifecycle:context.go:{New,From,FromContext,WithCancel,WithTimeout,WithDeadline,WithoutCancel,Background}
|Where to look:identity:info.go:{Info,CheckTenantID,CheckBKUsername,CheckLoginName,CheckMessageID}
|Conventions:Prefer pkg/contextx.IContext at service/internal boundaries when tenant/user/message metadata is required
|Conventions:Create or derive contexts through New, From, or FromContext instead of ad-hoc struct construction
|Conventions:Use contextx.WithTimeout, contextx.WithDeadline, or contextx.WithCancel when preserving existing context values
|Conventions:Run explicit identity guards at validation boundaries instead of assuming fields exist
|Conventions:Keep span-attribute behavior centralized in context.go; avoid duplicate tracing field injection in callers
|Conventions:Keep map propagation and clone semantics stable; avoid mutation patterns that corrupt shared context values
|Dependencies:Go stdlib context|time; OpenTelemetry go.opentelemetry.io/otel/{trace,attribute}
|Anti-patterns:Do not add service-specific metadata keys or business validation rules in this package
|Anti-patterns:Do not bypass this package with raw context.Context where tenant/user identity is mandatory
|Anti-patterns:Do not use plain context.WithTimeout or context.WithDeadline on IContext paths when values must be preserved
|Anti-patterns:Do not expose mutable internal state that lets callers corrupt shared context value semantics
|Anti-patterns:Do not duplicate context models in other packages; extend this package only when fields are truly shared
|Verification:go test ./pkg/contextx/...
