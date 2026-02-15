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
var skinMap = map[string][2]int{
	// 头部
	"head_front":  {8, 8},   // 头部正面
	"head_back":   {24, 8},  // 头部背面
	"head_top":    {8, 0},   // 头部顶面
	"head_bottom": {16, 0},  // 头部底面
	"head_right":  {0, 8},   // 头部右侧（视角正面看玩家的右侧）
	"head_left":   {16, 8},  // 头部左侧（视角正面看玩家的左侧）
	// 身体
	"body_front": {20, 20}, // 身体正面
	"body_back":  {32, 20}, // 身体背面
	"body_right": {16, 20}, // 身体右侧（视角正面看玩家的右侧）
	"body_left":  {28, 20}, // 身体左侧（视角正面看玩家的左侧）
	// 右臂
	"right_arm_front":   {44, 20}, // 右臂正面
	"right_arm_back":   {52, 20}, // 右臂背面
	"right_arm_right":  {40, 20}, // 右臂外侧
	"right_arm_left":   {48, 20}, // 右臂内侧
	"right_arm_top":    {44, 16}, // 右臂顶面
	"right_arm_bottom": {48, 16}, // 右臂底面
	// 左臂
	"left_arm_front":   {36, 52}, // 左臂正面
	"left_arm_back":    {44, 52}, // 左臂背面
	"left_arm_right":   {32, 52}, // 左臂外侧
	"left_arm_left":    {40, 52}, // 左臂内侧
	"left_arm_top":    {36, 48}, // 左臂顶面
	"left_arm_bottom":  {40, 48}, // 左臂底面
	// 右腿
	"right_leg_front":   {4, 20},  // 右腿正面
	"right_leg_back":   {12, 20},  // 右腿背面
	"right_leg_right":  {0, 20},   // 右腿外侧
	"right_leg_left":   {8, 20},   // 右腿内侧
	"right_leg_top":    {4, 16},   // 右腿顶面
	"right_leg_bottom": {8, 16},   // 右腿底面
	// 左腿
	"left_leg_front":   {20, 52}, // 左腿正面
	"left_leg_back":   {28, 52},  // 左腿背面
	"left_leg_right":  {16, 52},  // 左腿外侧
	"left_leg_left":   {24, 52},  // 左腿内侧
	"left_leg_top":    {20, 48},  // 左腿顶面
	"left_leg_bottom": {24, 48},  // 左腿底面
}