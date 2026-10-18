# 开发指南

## 本地验证

```sh
mise install
gofmt -w ./cmd ./internal ./pkg ./testsuite
go test ./...
make build
make lint
```

集成测试需要 Docker：

```sh
make integration-test
```

它是显式 opt-in 测试，会启动隔离的 FreshRSS Compose 实例并在结束时清理容器和卷。

## 新增服务能力

1. 在 `pkg/types` 确认或新增稳定的应用请求/响应类型；不要把 provider JSON tag 和协议 envelope 泄漏到公共类型。
2. 在 `pkg/thirdparty/<service>` 按 `Handler -> cli -> pkg/rest/client` 实现 provider；请求和响应 wire types 只在 provider 内部使用，并补齐双向转换与 transport stub 测试。
3. 在 `internal/cli` 增加该服务顶层命令，只消费 Handler 和 `pkg/types`。
4. 复用 `internal/config` 的 context/service 存储，不把 provider 字段放入 context 逻辑。
5. 使用 `internal/output` 输出 JSON、table 和错误结果。
6. 更新 `README.md`、`docs/README.md` 及对应 `docs/operation`、`docs/reference` 文档。
7. 运行格式化、单元测试、构建和 lint。

不同服务不要求采用相同的业务资源模型。只有多个服务确实共享同一语义和行为时，才提取跨服务抽象。

## 测试原则

- 单元测试不得依赖真实 RSS 服务。
- HTTP provider 测试必须验证 method、URL、请求体和响应解析。
- 类型转换测试必须验证 `pkg/types` 与 provider wire types 的请求、响应双向映射，并确保 wire types 不出现在 Handler API。
- 配置错误、未知 YAML 字段、非 2xx 响应和非法 JSON 应覆盖测试。
- 凭据测试必须验证加密文件不可直接读取，并验证保存、读取和删除流程。
