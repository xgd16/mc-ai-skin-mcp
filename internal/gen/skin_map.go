package gen

// AllParts 全部皮肤部位名列表（用于进度追踪与 Agent 分次调用）。
var AllParts = []string{
	"head_front", "head_back", "head_top", "head_bottom", "head_right", "head_left",
	"body_front", "body_back", "body_right", "body_left",
	"right_arm_front", "right_arm_back", "right_arm_right", "right_arm_left", "right_arm_top", "right_arm_bottom",
	"left_arm_front", "left_arm_back", "left_arm_right", "left_arm_left", "left_arm_top", "left_arm_bottom",
	"right_leg_front", "right_leg_back", "right_leg_right", "right_leg_left", "right_leg_top", "right_leg_bottom",
	"left_leg_front", "left_leg_back", "left_leg_right", "left_leg_left", "left_leg_top", "left_leg_bottom",
}

// skinMap 存储 Minecraft 皮肤各部位在皮肤贴图上的 UV 起始坐标 (x, y)。
// 坐标以 64×64 贴图左上角为原点，遵循 minotar/skin-spec (https://github.com/minotar/skin-spec)。
// 各部位实际像素范围：头部 8×8，躯干侧面 4×12 或 8×12，肢体 4×12，顶/底面 4×4。
var skinMap = map[string][2]int{
	// 头部 (8×8 each)
	"head_front":  {8, 8},   // Head Front (8,8,16,16)
	"head_back":   {24, 8},  // Head Back (24,8,32,16)
	"head_top":    {8, 0},   // Head Top (8,0,16,8)
	"head_bottom": {16, 0},  // Head Bottom (16,0,24,8)
	"head_right":  {0, 8},   // Head Right (0,8,8,16)
	"head_left":   {16, 8},  // Head Left (16,8,24,16)
	// 身体/Torso (侧面 4×12 或 8×12)
	"body_front": {20, 20}, // Torso Front (20,20,28,32)
	"body_back":  {32, 20}, // Torso Back (32,20,40,32)
	"body_right": {16, 20}, // Torso Right (16,20,20,32)
	"body_left":  {28, 20}, // Torso Left (28,20,32,32)
	// 右臂 Right Arm (4×12, 顶/底 4×4)
	"right_arm_front":   {44, 20}, // (44,20,48,32)
	"right_arm_back":   {52, 20},  // (52,20,56,32)
	"right_arm_right":  {40, 20},  // (40,20,44,32)
	"right_arm_left":   {48, 20},  // (48,20,52,32)
	"right_arm_top":    {44, 16},  // (44,16,48,20)
	"right_arm_bottom": {48, 16},  // (48,16,52,20)
	// 左臂 Left Arm (1.8+, 4×12)
	"left_arm_front":   {36, 52}, // (36,52,40,64)
	"left_arm_back":    {44, 52},  // (44,52,48,64)
	"left_arm_right":   {32, 52}, // (32,52,36,64)
	"left_arm_left":    {40, 52},  // (40,52,44,64)
	"left_arm_top":    {36, 48},   // (36,48,40,52)
	"left_arm_bottom":  {40, 48},  // (40,48,44,52)
	// 右腿 Right Leg (4×12)
	"right_leg_front":   {4, 20},  // (4,20,8,32)
	"right_leg_back":   {12, 20},   // (12,20,16,32)
	"right_leg_right":  {0, 20},    // (0,20,4,32)
	"right_leg_left":   {8, 20},    // (8,20,12,32)
	"right_leg_top":    {4, 16},    // (4,16,8,20)
	"right_leg_bottom": {8, 16},    // (8,16,12,20)
	// 左腿 Left Leg (1.8+, 4×12)
	"left_leg_front":   {20, 52}, // (20,52,24,64)
	"left_leg_back":   {28, 52},   // (28,52,32,64)
	"left_leg_right":  {16, 52},   // (16,52,20,64)
	"left_leg_left":   {24, 52},   // (24,52,28,64)
	"left_leg_top":    {20, 48},    // (20,48,24,52)
	"left_leg_bottom": {24, 48},   // (24,48,28,52)
}

// skinSizeMap 各部位在贴图上的实际像素尺寸 (宽, 高)，与 skinMap 一一对应。
// 头部 8×8；躯干正面/背面 8×12、侧面 4×12；四肢侧面 4×12、顶/底面 4×4。
var skinSizeMap = map[string][2]int{
	"head_front":  {8, 8}, "head_back": {8, 8}, "head_top": {8, 8}, "head_bottom": {8, 8},
	"head_right":  {8, 8}, "head_left": {8, 8},
	"body_front":  {8, 12}, "body_back": {8, 12}, "body_right": {4, 12}, "body_left": {4, 12},
	"right_arm_front": {4, 12}, "right_arm_back": {4, 12}, "right_arm_right": {4, 12}, "right_arm_left": {4, 12},
	"right_arm_top": {4, 4}, "right_arm_bottom": {4, 4},
	"left_arm_front": {4, 12}, "left_arm_back": {4, 12}, "left_arm_right": {4, 12}, "left_arm_left": {4, 12},
	"left_arm_top": {4, 4}, "left_arm_bottom": {4, 4},
	"right_leg_front": {4, 12}, "right_leg_back": {4, 12}, "right_leg_right": {4, 12}, "right_leg_left": {4, 12},
	"right_leg_top": {4, 4}, "right_leg_bottom": {4, 4},
	"left_leg_front": {4, 12}, "left_leg_back": {4, 12}, "left_leg_right": {4, 12}, "left_leg_left": {4, 12},
	"left_leg_top": {4, 4}, "left_leg_bottom": {4, 4},
}