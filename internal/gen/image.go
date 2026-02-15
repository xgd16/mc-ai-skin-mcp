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
// size 为输出贴图边长（如 64 或 128），内部按 size/64 对 UV 坐标与尺寸做缩放。
func (s *GenSkin) GenToPath(ColorMap map[string][64][4]uint8, size int, outPath string) (out string, err error) {
	if size <= 0 {
		size = 64
	}
	scale := size / 64
	if scale <= 0 {
		scale = 1
	}
	img := imaging.New(size, size, color.RGBA{0, 0, 0, 0})
	for k, v := range ColorMap {
		x, y, e := s.getSkinXYMap(k)
		if e != nil {
			return "", e
		}
		w, h := s.getSkinSize(k)
		s.printColorBlock(x*scale, y*scale, w*scale, h*scale, v, img)
	}
	if s.options.PrintOutPath {
		fmt.Printf("save to: %s\n", outPath)
	}
	return outPath, imaging.Save(img, outPath)
}

// GenMerge 在已有皮肤图像上合并绘制新部位，用于分次调用、避免长时间断连。
// existingPath 为已存在的皮肤图片路径，colorMap 为本次要绘制的新部位。
// size 为贴图边长，与 GenToPath 一致时缩放正确；返回更新后的保存路径（与 existingPath 相同）。
func (s *GenSkin) GenMerge(existingPath string, colorMap map[string][64][4]uint8, size int) (out string, err error) {
	if len(colorMap) == 0 {
		return existingPath, nil
	}
	if size <= 0 {
		size = 64
	}
	scale := size / 64
	if scale <= 0 {
		scale = 1
	}
	img, err := imaging.Open(existingPath)
	if err != nil {
		return "", err
	}
	// imaging.Clone 返回 *image.NRGBA，可直接用于绘制
	nrgba := imaging.Clone(img)
	for k, v := range colorMap {
		x, y, e := s.getSkinXYMap(k)
		if e != nil {
			return "", e
		}
		w, h := s.getSkinSize(k)
		s.printColorBlock(x*scale, y*scale, w*scale, h*scale, v, nrgba)
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
	v, ok := skinMap[k]
	if !ok {
		return 0, 0, fmt.Errorf("skin %s not found in skinMap", k)
	}
	return v[0], v[1], nil
}

// getSkinSize 返回部位在贴图上的像素尺寸 (宽, 高)。若 key 不存在则返回 8,8 作为默认。
func (s *GenSkin) getSkinSize(k string) (w, h int) {
	v, ok := skinSizeMap[k]
	if !ok {
		return 8, 8
	}
	return v[0], v[1]
}

// printColorBlock 将 64 像素（8×8）的 arr 按尺寸 (w,h) 缩放绘制到 img 的 (x,y) 起始区域。
// 头部 8×8 原样映射；躯干/四肢按 skinSizeMap 的 8×12、4×12、4×4 等尺寸绘制，不越界。
// x, y: 绘制起点；w, h: 该部位在贴图上的宽高；arr: 64 个 RGBA（行优先 8×8）；img: 目标图。
func (s *GenSkin) printColorBlock(x, y, w, h int, arr [64][4]uint8, img *image.NRGBA) {
	if w <= 0 || h <= 0 {
		return
	}
	const srcSize = 8
	bounds := img.Bounds()
	maxX := bounds.Max.X
	maxY := bounds.Max.Y
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			px := x + dx
			py := y + dy
			if px < bounds.Min.X || px >= maxX || py < bounds.Min.Y || py >= maxY {
				continue
			}
			srcX := dx * srcSize / w
			if srcX >= srcSize {
				srcX = srcSize - 1
			}
			srcY := dy * srcSize / h
			if srcY >= srcSize {
				srcY = srcSize - 1
			}
			idx := srcY*srcSize + srcX
			item := arr[idx]
			img.Set(px, py, color.RGBA{item[0], item[1], item[2], item[3]})
		}
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
