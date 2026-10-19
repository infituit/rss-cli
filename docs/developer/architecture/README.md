# 架构设计

## 基本模块构成

```text
rss-cli/
├── cmd/                # 进程入口
├── internal/           # 命令、配置、凭据、输出和 telemetry 组装
└── pkg/                # 可复用库、协议 client 和 provider
```

## 模块职责

### `cmd`

负责进程入口、标准 I/O 和退出码，不放业务逻辑、配置解析或 HTTP 请求。

### `internal`

负责命令、配置、凭据、输出和进程内服务组装。它消费 `pkg` 提供的能力，不向 `pkg` 暴露自身实现。

### `pkg`

提供可复用库、外部服务 provider 和项目共享类型。`pkg` 不依赖 `internal`。

## 依赖关系

```text
cmd -> internal
internal -> pkg
```

`internal` 消费 `pkg` 提供的能力；`pkg` 完全不知道 `internal` 的存在。一级目录内部的包依赖由各自目录规则和代码负责，不在本架构文档中固化。

## Package Architecture

- [`pkg` 一级包](pkg/README.md)：`pkg` 下一级包的定位和边界。
- [`pkg/thirdparty`](pkg/thirdparty.md)：外部服务 provider 的目录结构、职责和边界。
