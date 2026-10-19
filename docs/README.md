# 文档导航

本文档面向 `rss-cli` 的使用者、维护者和协议集成开发者。

## Concepts

- [项目目标与边界](concepts/project-goals.md)：定位、核心概念和当前范围。

## Operation

- [使用与配置](operation/usage.md)：构建、context、FreshRSS 配置和命令输出。
- [Release 源码包验证](operation/release_verification.md)：使用 GPG 公钥和 detached signature 验证源码包。

## Developer

- [架构设计](developer/architecture/README.md)：代码职责、依赖关系和 `pkg` 一级包边界。
- [代码目录](developer/code-structure.md)：目录结构与放置规则。
- [开发指南](developer/development.md)：测试、构建和新增服务能力。
- [版本与构建规则](developer/versioning.md)：版本元数据和构建产物。
- [Git 协作规范](developer/git-contribution.md)：贡献和验证要求。

## Reference

- [FreshRSS Fever API](reference/fever-api.md)：认证、请求和响应约定。

文档中的“当前”只描述仓库中已经实现并可验证的行为；规划内容不会被写成稳定命令或公开 API。
