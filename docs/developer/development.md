# 开发指南

## 本地验证

```sh
mise install
gofmt -w ./cmd ./internal ./pkg ./testsuite
go test ./...
make build
```

构建产物写入 `build/<VERSION>/`，该目录已加入 Git 忽略规则。`make clean` 会删除整个 `build/`。

跨平台构建：

- `make platform-builds` 构建预设的全部 `GOOS/GOARCH` 组合。
- `make build-matrix GOOS=linux GOARCH=amd64` 构建单个指定组合。
- 交叉构建固定使用 `CGO_ENABLED=0` 和 `-trimpath`，产物名包含目标平台和架构。

UPX 压缩：

- `UPX_ENABLED` 默认为空；设置为 `1` 后，构建完成自动压缩二进制。
- `UPX_LEVEL` 默认是 `9`，也可通过 `UPX_ARGS` 直接覆盖 UPX 参数。
- `UPX_SUPPORTED_PLATFORMS` 限制跨平台压缩范围，不支持的平台只跳过压缩，不影响构建。
- 项目通过 `mise.toml` 固定 `upx` 版本；首次使用前运行 `mise install`。
- 启用压缩前由 `check-upx` 检查 `upx` 是否存在；缺少工具时直接失败，避免产出未按要求处理的构建结果。

版本元数据：

- `pkg/version` 保存 `VERSION`、`BUILDTIME` 和 `GITHASH` 三项构建信息。
- `Makefile` 使用 `-ldflags -X` 注入发布值；默认版本来自 Git tag、日期和当前 commit。
- 不使用 Makefile 时允许直接 `go build`，但得到的是开发默认值。

静态检查：

```sh
make lint
```

集成测试需要 Docker：

```sh
make integration-test
```

它是显式 opt-in 测试，会启动隔离的 FreshRSS Compose 实例并在结束时清理容器和卷。

## 新增服务能力

1. 在 `pkg/types` 确认或新增稳定的应用请求/响应类型；不要把 provider JSON tag 和协议 envelope 泄漏到公共类型。
2. 在 `pkg/thirdparty/<service>` 实现 provider；client 使用 `pkg/rest/client`，Handler 转换 `pkg/types`，请求和响应 wire types 只在 provider 内部使用，并补齐 transport stub 测试。
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
