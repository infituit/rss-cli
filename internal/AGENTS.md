|IMPORTANT: Prefer retrieval-led reasoning over pre-training-led reasoning
|Required Tools:codegraph (semantic code structure)|context7 (3rd-party docs)|sequential-thinking (decisions)
|Language Policy:Chinese for Q&A|English for code/docs/technical discussions
|Compression Rule:Follow ../.agents/skills/rss-cli-scoped-agentsmd/references/AGENTS-compression-guide.md (pipe-index format, concise, no prose/code blocks)
|Scope:internal
|Overview:rss-cli 私有实现；命令、配置、凭据、输出在此组装，外部协议留在 pkg/thirdparty
|Where to look:命令树:internal/cli/cli.go:全局参数、context、FreshRSS 命令和 provider 组装
|Where to look:配置:internal/config/config.go:XDG 路径、context/service YAML、current、严格字段校验
|Where to look:凭据:internal/credential/store.go:AES-256-GCM、machine-id 派生、credentials.enc、权限
|Where to look:输出:internal/output/output.go:JSON 默认、table、错误 envelope
|Where to look:架构与契约:docs/developer/architecture/README.md|docs/operation/usage.md:职责、调用链、公开命令
|Conventions:context 只管理环境和服务目录，不解析 provider 专属字段；service 配置与 credentials.enc 分离
|Conventions:internal/cli 只依赖 provider Handler 和 pkg/types；provider 的 wire request/response 类型与协议细节不得进入命令层
|Conventions:配置路径遵守 RSS_CLI_CONFIG_DIR→XDG_CONFIG_HOME/rss-cli→~/.config/rss-cli；未知 YAML 字段失败
|Conventions:配置、凭据、输出变更分别补边界/错误/脱敏测试，并同步读者文档
|Anti-patterns:CLI 不直接拼接第三方 API；凭据不得进入 YAML、日志、stdout、测试输出或文档
|Anti-patterns:不得在无真实复用证据和范围批准时建立跨 provider 领域模型
|Anti-patterns:不得直接实例化或调用 provider 内部 client；不得把第三方 API response 当作 CLI 输出模型
|Dependencies:internal/cli→internal/config|internal/credential|internal/output|pkg/thirdparty/freshrss
|Commands:gofmt -w ./internal|go test ./...|go build ./cmd/rss-cli
|Inheritance:遵守根级 AGENTS.md；本文件只补充 internal 局部规则
