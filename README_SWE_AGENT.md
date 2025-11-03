# SWE-agent 配置指南

本项目已配置 [SWE-agent](https://github.com/SWE-agent/SWE-agent)，一个由 Princeton 和 Stanford 开源的 AI 代码修复工具。当你在 GitHub Issue 上打上 `bot:fix` 标签时，SWE-agent 会自动尝试修复问题并创建 PR。

## 🚀 快速开始

### 1. 配置 GitHub Secrets

在你的 GitHub 仓库设置中添加以下 Secrets：

1. 进入仓库：`Settings` → `Secrets and variables` → `Actions`
2. 添加以下 Secrets：

   **必需：**
   - `ANTHROPIC_API_KEY`: 你的 Anthropic API Key（如果使用 Claude）
   - 或 `OPENAI_API_KEY`: 你的 OpenAI API Key（如果使用 GPT）
   - 或 `TOGETHER_API_KEY`: 你的 Together API Key（如果使用 Together 模型）
   
   **可选（如果 SWE-agent 需要创建 PR）：**
   - `SWE_AGENT_GH_TOKEN`: GitHub Personal Access Token（需要 `repo` 权限）
     - 如果不设置，会使用默认的 `GITHUB_TOKEN`（权限可能受限）

### 2. 配置模型

编辑 `.github/swe-agent.yaml` 文件，设置你想要使用的模型：

```yaml
agent:
  model:
    name: "claude-sonnet-4-20250514"  # 或 "gpt-4o", "gpt-4-turbo" 等
    provider: "anthropic"  # 或 "openai", "together" 等
```

### 3. 使用方式

1. 创建一个 GitHub Issue 描述需要修复的问题
2. 给 Issue 打上 `bot:fix` 标签
3. GitHub Actions 会自动触发 SWE-agent
4. SWE-agent 会：
   - 在隔离的 Docker 环境中运行
   - 分析问题并修改代码
   - 创建新分支并推送更改
   - 自动创建 Pull Request

## 📋 配置说明

### 模型选择

SWE-agent 支持多种模型提供商：

- **Anthropic**: Claude Sonnet 4, Claude Opus 等
- **OpenAI**: GPT-4o, GPT-4 Turbo 等
- **Together**: 各种开源模型
- **本地模型**: 通过自定义 `api_base` 配置

### Docker 环境

默认使用 Docker 作为执行后端，提供隔离环境。如果不想使用 Docker（需要 Docker 安装在 runner 上），可以在 `.github/swe-agent.yaml` 中修改：

```yaml
environment:
  execution_type: "local"  # 改为 local
```

### 成本控制

在 `.github/swe-agent.yaml` 中可以配置：
- 最大迭代次数（`task.max_iterations`）
- 成本追踪开关（`agent.cost_tracking.enabled`）

## 🔧 高级配置

### 自定义 API 端点

如果你使用本地模型或代理，可以配置自定义 API 端点：

```yaml
agent:
  model:
    name: "custom-model"
    api_base: "https://your-api-endpoint.com/v1"
```

### 工具配置

可以在 `.github/swe-agent.yaml` 中启用/禁用特定工具，或调整工具参数。

## 📚 更多资源

- [SWE-agent 官方文档](https://swe-agent.com/docs)
- [SWE-agent GitHub](https://github.com/SWE-agent/SWE-agent)
- [配置参考](https://swe-agent.com/docs/configuration)

## ⚠️ 注意事项

1. **Docker 必需性**: 默认配置使用 Docker，GitHub Actions 的 Ubuntu runner 已包含 Docker。如果要本地运行，确保 Docker 已安装。

2. **权限要求**: 确保 GitHub Token 有足够权限（`contents: write`, `pull-requests: write`, `issues: write`）

3. **成本**: AI 模型调用会产生费用，注意监控使用情况

4. **安全性**: SWE-agent 会在隔离环境中执行代码，但仍需审查自动生成的 PR

## 🐛 故障排除

### Action 未触发
- 检查 Issue 标签是否为 `bot:fix`（区分大小写）
- 查看 Actions 日志

### API 调用失败
- 检查 API Key 是否正确配置
- 检查账户余额和配额

### Docker 相关问题
- GitHub Actions 的 Ubuntu runner 自带 Docker，无需额外配置
- 本地运行时确保 Docker 已安装并运行

