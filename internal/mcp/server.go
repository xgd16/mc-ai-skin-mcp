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

// toolResponse 返回给 Agent 的结构化数据
type toolResponse struct {
	TaskID         string   `json:"task_id"`          // 任务随机码，下次请求需携带
	Path           string   `json:"path"`             // 当前皮肤图片保存路径
	CompletedParts []string `json:"completed_parts"`  // 已绘制的部位
	PendingParts   []string `json:"pending_parts"`    // 尚未绘制的部位
	Message        string   `json:"message,omitempty"` // 简要说明
}

// Run 启动 MCP 服务器
// useHTTP: true 使用 HTTP 模式，false 使用 stdio 模式
func Run(useHTTP bool) error {
	mcpServer := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithToolCapabilities(true),
		server.WithRecovery(),
	)

	// 注册 generate_minecraft_skin 工具（支持分部位分次调用，避免长时间断连）
	tool := mcp.NewTool("generate_minecraft_skin",
		mcp.WithDescription("根据文字描述生成《我的世界》角色皮肤。支持 64x64 及分部位分次调用。首次调用传 schema，后续调用传 task_id 与 schema；仅传 task_id 可查询进度。返回 task_id、path、completed_parts、pending_parts。"),
		mcp.WithString("schema", mcp.Description("皮肤部位名到颜色块数组的映射，JSON 字符串，每个数组 64 个 RGBA，key 为部位名。可分批传入，后续请求需带 task_id")),
		mcp.WithString("task_id", mcp.Description("任务随机码。续传或查询进度时必填，从上次返回中获取")),
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
// 支持分部位分次调用：首次传 schema，后续传 task_id + schema；仅 task_id 可查询进度
func generateMinecraftSkinHandler(ctx context.Context, request mcp.CallToolRequest) (out *mcp.CallToolResult, err error) {
	taskID := request.GetString("task_id", "")
	schema := request.GetString("schema", "")

	opts := types.DefaultGenSkinOptions()
	g := gen.NewGenSkin(opts)

	// 解析 schema（若有）
	var colorMap map[string][64][4]uint8
	if schema != "" {
		if err = json.Unmarshal([]byte(schema), &colorMap); err != nil {
			return
		}
	}

	// 计算待绘制部位
	completedSet := make(map[string]bool)
	var pendingParts []string
	computePending := func(completed []string) {
		completedSet = make(map[string]bool)
		for _, p := range completed {
			completedSet[p] = true
		}
		pendingParts = nil
		for _, p := range gen.AllParts {
			if !completedSet[p] {
				pendingParts = append(pendingParts, p)
			}
		}
	}

	var path string
	var completedParts []string

	if taskID == "" {
		// 新任务：必须有 schema
		if schema == "" || len(colorMap) == 0 {
			err = fmt.Errorf("新建任务时 schema 必填且不能为空")
			return
		}
		taskID = NewTaskID()
		dir, e := gen.OutputDir()
		if e != nil {
			err = e
			return
		}
		path = dir + "/skin_" + taskID[:8] + ".png"
		if path, err = g.GenToPath(colorMap, 64, path); err != nil {
			return
		}
		for k := range colorMap {
			completedParts = append(completedParts, k)
		}
		SaveTask(&TaskState{TaskID: taskID, OutputPath: path, CompletedParts: completedParts})
		computePending(completedParts)
	} else {
		// 续传或查询进度
		t, ok := GetTask(taskID)
		if !ok {
			err = fmt.Errorf("任务码无效或已过期: %s", taskID)
			return
		}
		path = t.OutputPath
		completedParts = append([]string{}, t.CompletedParts...)
		for _, p := range completedParts {
			completedSet[p] = true
		}

		if len(colorMap) > 0 {
			if path, err = g.GenMerge(path, colorMap, 64); err != nil {
				return
			}
			for k := range colorMap {
				if !completedSet[k] {
					completedParts = append(completedParts, k)
					completedSet[k] = true
				}
			}
			t.OutputPath = path
			t.CompletedParts = completedParts
			SaveTask(t)
		}
		computePending(completedParts)
	}

	resp := toolResponse{
		TaskID:         taskID,
		Path:           path,
		CompletedParts: completedParts,
		PendingParts:   pendingParts,
		Message:        fmt.Sprintf("已绘制 %d 个部位，剩余 %d 个待绘制", len(completedParts), len(pendingParts)),
	}
	data, _ := json.Marshal(resp)
	out = mcp.NewToolResultText(string(data))
	return
}
