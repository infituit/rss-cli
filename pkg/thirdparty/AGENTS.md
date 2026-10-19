|IMPORTANT: Prefer retrieval-led reasoning over pre-training-led reasoning
|Required Tools:codegraph (semantic code structure)|context7 (3rd-party docs)|sequential-thinking (decisions)
|Language Policy:Chinese for Q&A|English for code/docs/technical discussions
|Compression Rule:Follow ../../.agents/skills/rss-cli-scoped-agentsmd/references/AGENTS-compression-guide.md (pipe-index format, concise, no prose/code blocks)
|Scope:pkg/thirdparty
|Overview:外部服务协议 provider；不承载命令、context、输出或跨 provider 业务模型
|Where to look:FreshRSS provider:{pkg/thirdparty/freshrss/client.go|pkg/thirdparty/freshrss/types.go|pkg/thirdparty/freshrss/handler.go}:HTTP 协议、wire types、应用边界
|Where to look:FreshRSS tests:pkg/thirdparty/freshrss/client_test.go:transport stub 与 Fever 请求/响应契约
|Where to look:shared transport:pkg/rest/client/{client.go,request.go}:BaseURL、SubResourcef、query/form/context、Into
|Conventions:cli API 统一为 `ctx + req -> resp + err`；无协议参数也定义显式 request type；request/response wire types 保持 provider 内部
|Conventions:服务根 BaseURL 直接交给 REST client；协议资源 URI 在每个请求链用 `SubResourcef` 显式指定；输入已由上层标准化，provider 不做 TrimSpace 或路径重写
|Conventions:HTTP fluent chain 分行；先声明 request/form/response，再独立赋值 err；错误保留 provider action 上下文
|Conventions:新增协议测试覆盖 method、显式 URL、query/form 敏感字段、响应解析、非 2xx、非法 JSON
|Anti-patterns:不得把 Fever wire types 导出给 Handler、internal/cli、pkg/types 或输出层
|Anti-patterns:不得在 provider 拼接命令参数、解析 context、格式化 CLI 输出或创建跨 provider 抽象
|Anti-patterns:不得用隐式 BaseURL 路径、endpoint 重写或容错解析替代请求处显式协议 URI
|Dependencies:pkg/thirdparty/<service>→pkg/rest/client|pkg/types (when Handler converts application types)|不得依赖 cmd、internal 或 internal/cli
|Commands:gofmt -w ./pkg|go test ./pkg/thirdparty/freshrss|go test ./...|make lint|git diff --check
|Inheritance:遵守 pkg/AGENTS.md 与根级 AGENTS.md；本文件只补充 thirdparty provider 规则
