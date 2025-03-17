package graph

import (
	"context"
	"fmt"
	"sync"

	"github.com/HerSophia/go-lightrag/pkg/interfaces"
)

// MemoryGraphStorage 内存图存储实现
type MemoryGraphStorage struct {
	mu    sync.RWMutex
	nodes map[string]interfaces.Node
	edges map[string]interfaces.Edge
	// 节点邻接表，用于快速查找节点的边和邻居
	adjList map[string]map[string][]string // nodeID -> edgeType -> []edgeID
}

// NewMemoryGraphStorage 创建新的内存图存储
func NewMemoryGraphStorage() *MemoryGraphStorage {
	return &MemoryGraphStorage{
		nodes:   make(map[string]interfaces.Node),
		edges:   make(map[string]interfaces.Edge),
		adjList: make(map[string]map[string][]string),
	}
}

// AddNode 添加节点
func (s *MemoryGraphStorage) AddNode(ctx context.Context, node interfaces.Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查节点ID是否为空
	if node.ID == "" {
		return interfaces.ErrInvalidID
	}

	// 存储节点
	s.nodes[node.ID] = node

	// 初始化邻接表
	if _, exists := s.adjList[node.ID]; !exists {
		s.adjList[node.ID] = make(map[string][]string)
	}

	return nil
}

// GetNode 获取节点
func (s *MemoryGraphStorage) GetNode(ctx context.Context, id string) (interfaces.Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	node, exists := s.nodes[id]
	if !exists {
		return interfaces.Node{}, interfaces.ErrNotFound
	}

	return node, nil
}

// AddEdge 添加边
func (s *MemoryGraphStorage) AddEdge(ctx context.Context, edge interfaces.Edge) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查边ID是否为空
	if edge.ID == "" {
		return interfaces.ErrInvalidID
	}

	// 检查源节点和目标节点是否存在
	if _, exists := s.nodes[edge.SourceID]; !exists {
		return fmt.Errorf("source node %s not found", edge.SourceID)
	}

	if _, exists := s.nodes[edge.TargetID]; !exists {
		return fmt.Errorf("target node %s not found", edge.TargetID)
	}

	// 存储边
	s.edges[edge.ID] = edge

	// 更新邻接表
	if _, exists := s.adjList[edge.SourceID][edge.Type]; !exists {
		if s.adjList[edge.SourceID] == nil {
			s.adjList[edge.SourceID] = make(map[string][]string)
		}
		s.adjList[edge.SourceID][edge.Type] = []string{edge.ID}
	} else {
		s.adjList[edge.SourceID][edge.Type] = append(s.adjList[edge.SourceID][edge.Type], edge.ID)
	}

	// 如果是无向边，也更新目标节点的邻接表
	if !edge.Directional {
		if _, exists := s.adjList[edge.TargetID][edge.Type]; !exists {
			if s.adjList[edge.TargetID] == nil {
				s.adjList[edge.TargetID] = make(map[string][]string)
			}
			s.adjList[edge.TargetID][edge.Type] = []string{edge.ID}
		} else {
			s.adjList[edge.TargetID][edge.Type] = append(s.adjList[edge.TargetID][edge.Type], edge.ID)
		}
	}

	return nil
}

// GetEdge 获取边
func (s *MemoryGraphStorage) GetEdge(ctx context.Context, id string) (interfaces.Edge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	edge, exists := s.edges[id]
	if !exists {
		return interfaces.Edge{}, interfaces.ErrNotFound
	}

	return edge, nil
}

// GetNodeNeighbors 获取节点的邻居
func (s *MemoryGraphStorage) GetNodeNeighbors(ctx context.Context, nodeID string, edgeType string, direction string, limit int) ([]interfaces.Node, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查节点是否存在
	if _, exists := s.nodes[nodeID]; !exists {
		return nil, interfaces.ErrNotFound
	}

	neighbors := make([]interfaces.Node, 0)
	visited := make(map[string]bool)

	// 处理出边
	if direction == "out" || direction == "both" {
		for _, edgeID := range s.adjList[nodeID][edgeType] {
			edge := s.edges[edgeID]
			if !visited[edge.TargetID] {
				visited[edge.TargetID] = true
				neighbors = append(neighbors, s.nodes[edge.TargetID])
			}
		}
	}

	// 处理入边
	if direction == "in" || direction == "both" {
		for _, edge := range s.edges {
			if edge.TargetID == nodeID && (edgeType == "" || edge.Type == edgeType) {
				if !visited[edge.SourceID] {
					visited[edge.SourceID] = true
					neighbors = append(neighbors, s.nodes[edge.SourceID])
				}
			}
		}
	}

	// 限制结果数量
	if limit > 0 && len(neighbors) > limit {
		neighbors = neighbors[:limit]
	}

	return neighbors, nil
}

// GetNodeEdges 获取节点的边
func (s *MemoryGraphStorage) GetNodeEdges(ctx context.Context, nodeID string, edgeType string, direction string, limit int) ([]interfaces.Edge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查节点是否存在
	if _, exists := s.nodes[nodeID]; !exists {
		return nil, interfaces.ErrNotFound
	}

	edges := make([]interfaces.Edge, 0)

	// 处理出边
	if direction == "out" || direction == "both" {
		for _, edgeID := range s.adjList[nodeID][edgeType] {
			edges = append(edges, s.edges[edgeID])
		}
	}

	// 处理入边
	if direction == "in" || direction == "both" {
		for _, edge := range s.edges {
			if edge.TargetID == nodeID && (edgeType == "" || edge.Type == edgeType) {
				edges = append(edges, edge)
			}
		}
	}

	// 限制结果数量
	if limit > 0 && len(edges) > limit {
		edges = edges[:limit]
	}

	return edges, nil
}

// DeleteNode 删除节点
func (s *MemoryGraphStorage) DeleteNode(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查节点是否存在
	if _, exists := s.nodes[id]; !exists {
		return interfaces.ErrNotFound
	}

	// 删除与该节点相关的所有边
	for edgeID, edge := range s.edges {
		if edge.SourceID == id || edge.TargetID == id {
			delete(s.edges, edgeID)
		}
	}

	// 删除邻接表中的记录
	delete(s.adjList, id)

	// 删除节点
	delete(s.nodes, id)

	return nil
}

// DeleteEdge 删除边
func (s *MemoryGraphStorage) DeleteEdge(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查边是否存在
	edge, exists := s.edges[id]
	if !exists {
		return interfaces.ErrNotFound
	}

	// 从邻接表中删除边
	for i, edgeID := range s.adjList[edge.SourceID][edge.Type] {
		if edgeID == id {
			s.adjList[edge.SourceID][edge.Type] = append(
				s.adjList[edge.SourceID][edge.Type][:i],
				s.adjList[edge.SourceID][edge.Type][i+1:]...,
			)
			break
		}
	}

	// 如果是无向边，也从目标节点的邻接表中删除
	if !edge.Directional {
		for i, edgeID := range s.adjList[edge.TargetID][edge.Type] {
			if edgeID == id {
				s.adjList[edge.TargetID][edge.Type] = append(
					s.adjList[edge.TargetID][edge.Type][:i],
					s.adjList[edge.TargetID][edge.Type][i+1:]...,
				)
				break
			}
		}
	}

	// 删除边
	delete(s.edges, id)

	return nil
}

// Close 关闭存储
func (s *MemoryGraphStorage) Close() error {
	// 内存实现不需要特殊关闭操作
	return nil
}
