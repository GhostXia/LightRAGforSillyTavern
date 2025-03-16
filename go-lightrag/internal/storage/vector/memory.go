package vector

import (
	"context"
	"sync"

	"github.com/HerSophia/go-lightrag/pkg/interfaces"
)

// MemoryVectorStorage 内存向量存储实现
type MemoryVectorStorage struct {
	mu      sync.RWMutex
	vectors map[string]vectorEntry
}

type vectorEntry struct {
	Vector   []float32
	Metadata map[string]interface{}
}

// NewMemoryVectorStorage 创建新的内存向量存储
func NewMemoryVectorStorage() *MemoryVectorStorage {
	return &MemoryVectorStorage{
		vectors: make(map[string]vectorEntry),
	}
}

// Upsert 插入或更新向量
func (s *MemoryVectorStorage) Upsert(ctx context.Context, id string, vector []float32, metadata map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.vectors[id] = vectorEntry{
		Vector:   vector,
		Metadata: metadata,
	}

	return nil
}

// Query 查询最相似的向量
func (s *MemoryVectorStorage) Query(ctx context.Context, vector []float32, topK int) ([]interfaces.SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 计算所有向量的相似度
	scores := make([]interfaces.SearchResult, 0, len(s.vectors))
	for id, entry := range s.vectors {
		score := cosineSimilarity(vector, entry.Vector)
		scores = append(scores, interfaces.SearchResult{
			ID:       id,
			Score:    score,
			Metadata: entry.Metadata,
		})
	}

	// 按相似度排序
	sortSearchResults(scores)

	// 返回前topK个结果
	if len(scores) > topK {
		scores = scores[:topK]
	}

	return scores, nil
}

// Delete 删除向量
func (s *MemoryVectorStorage) Delete(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.vectors[id]; !exists {
		return interfaces.ErrNotFound
	}

	delete(s.vectors, id)
	return nil
}

// Close 关闭存储
func (s *MemoryVectorStorage) Close() error {
	// 内存实现不需要特殊关闭操作
	return nil
}

// 辅助函数：计算余弦相似度
func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) {
		return 0
	}

	var dotProduct float32
	var normA float32
	var normB float32

	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}

	if normA == 0 || normB == 0 {
		return 0
	}

	return dotProduct / (float32(sqrt(float64(normA))) * float32(sqrt(float64(normB))))
}

// 辅助函数：简单平方根计算
func sqrt(x float64) float64 {
	// 简单实现，实际应用中应使用math.Sqrt
	z := 1.0
	for i := 0; i < 10; i++ { // 牛顿迭代法
		z -= (z*z - x) / (2 * z)
	}
	return z
}

// 辅助函数：排序搜索结果
func sortSearchResults(results []interfaces.SearchResult) {
	// 简单的冒泡排序，实际应用中可使用更高效的排序算法
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].Score < results[j].Score {
				results[i], results[j] = results[j], results[i]
			}
		}
	}
}
