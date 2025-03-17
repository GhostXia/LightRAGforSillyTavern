package lightrag

import (
	"context"
	"errors"
)

// 常见错误定义
var (
	ErrNotFound  = errors.New("item not found")
	ErrInvalidID = errors.New("invalid ID")
)

// VectorStorage 向量存储接口
type VectorStorage interface {
	// Upsert 插入或更新向量
	Upsert(ctx context.Context, id string, vector []float32, metadata map[string]interface{}) error

	// Query 查询最相似的向量
	Query(ctx context.Context, vector []float32, topK int) ([]SearchResult, error)

	// Delete 删除向量
	Delete(ctx context.Context, id string) error

	// Close 关闭存储
	Close() error
}

// SearchResult 搜索结果
type SearchResult struct {
	ID       string
	Score    float32
	Metadata map[string]interface{}
}

// KVStorage 键值存储接口
type KVStorage interface {
	// Get 获取值
	Get(ctx context.Context, key string) (interface{}, error)

	// GetMany 批量获取值
	GetMany(ctx context.Context, keys []string) (map[string]interface{}, error)

	// Set 设置值
	Set(ctx context.Context, key string, value interface{}) error

	// Delete 删除键
	Delete(ctx context.Context, key string) error

	// Keys 获取所有键
	Keys(ctx context.Context) ([]string, error)

	// Close 关闭存储
	Close() error
}

// Node 图节点
type Node struct {
	ID         string
	Type       string
	Properties map[string]interface{}
}

// Edge 图边
type Edge struct {
	ID          string
	Type        string
	SourceID    string
	TargetID    string
	Properties  map[string]interface{}
	Directional bool
}

// GraphStorage 图存储接口
type GraphStorage interface {
	// AddNode 添加节点
	AddNode(ctx context.Context, node Node) error

	// GetNode 获取节点
	GetNode(ctx context.Context, id string) (Node, error)

	// AddEdge 添加边
	AddEdge(ctx context.Context, edge Edge) error

	// GetEdge 获取边
	GetEdge(ctx context.Context, id string) (Edge, error)

	// GetNodeNeighbors 获取节点的邻居
	GetNodeNeighbors(ctx context.Context, nodeID string, edgeType string, direction string, limit int) ([]Node, error)

	// GetNodeEdges 获取节点的边
	GetNodeEdges(ctx context.Context, nodeID string, edgeType string, direction string, limit int) ([]Edge, error)

	// DeleteNode 删除节点
	DeleteNode(ctx context.Context, id string) error

	// DeleteEdge 删除边
	DeleteEdge(ctx context.Context, id string) error

	// Close 关闭存储
	Close() error
}
