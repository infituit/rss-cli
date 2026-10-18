# 架构设计

## 运行链路

```text
cmd/rss-cli/main.go
        |
        v
internal/cli
   |          |
   v          v
internal/  internal/
config     credential/output
   |
    v
pkg/thirdparty/freshrss
        |
        v
   provider Handler
        |
        v
   provider cli
        |
        v
pkg/rest/client
        |
        v
FreshRSS Fever API
```

## 模块职责

### `cmd/rss-cli`

只负责进程入口、标准 I/O 和退出码，不放业务逻辑、配置解析或 HTTP 请求。

### `internal/cli`

负责全局参数、context 命令、服务命令、帮助和 capabilities。它把服务配置和凭据组装成 provider client，并将结果交给统一输出层。

### `internal/config`

负责 XDG 配置目录解析、context/service YAML 文件、当前 context 和严格字段校验。context 管理只理解环境和服务目录，不理解 provider 专属认证字段。

### `internal/credential`

负责 `credentials.enc` 的应用层加密、解密和文件权限。凭据不进入普通 YAML、日志或输出。

### `internal/output`

负责 JSON 默认输出、table 输出和 stderr 错误 envelope。

### `pkg/types`

定义跨 provider 和命令层使用的应用类型。它不携带 FreshRSS JSON tag 以外的 provider 协议细节；所有外部请求和响应都必须在 provider 边界转换为或转换自这些类型。

### `pkg/rest/client`

提供共享 REST transport 和链式 request/result 能力，包括 context、headers、query、form、JSON、重试、延迟观测和敏感字段脱敏。它不理解 RSS 业务或具体 provider 类型。

### `pkg/thirdparty/freshrss`

采用 `Handler -> provider cli -> pkg/rest/client` 分层。Handler 对外只暴露 `pkg/types`；provider 内部的 Fever wire request/response 负责 JSON 映射，cli 负责请求组装和协议错误转换，Handler 负责最终类型转换。该包不依赖 CLI 或输出实现。

## 依赖方向

```text
cmd -> internal/cli
internal/cli -> internal/config
internal/cli -> internal/credential
internal/cli -> internal/output
internal/cli -> pkg/thirdparty/freshrss
pkg/thirdparty/freshrss -> pkg/rest/client
pkg/thirdparty/freshrss -> pkg/types
```

`pkg/thirdparty` 不得反向依赖 `cmd` 或 `internal/cli`。不同服务的协议差异保留在各自 provider 包中；只有出现真实复用证据时才新增跨服务领域抽象。

### 类型转换边界

- 发往外部服务的数据：`pkg/types` -> provider wire request -> REST client。
- 外部服务返回的数据：REST client -> provider wire response -> `pkg/types`。
- `internal/cli`、输出层和其他 provider 不得依赖 Fever wire types。
- provider 特有字段只能在 provider 类型中保留；只有具有稳定跨服务语义的字段才进入 `pkg/types`。

## 错误处理

- 业务命令返回 Go error，由入口转换为机器可读 JSON 并写入 stderr。
- 成功结果只写入 stdout。
- provider 负责保留 HTTP 状态和协议解析上下文。
- 凭据解密失败不得回显密文或认证内容。
