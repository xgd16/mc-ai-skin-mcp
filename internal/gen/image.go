package gen

import (
	"fmt"
	"image"
	"image/color"
	"mc-ai-skin/internal/types"

	"github.com/disintegration/imaging"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

type GenSkin struct {
	options *types.GenSkinOptions
}

func NewGenSkin(options *types.GenSkinOptions) *GenSkin {
	return &GenSkin{
		options: options,
	}
}

// Gen 生成皮肤贴图
//
// 参数：
//   ColorMap map[string][64][4]uint8 —— 皮肤部位名到颜色块数组的映射，每个数组对应一个皮肤部位的 8×8 像素 RGBA 色块，
//                                     key 为部位名（见 skinMap）
//   size int —— 输出图片的尺寸宽高（像素），例如 64 或 128，即生成 size×size 的皮肤贴图
//
// 返回值：
//   out string —— 生成图片的实际保存路径
//   err error  —— 生成或保存图片失败时返回错误，成功为 nil
func (s *GenSkin) Gen(ColorMap map[string][64][4]uint8, size int) (out string, err error) {
	path, err := OutputDir()
	if err != nil {
		return "", err
	}
	outPath := path + "/" + s.createFileName()
	return s.GenToPath(ColorMap, size, outPath)
}

// GenToPath 将皮肤贴图生成并保存到指定路径，用于任务化分次调用。
func (s *GenSkin) GenToPath(ColorMap map[string][64][4]uint8, size int, outPath string) (out string, err error) {
	img := imaging.New(size, size, color.RGBA{0, 0, 0, 0})
	for k, v := range ColorMap {
		x, y, e := s.getSkinXYMap(k)
		if e != nil {
			return "", e
		}
		s.printColorBlock(x, y, v, img)
	}
	if s.options.PrintOutPath {
		fmt.Printf("save to: %s\n", outPath)
	}
	return outPath, imaging.Save(img, outPath)
}

// GenMerge 在已有皮肤图像上合并绘制新部位，用于分次调用、避免长时间断连。
// existingPath 为已存在的皮肤图片路径，colorMap 为本次要绘制的新部位。
// 返回更新后的保存路径（与 existingPath 相同）。
func (s *GenSkin) GenMerge(existingPath string, colorMap map[string][64][4]uint8, size int) (out string, err error) {
	if len(colorMap) == 0 {
		return existingPath, nil
	}
	img, err := imaging.Open(existingPath)
	if err != nil {
		return "", err
	}
	nrgba := imaging.Clone(img)
	for k, v := range colorMap {
		x, y, e := s.getSkinXYMap(k)
		if e != nil {
			return "", e
		}
		s.printColorBlock(x, y, v, nrgba)
	}
	if s.options.PrintOutPath {
		fmt.Printf("merge save to: %s\n", existingPath)
	}
	err = imaging.Save(nrgba, existingPath)
	return existingPath, err
}

// getSkinXYMap 用于查询指定皮肤部位的 UV 起始坐标（x, y）。
// 参数：
//
//	k string —— skinMap 的 key，代表皮肤的各个部位（如 "head_front", "body_right" 等）
//
// 返回值：
//
//	x int  —— 该部位在皮肤贴图上起始的横向（x）坐标，单位为像素
//	y int  —— 该部位在皮肤贴图上起始的纵向（y）坐标，单位为像素
//	err error —— 若未找到 key，则返回错误；否则 err 为 nil
//
// 使用说明：
//
//	若 key 在 skinMap 中存在，则返回对应的 (x, y)，否则返回错误。
//	例如，getSkinXYMap("head_front") 会返回 8, 8 并 err 为 nil，
//	若查询不存在的 key，则 err 不为 nil，x、y 均为 0。
func (s *GenSkin) getSkinXYMap(k string) (x, y int, err error) {
	// 从 skinMap 查询对应的坐标值
	v, ok := skinMap[k]
	if !ok {
		// 若未找到该 key，则返回错误提示
		return 0, 0, fmt.Errorf("skin %s not found in skinMap", k)
	}
	// 正常情况，解构返回 UV 坐标
	return v[0], v[1], nil
}

// printColorBlock 在指定坐标 (x, y) 处，按 8 列方式依次绘制颜色块数组 arr 到 img 图像上。
// x, y: 绘制区域的起始坐标，分别表示第一个块的横向与纵向起点（像素单位）
// arr: 颜色块数组，每个元素为 RGBA 四元组，长度为 64，代表 64 个颜色块
// img: 目标图像指针，类型为 *image.NRGBA，用于绘制颜色块
func (s *GenSkin) printColorBlock(x, y int, arr [64][4]uint8, img *image.NRGBA) {
	const cols = 8 // 每行显示的颜色块个数
	/*
		循环遍历颜色块数组 arr，将每个颜色块按顺序绘制到目标图像 img 上。
		排列方式为每行 8 个块，超出自动换行。
		计算方法：
			col = k % cols 获取每个颜色块在当前行的列索引（0~7）
			row = k / cols 获取每个颜色块所在的行索引
			最终像素坐标为 (x + col, y + row)
		颜色使用 color.RGBA 设置，参数分别对应 R、G、B、A 通道
	*/
	for k, item := range arr {
		col := k % cols // 当前颜色块在行内的列号
		row := k / cols // 当前颜色块所在的行号
		img.Set(
			x+col, // 水平坐标：起始 x + 列号
			y+row, // 垂直坐标：起始 y + 行号
			color.RGBA{
				item[0], // 红色通道
				item[1], // 绿色通道
				item[2], // 蓝色通道
				item[3], // alpha 通道（透明度）
			},
		)
	}
}

// OutputDir 返回 output 目录绝对路径并确保存在，供 MCP 任务化调用使用。
func OutputDir() (string, error) {
	path := gfile.Pwd() + "/output"
	if !gfile.Exists(path) {
		if err := gfile.Mkdir(path); err != nil {
			return "", err
		}
	}
	return path, nil
}


// 创建文件名
func (s *GenSkin) createFileName() string {
	if s.options.DateFileName {
		return fmt.Sprintf("%s_%s.png", s.options.BaseName, gtime.Now().Format("Y_m_d_H_i_s"))
	}
	return fmt.Sprintf("%s.png", s.options.BaseName)
}
