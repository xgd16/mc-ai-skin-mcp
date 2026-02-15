# mc-ai-skin

基于 MCP (Model Context Protocol) 的 **Minecraft 皮肤生成服务**。AI Agent 可通过 `generate_minecraft_skin` 工具，根据文字描述生成符合 Minecraft UV 规范的 64×64 皮肤 PNG 贴图。

## 特性

- **MCP 工具**：作为 MCP 服务器运行，可被 Cursor、Claude Desktop 等支持 MCP 的 Agent 调用
- **分部位分次绘制**：每次仅绘制 1 个部位，支持 `task_id` 续传，避免长时间连接超时
- **标准 UV 布局**：严格遵循 Minecraft 64×64 皮肤贴图布局，覆盖 34 个身体部位
- **任务持久化**：支持任务续传与进度查询

## 快速开始

### 环境要求

- Go 1.25+
- （可选）支持 MCP 的 AI 客户端（如 Cursor、Claude Desktop）

### 构建与运行

```bash
# 构建
make build

# 运行（默认 stdio 模式，供 MCP 客户端调用）
./mc-ai-skin

# HTTP 模式（端口 19283，供远程调用）
./mc-ai-skin -http
```

### MCP 客户端配置

将 `mc-ai-skin` 添加到 MCP 配置，例如 Cursor 的 `~/.cursor/mcp.json`：

```json
{
  "mcpServers": {
    "mc-ai-skin": {
      "command": "/path/to/mc-ai-skin",
      "args": []
    }
  }
}
```

HTTP 模式示例：

```json
{
  "mcpServers": {
    "mc-ai-skin": {
      "url": "http://localhost:19283/mcp"
    }
  }
}
```

## 工具说明

### `generate_minecraft_skin`

根据文字描述生成 Minecraft 角色皮肤。

| 参数 | 必填 | 说明 |
|------|------|------|
| `schema` | 首次/续传时 | 皮肤部位到 RGBA 像素数组的 JSON 映射，每次仅可包含 **1 个** 部位 |
| `task_id` | 续传/查询时 | 任务标识，从上次返回中获取 |

**schema 格式**：`map[string][64][4]uint8`，每个部位 64 个 RGBA 像素（8×8），例如：

```json
{
  "head_front": [[r,g,b,a], [r,g,b,a], ...]
}
```

**返回示例**：

```json
{
  "task_id": "550e8400-...",
  "path": "/path/to/output/skin_550e8400.png",
  "completed_parts": ["head_front", "head_back"],
  "pending_parts": ["head_top", "head_bottom", ...],
  "message": "已绘制 2 个部位，剩余 32 个待绘制"
}
```

- **新建任务**：传 `schema`，不传 `task_id`
- **续传绘制**：传 `task_id` + `schema`
- **查询进度**：仅传 `task_id`

## 支持的皮肤部位

| 部位类型 | 部位 key |
|----------|----------|
| 头部 | `head_front`, `head_back`, `head_top`, `head_bottom`, `head_right`, `head_left` |
| 躯干 | `body_front`, `body_back`, `body_right`, `body_left` |
| 右臂 | `right_arm_front`, `right_arm_back`, `right_arm_right`, `right_arm_left`, `right_arm_top`, `right_arm_bottom` |
| 左臂 | `left_arm_front`, `left_arm_back`, `left_arm_right`, `left_arm_left`, `left_arm_top`, `left_arm_bottom` |
| 右腿 | `right_leg_front`, `right_leg_back`, `right_leg_right`, `right_leg_left`, `right_leg_top`, `right_leg_bottom` |
| 左腿 | `left_leg_front`, `left_leg_back`, `left_leg_right`, `left_leg_left`, `left_leg_top`, `left_leg_bottom` |

共 34 个部位，每个部位 8×8 像素。

## 输出

生成的皮肤保存到 `output/` 目录，文件名为 `skin_<task_id前8位>.png`。

## Agent 提示词

使用 AI Agent 调用时，可参考 [`docs/ai_agent_prompt.md`](docs/ai_agent_prompt.md) 中的画师准则、像素技巧和完整 Prompt 模板，以获得细节更丰富的皮肤效果。

## 开发

```bash
# 开发构建（含调试信息）
make build-dev

# 运行测试
make test

# 清理
make clean
```

## 技术栈

- [mcp-go](https://github.com/mark3labs/mcp-go) - MCP 协议实现
- [gogf/gf](https://github.com/gogf/gf) - 工具库
- [disintegration/imaging](https://github.com/disintegration/imaging) - 图像处理

## License

MIT
