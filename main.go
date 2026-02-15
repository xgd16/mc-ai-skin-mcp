package main

import (
	"flag"
	"log"

	"mc-ai-skin/internal/mcp"
)

func main() {
	// gen.NewGenSkin(types.DefaultGenSkinOptions()).Gen(map[string][64][4]uint8{
	//  "head_top": gen.ColorBlock,
	// }, 64)
	useHTTP := flag.Bool("http", false, "启用 HTTP 模式，默认使用 stdio 模式")
	flag.Parse()

	if err := mcp.Run(*useHTTP); err != nil {
		log.Fatalf("MCP 服务器错误: %v", err)
	}
}
