# `pkg/thirdparty`

## Overview

`pkg/thirdparty/` 存放外部服务集成和 SDK wrapper。每个外部服务使用独立 provider 包，外部服务的协议概念封装在对应包内。

## Structure

```text
pkg/thirdparty/
└── <provider>/         # 每个外部服务一个独立 provider 包
```

## Where To Look

| Task                       | Location                     | Notes                                              |
| -------------------------- | ---------------------------- | -------------------------------------------------- |
| Add an external service    | `pkg/thirdparty/<provider>/` | Keep the provider isolated in its own package      |
| Provider-specific behavior | `pkg/thirdparty/<provider>/` | Keep protocol concepts inside the provider package |

## Responsibilities

- 封装外部服务的认证、请求、响应和协议错误。
- 在 provider 包内维护外部服务的协议类型和转换。

## Boundaries

- provider 之间保持独立，不在一个包中混合多个外部服务。
- 外部服务 payload 和协议概念保持在对应 provider 包内。
- 每个 provider 包独立维护自身的协议适配，不在 `pkg/thirdparty` 中混合多个外部服务。
