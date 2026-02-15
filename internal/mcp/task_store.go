package mcp

import (
	"sync"

	"github.com/google/uuid"
)

// TaskState 任务状态：用于分次绘制时追踪进度
type TaskState struct {
	TaskID        string   // 任务随机码
	OutputPath    string   // 当前皮肤图片保存路径
	CompletedParts []string // 已绘制的部位列表
}

var (
	taskStore = make(map[string]*TaskState)
	taskMu    sync.RWMutex
)

// NewTaskID 生成新的任务随机码
func NewTaskID() string {
	return uuid.New().String()
}

// SaveTask 保存或更新任务状态
func SaveTask(t *TaskState) {
	taskMu.Lock()
	defer taskMu.Unlock()
	taskStore[t.TaskID] = t
}

// GetTask 根据任务码获取任务状态
func GetTask(taskID string) (*TaskState, bool) {
	taskMu.RLock()
	defer taskMu.RUnlock()
	t, ok := taskStore[taskID]
	if !ok || t == nil {
		return nil, false
	}
	return t, true
}
