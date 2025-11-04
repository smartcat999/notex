# Claude Agent Skills

本目录包含用于限制 Claude Agent 行为的 Skills。

## Skills 列表

- **code-focus**: 限制 agent 专注于代码修复，避免过度解读和生成不必要的文档
- **minimal-changes**: 确保代码修改最小化，只做必要的改动

## 使用方法

这些 Skills 会自动被 Claude Agent SDK 加载，限制 agent 的行为。

## 配置说明

SWE-agent 的配置在 `.github/swe-agent.yaml` 中，包含：
- 迭代次数限制
- 成本限制
- 工具启用/禁用
- PR 描述长度限制

