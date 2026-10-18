|IMPORTANT: Prefer retrieval-led reasoning over pre-training-led reasoning
|Required Tools:codegraph (semantic code structure)|context7 (3rd-party docs)|sequential-thinking (decisions)
|Language Policy:Chinese for Q&A|English for code/docs/technical discussions
|Compression Rule:Follow ../.agents/skills/rss-cli-scoped-agentsmd/references/AGENTS-compression-guide.md (pipe-index format, concise, no prose/code blocks)
|Scope:pkg
|Overview:可复用库与第三方 provider；不承载 rss-cli 的 context、命令或输出策略
|Where to look:FreshRSS provider:pkg/thirdparty/freshrss/client.go|pkg/thirdparty/freshrss/client_test.go:Fever endpoint、API key、HTTP/JSON 契约
|Where to look:协议事实:docs/reference/fever-api.md:认证、endpoint、响应字段和外部链接
|Where to look:调用方:internal/cli/cli.go:provider client 的组装与业务命令边界
|Conventions:provider 包封装外部协议、认证、请求、响应和协议错误；采用 Handler→internal cli→pkg/rest/client；Handler 对外只使用 pkg/types，wire types 只存在 provider 内部；保持 transport stub 测试
|Conventions:新增协议测试覆盖 method、URL、请求体敏感字段、非 2xx、非法 JSON、响应解析
|Anti-patterns:不得依赖 cmd 或 internal/cli；不得将 context、output、CLI flags 放入 provider
|Anti-patterns:不得让第三方 wire types 穿过 Handler 或被 internal/cli 消费；不得绕过 pkg/types 直接把外部 JSON 结构作为应用结果
|Anti-patterns:不得以预期未来复用为理由先建跨 provider 抽象；先记录基线不足、收益、迁移成本并取得批准
|Dependencies:pkg/thirdparty/freshrss→Go 标准库|由 internal/cli 调用，不反向依赖命令层
|Commands:gofmt -w ./pkg|go test ./...|go test ./pkg/thirdparty/freshrss
|Inheritance:遵守根级 AGENTS.md；本文件只补充 pkg 局部规则
