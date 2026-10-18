# 代码目录

```text
rss-cli/
├── cmd/rss-cli/main.go                 # 进程入口
├── internal/
│   ├── cli/cli.go                       # 命令树与服务组装
│   ├── config/config.go                 # XDG、context 和 service 配置
│   ├── credential/store.go              # credentials.enc
│   └── output/output.go                 # JSON、table 和错误输出
├── pkg/
│   ├── rest/client/                     # shared REST request/result client
│   ├── types/                           # provider-neutral application types
│   ├── thirdparty/freshrss/             # Handler, provider cli and Fever wire types
│   └── version/                         # 构建版本元数据
├── testsuite/                           # opt-in FreshRSS 集成测试
├── docs/                                # 读者文档
├── AGENTS.md                            # 仓库级实现规则
└── go.mod
```

## 放置规则

- 新可执行程序放入 `cmd/<name>`。
- 项目私有命令、配置、凭据和输出放入 `internal`。
- 外部服务协议放入 `pkg/thirdparty/<service>`。
- provider 对外通过 Handler 暴露能力，内部按 `Handler -> cli -> pkg/rest/client` 分层。
- 跨 provider 和命令层的稳定应用类型放入 `pkg/types`；外部 wire types 只放在对应 provider 内，并必须显式转换。
- 只有已证明跨命令复用的通用库才放入其他 `pkg` 子目录。
- 文档放入对应的 `docs` 分区，不把实现规则写入读者文档。

## 命名规则

- 项目、二进制和根命令名称为 `rss-cli`。
- 服务使用顶层命令名，例如 `freshrss`。
- 配置覆盖变量为 `RSS_CLI_CONFIG_DIR`；XDG 变量使用标准名称。
- Go package 使用简短、小写、单一职责名称。
