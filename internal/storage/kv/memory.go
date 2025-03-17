package kv

import (
	"context"
	"sync"

	"github.com/HerSophia/go-lightrag/pkg/interfaces"
)

// MemoryKVStorage 内存键值存储实现
type MemoryKVStorage struct {
	mu    sync.RWMutex
	items map[string]interface{}
}

// NewMemoryKVStorage 创建新的内存键值存储
func NewMemoryKVStorage() *MemoryKVStorage {
	return &MemoryKVStorage{
		items: make(map[string]interface{}),
	}
}

// Get 获取值
func (s *MemoryKVStorage) Get(ctx context.Context, key string) (interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.items[key]
	if !exists {
		return nil, interfaces.ErrNotFound
	}

	return value, nil
}

// GetMany 批量获取值
func (s *MemoryKVStorage) GetMany(ctx context.Context, keys []string) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]interface{}, len(keys))
	for _, key := range keys {
		if value, exists := s.items[key]; exists {
			result[key] = value
		}
	}

	return result, nil
}

// Set 设置值
func (s *MemoryKVStorage) Set(ctx context.Context, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = value
	return nil
}

// Delete 删除键
func (s *MemoryKVStorage) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.items[key]; !exists {
		return interfaces.ErrNotFound
	}

	delete(s.items, key)
	return nil
}

// Keys 获取所有键
func (s *MemoryKVStorage) Keys(ctx context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, 0, len(s.items))
	for key := range s.items {
		keys = append(keys, key)
	}

	return keys, nil
}

// Close 关闭存储
func (s *MemoryKVStorage) Close() error {
	// 内存实现不需要特殊关闭操作
	return nil
}
