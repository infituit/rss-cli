# rss-cli

面向个人重度用户和 AI agent 的 RSS 服务命令行工具。`rss-cli` 通过独立的服务子命令连接不同 RSS 应用；当前提供 FreshRSS Fever API 的配置、认证和连通性检查能力。

## 快速开始

```sh
go run ./cmd/rss-cli help
RSS_CLI_CONFIG_DIR=/tmp/rss-cli go run ./cmd/rss-cli capabilities
```

创建 context 并配置 FreshRSS：

```sh
rss-cli context create home
rss-cli context use home
rss-cli freshrss config init --url https://freshrss.example.net
printf '%s' 'your-api-password' | rss-cli freshrss auth set --user alice --api-password-stdin=true
rss-cli freshrss auth check
```

密码和 API key 不作为普通命令行参数传递。示例中的地址、用户名和密码均为占位值。

## 命令输出

命令默认输出结构化 JSON。需要终端表格时使用 `--output table`：

```sh
rss-cli --output table context list
rss-cli --context home freshrss auth check
```

成功结果写入 stdout；错误以机器可读 JSON 写入 stderr，并使用非零退出码。

## 配置位置

配置目录按以下优先级选择：

```text
RSS_CLI_CONFIG_DIR
XDG_CONFIG_HOME/rss-cli
~/.config/rss-cli
```

目录结构为：

```text
<config-dir>/
├── current
└── contexts/<context>/
    ├── config.yaml
    └── services/<service>/
        ├── config.yaml
        └── credentials.enc
```

每个 context 可以配置多个服务，但同一 context 内同一服务最多保存一份认证。`XDG_DATA_HOME` 和 `XDG_CACHE_HOME` 当前未使用。

## 文档

- [文档导航](docs/README.md)
- [项目目标与边界](docs/concepts/project-goals.md)
- [使用与配置](docs/operation/usage.md)
- [架构设计](docs/developer/architecture/README.md)
- [代码目录](docs/developer/code-structure.md)
- [开发指南](docs/developer/development.md)
- [FreshRSS Fever API 参考](docs/reference/fever-api.md)

## 开发

```sh
mise install
go test ./...
make build
make lint
```
