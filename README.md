<div align="center">
  <img src="https://avatars.githubusercontent.com/u/245985800?s=200&v=4" alt="logo" width="100" />
  <h2>Example</h2>
  <h3>FasterEdge 主项目完整使用示例（Go Demo）</h3>
</div>

### 一、项目简介

- **用途**: 基于 [FasterEdge](https://github.com/FasterEdge/FasterEdge) 主项目的**完整可运行示例**，展示如何初始化 Atom、挂载内置能力、调用各能力命令、执行鉴权命令链路以及管理组件生命周期。
- **依赖版本**: `github.com/FasterEdge/FasterEdge v1.0.20261022`（见 `go.mod`）

### 二、快速开始

```bash
# 1. 拉取依赖（含主项目 FasterEdge）
go mod download

# 2. 直接运行 demo
go run .

# 3. 编译为本地可执行文件
go build .
```

### 三、Demo 演示内容

| 能力 | 命令 | 说明 |
|---|---|---|
| **TimeAbility** | `get_time` | 获取当前时间快照（来源 / 时间） |
| **NetMapAbility** | `register_peer` / `list_peers` | 注册对等节点并查询网络拓扑 |
| **OneKeyAbility** | `issue_token` / `verify_token` | 为对等节点签发短期访问令牌并校验 |
| **CmdAbility** | `set_allowlist` / `run` | 配置命令执行白名单并运行受控命令 |
| **鉴权链路** | `AuthenticatedCommand` | 用 OneKey 凭证走安全边界调用命令；篡改签名会被拒绝 |
| **生命周期** | `RunAtom` / `CloseAtom` | 运行全部 Runner 能力，取消后优雅关闭 |

### 四、跨平台自动打包

仓库配置了 GitHub Actions（`.github/workflows/release.yml`）：

- 触发方式：推送 `v*` tag 或手动 `workflow_dispatch`
- 构建矩阵：**Windows / macOS / Linux × amd64 / arm64**，共 6 个平台
- 发布方式：自动打包 `tar.gz`（Unix）与 `zip`（Windows）并发布到 GitHub Release

```bash
# 打 tag 触发自动打包发布
git tag v1.0.0 && git push origin v1.0.0
```

### 五、许可证

- **许可协议**: 仓库根目录 `LICENSE` 为 Apache License 2.0。
- **适用说明**: 示例源码保留 FasterEdge 组织归属与版权声明。
