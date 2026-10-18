# 使用与配置

## 环境要求

- Go 1.25 或兼容版本，用于从源码运行和构建。
- 可访问 FreshRSS Fever API 的网络环境，用于认证检查。

## 构建与运行

```sh
go run ./cmd/rss-cli help
make build VERSION=v0.1.0-alpha.1
./build/v0.1.0-alpha.1/rss-cli version
```

跨平台构建：

```sh
make platform-builds VERSION=v0.1.0-alpha.1
make build-matrix VERSION=v0.1.0-alpha.1 GOOS=linux GOARCH=arm64
```

构建产物位于 `build/<VERSION>/`。`make build` 会注入版本、构建时间和 Git commit 信息。

UPX 压缩默认关闭。启用后，构建完成会自动压缩产物：

```sh
make build VERSION=v0.1.0-alpha.1 UPX_ENABLED=1
make platform-builds VERSION=v0.1.0-alpha.1 UPX_ENABLED=1 UPX_LEVEL=9
```

项目通过 `mise.toml` 固定 UPX 版本，启用压缩前运行 `mise install`。UPX 压缩支持 Linux `amd64`、`arm64` 和 Windows `amd64`；Darwin 以及 Windows `arm64` 会跳过压缩。缺少 `upx` 时构建会直接失败。

## 配置目录

配置目录按以下优先级选择：

```text
RSS_CLI_CONFIG_DIR
XDG_CONFIG_HOME/rss-cli
~/.config/rss-cli
```

例如：

```sh
RSS_CLI_CONFIG_DIR=/tmp/rss-cli rss-cli context list
XDG_CONFIG_HOME=/tmp/config rss-cli context list
```

实际结构为：

```text
<config-dir>/
├── current
└── contexts/<context>/
    ├── config.yaml
    └── services/<service>/
        ├── config.yaml
        └── credentials.enc
```

配置文件使用严格 YAML 解析，未知字段会被拒绝。敏感认证只写入应用层加密的 `credentials.enc`。

## Context

```sh
rss-cli context create home
rss-cli context use home
rss-cli context list
rss-cli context show
rss-cli context delete home --confirm
```

一个 context 可以配置多个服务，也可以只配置部分服务。同一个 context 内同一个服务最多一份认证。`--context` 可以临时覆盖 `current`：

```sh
rss-cli --context home context show
```

## FreshRSS

```sh
rss-cli freshrss config init --url https://freshrss.example.net
printf '%s' 'api-password' | rss-cli freshrss auth set --user alice --api-password-stdin
rss-cli freshrss auth check
rss-cli freshrss auth clear --confirm
```

FreshRSS 配置和认证分开管理。重复配置或认证默认失败，需要显式使用 `--force` 覆盖。密码和 API key 不作为普通命令行参数传递。

## 输出

默认输出 JSON：

```sh
rss-cli context list
rss-cli freshrss capabilities
```

需要表格时使用：

```sh
rss-cli --output table context list
```

成功结果写入 stdout；错误以包含 `code` 和 `message` 的 JSON 写入 stderr，并返回非零退出码。

## 常见错误

- `no current context`：先创建并选择 context，或使用 `--context`。
- `--url is required`：先运行 `freshrss config init --url URL`。
- `fever API key or username and API password are required`：补充有效凭据。
- `authentication failed`：检查地址、用户名和 FreshRSS API password。

示例中的密码和 API key 均为占位值；不要把真实凭据写入 shell 历史、仓库或问题报告。
