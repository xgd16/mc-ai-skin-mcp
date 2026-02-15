package types

// GenSkinOptions 用于生成皮肤的配置信息
type GenSkinOptions struct {
	DateFileName bool   // 是否使用日期作为文件名
	PrintOutPath bool   // 是否打印输出路径
	BaseName     string // 生成文件的基础名称
}

func DefaultGenSkinOptions() *GenSkinOptions {
	return &GenSkinOptions{
		DateFileName: false,
		PrintOutPath: true,
		BaseName:     "skin",
	}
}
