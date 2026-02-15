package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mc-ai-skin/internal/gen"
	"mc-ai-skin/internal/types"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

const (
	serverName    = "mc-ai-skin"
	serverVersion = "1.0.0"
	defaultAddr   = ":19283"
)

// Run 启动 MCP 服务器
// useHTTP: true 使用 HTTP 模式，false 使用 stdio 模式
func Run(useHTTP bool) error {
	mcpServer := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// 注册 generate_minecraft_skin 工具
	tool := mcp.NewTool("generate_minecraft_skin",
		mcp.WithDescription("根据文字描述生成《我的世界》角色皮肤。支持 64x64。输出为 PNG 格式，遵循 Minecraft 皮肤 UV 布局（头部、躯干、手臂、腿部展开图）。"),
		mcp.WithString("schema", mcp.Required(), mcp.Description("皮肤部位名到颜色块数组的映射，每个数组对应一个皮肤部位的 8×8 像素 RGBA 色块，key 为部位名（见 skinMap）")),
	)

	mcpServer.AddTool(tool, generateMinecraftSkinHandler)

	if useHTTP {
		return runHTTP(mcpServer)
	}
	return runStdio(mcpServer)
}

func runStdio(mcpServer *server.MCPServer) error {
	log.Printf("[%s v%s] MCP 服务器已启动 (stdio 模式)", serverName, serverVersion)
	defer log.Printf("[%s] MCP 服务器已退出", serverName)
	return server.ServeStdio(mcpServer)
}

func runHTTP(mcpServer *server.MCPServer) error {
	log.Printf("[%s v%s] MCP 服务器已启动 (HTTP 模式): http://localhost%s/mcp", serverName, serverVersion, defaultAddr)
	defer log.Printf("[%s] MCP 服务器已退出", serverName)
	httpServer := server.NewStreamableHTTPServer(mcpServer)
	return httpServer.Start(defaultAddr)
}

// generateMinecraftSkinHandler 处理 generate_minecraft_skin 工具调用
func generateMinecraftSkinHandler(ctx context.Context, request mcp.CallToolRequest) (out *mcp.CallToolResult, err error) {

	schema := request.GetString("schema", "")
	if schema == "" {
		err = fmt.Errorf("schema is required")
		return
	}

	var colorMap map[string][64][4]uint8

	if err = json.Unmarshal([]byte(schema), &colorMap); err != nil {
		return
	}

	path, err := gen.NewGenSkin(types.DefaultGenSkinOptions()).Gen(colorMap, 64)
	if err != nil {
		return
	}

	out = mcp.NewToolResultText(path)
	return
}
