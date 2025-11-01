package dag

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// BasicFunctionalityTestSuite tests core DAG functionality
type BasicFunctionalityTestSuite struct {
	suite.Suite
}

func (s *BasicFunctionalityTestSuite) TestNewAdjacencyGroups() {
	ag := New()

	s.Require().NotNil(ag)
	s.Require().NotNil(ag.groups)
	s.Require().NotNil(ag.backRefs)
	s.Require().NotNil(ag.adjacency)
	s.Require().Equal(0, len(ag.groups))
}

func (s *BasicFunctionalityTestSuite) TestAddGroup() {
	ag := New()

	err := ag.AddGroup("users")
	s.Require().NoError(err)

	// Adding same group should return error
	err = ag.AddGroup("users")
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrGroupAlreadyExists)
}

func (s *BasicFunctionalityTestSuite) TestAddNode() {
	ag := New()
	_ = ag.AddGroup("users")

	node := GroupNode{ID: 1, Group: "users"}
	err := ag.AddNode(node)
	s.Require().NoError(err)

	// Adding to non-existent group should fail
	invalidNode := GroupNode{ID: 2, Group: "nonexistent"}
	err = ag.AddNode(invalidNode)
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrGroupNotFound)
}

func (s *BasicFunctionalityTestSuite) TestAddNode_Idempotent() {
	ag := New()
	_ = ag.AddGroup("users")

	node := GroupNode{ID: 1, Group: "users"}
	err1 := ag.AddNode(node)
	err2 := ag.AddNode(node)

	s.Require().NoError(err1)
	s.Require().NoError(err2)
}

func (s *BasicFunctionalityTestSuite) TestHasNode() {
	ag := New()
	_ = ag.AddGroup("users")

	node := GroupNode{ID: 1, Group: "users"}
	s.Require().False(ag.HasNode(node))

	_ = ag.AddNode(node)
	s.Require().True(ag.HasNode(node))
}

func (s *BasicFunctionalityTestSuite) TestAddEdge() {
	ag := New()
	_ = ag.AddGroup("users")

	from := GroupNode{ID: 1, Group: "users"}
	to := GroupNode{ID: 2, Group: "users"}
	_ = ag.AddNode(from)
	_ = ag.AddNode(to)

	err := ag.AddEdge(from, to)
	s.Require().NoError(err)
	s.Require().True(ag.HasEdge(from, to))
}

func (s *BasicFunctionalityTestSuite) TestAddEdge_NonExistentNode() {
	ag := New()
	_ = ag.AddGroup("users")

	from := GroupNode{ID: 1, Group: "users"}
	to := GroupNode{ID: 2, Group: "users"}
	_ = ag.AddNode(from)

	err := ag.AddEdge(from, to)
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrInvalidEdge)
}

func (s *BasicFunctionalityTestSuite) TestRemoveEdge() {
	ag := New()
	_ = ag.AddGroup("users")

	from := GroupNode{ID: 1, Group: "users"}
	to := GroupNode{ID: 2, Group: "users"}
	_ = ag.AddNode(from)
	_ = ag.AddNode(to)
	_ = ag.AddEdge(from, to)

	s.Require().True(ag.HasEdge(from, to))

	err := ag.RemoveEdge(from, to)
	s.Require().NoError(err)
	s.Require().False(ag.HasEdge(from, to))
}

func (s *BasicFunctionalityTestSuite) TestRemoveNode() {
	ag := New()
	_ = ag.AddGroup("users")

	node1 := GroupNode{ID: 1, Group: "users"}
	node2 := GroupNode{ID: 2, Group: "users"}
	node3 := GroupNode{ID: 3, Group: "users"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)
	_ = ag.AddNode(node3)

	_ = ag.AddEdge(node1, node2)
	_ = ag.AddEdge(node2, node3)
	_ = ag.AddEdge(node3, node1)

	err := ag.RemoveNode(node2)
	s.Require().NoError(err)
	s.Require().False(ag.HasNode(node2))
	s.Require().False(ag.HasEdge(node1, node2))
	s.Require().False(ag.HasEdge(node2, node3))
}

// MemoryConsistencyTestSuite tests memory cleanup and consistency
type MemoryConsistencyTestSuite struct {
	suite.Suite
}

func (s *MemoryConsistencyTestSuite) TestRemoveEdge_CleansUpEmptyMaps() {
	ag := New()
	_ = ag.AddGroup("test")

	from := GroupNode{ID: 1, Group: "test"}
	to := GroupNode{ID: 2, Group: "test"}
	_ = ag.AddNode(from)
	_ = ag.AddNode(to)

	_ = ag.AddEdge(from, to)
	s.Require().Equal(1, len(ag.adjacency))
	s.Require().Equal(1, len(ag.backRefs))

	_ = ag.RemoveEdge(from, to)

	// Empty maps should be cleaned up
	s.Require().Equal(0, len(ag.adjacency))
	s.Require().Equal(0, len(ag.backRefs))
}

func (s *MemoryConsistencyTestSuite) TestRemoveNode_CleansUpAllReferences() {
	ag := New()
	_ = ag.AddGroup("test")

	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	node3 := GroupNode{ID: 3, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)
	_ = ag.AddNode(node3)

	// Create edge structure: node1 -> node2 -> node3
	_ = ag.AddEdge(node1, node2)
	_ = ag.AddEdge(node2, node3)

	_ = ag.RemoveNode(node2)

	// Verify node2 is removed from the group
	s.Require().False(ag.HasNode(node2), "node2 should be removed from group")

	// Verify all outgoing edges from node2 are cleaned up
	_, hasAdjacency := ag.adjacency[node2.ID]
	s.Require().False(hasAdjacency, "node2 should not have adjacency entries")

	// Verify edge from node2 to node3 is cleaned up in node3's backRefs
	if backRefs, exists := ag.backRefs[node3.ID]; exists {
		_, hasNode2Ref := backRefs[node2.ID]
		s.Require().False(hasNode2Ref, "node3 should not have backRef from node2")
	}
}

func (s *MemoryConsistencyTestSuite) TestBackRefsConsistency() {
	ag := New()
	_ = ag.AddGroup("test")

	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)

	_ = ag.AddEdge(node1, node2)

	// BackRefs should be consistent with adjacency
	backRefs, exists := ag.backRefs[node2.ID]
	s.Require().True(exists)
	_, hasRef := backRefs[node1.ID]
	s.Require().True(hasRef)
}

func (s *MemoryConsistencyTestSuite) TestMultipleEdgeAdditions_Consistency() {
	ag := New()
	_ = ag.AddGroup("test")

	from := GroupNode{ID: 1, Group: "test"}
	to := GroupNode{ID: 2, Group: "test"}
	_ = ag.AddNode(from)
	_ = ag.AddNode(to)

	// Add same edge multiple times
	_ = ag.AddEdge(from, to)
	_ = ag.AddEdge(from, to)
	_ = ag.AddEdge(from, to)

	// Should only have one edge
	edges := ag.adjacency[from.ID]
	s.Require().Equal(1, len(edges))

	backRefs := ag.backRefs[to.ID]
	s.Require().Equal(1, len(backRefs))
}

// IsAcyclicCorrectnessTestSuite tests cycle detection correctness
type IsAcyclicCorrectnessTestSuite struct {
	suite.Suite
}

func (s *IsAcyclicCorrectnessTestSuite) TestIsAcyclic_EmptyGraph() {
	ag := New()

	result := <-ag.IsAcyclic()
	s.Require().True(result, "empty graph should be acyclic")
}

func (s *IsAcyclicCorrectnessTestSuite) TestIsAcyclic_SingleNode() {
	ag := New()
	_ = ag.AddGroup("test")

	node := GroupNode{ID: 1, Group: "test"}
	_ = ag.AddNode(node)

	result := <-ag.IsAcyclic()
	s.Require().True(result, "single node with no edges should be acyclic")
}

func (s *IsAcyclicCorrectnessTestSuite) TestIsAcyclic_TwoNodesNoEdges() {
	ag := New()
	_ = ag.AddGroup("test")

	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)

	result := <-ag.IsAcyclic()
	s.Require().True(result, "nodes without edges should be acyclic")
}

func (s *IsAcyclicCorrectnessTestSuite) TestIsAcyclic_SimpleChain() {
	ag := New()
	_ = ag.AddGroup("test")

	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	node3 := GroupNode{ID: 3, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)
	_ = ag.AddNode(node3)

	_ = ag.AddEdge(node1, node2)
	_ = ag.AddEdge(node2, node3)

	result := <-ag.IsAcyclic()
	s.Require().True(result, "simple chain should be acyclic")
}

func (s *IsAcyclicCorrectnessTestSuite) TestIsAcyclic_SimpleCycle() {
	ag := New()
	_ = ag.AddGroup("test")

	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	node3 := GroupNode{ID: 3, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)
	_ = ag.AddNode(node3)

	// Create cycle: 1 -> 2 -> 3 -> 1
	_ = ag.AddEdge(node1, node2)
	_ = ag.AddEdge(node2, node3)
	_ = ag.AddEdge(node3, node1)

	result := <-ag.IsAcyclic()
	s.Require().False(result, "cycle should be detected")
}

func (s *IsAcyclicCorrectnessTestSuite) TestIsAcyclic_SelfLoop() {
	ag := New()
	_ = ag.AddGroup("test")

	node := GroupNode{ID: 1, Group: "test"}
	_ = ag.AddNode(node)
	_ = ag.AddEdge(node, node)

	result := <-ag.IsAcyclic()
	s.Require().False(result, "self-loop should be detected as cycle")
}

func (s *IsAcyclicCorrectnessTestSuite) TestIsAcyclic_DAGWithMultiplePaths() {
	ag := New()
	_ = ag.AddGroup("test")

	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	node3 := GroupNode{ID: 3, Group: "test"}
	node4 := GroupNode{ID: 4, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)
	_ = ag.AddNode(node3)
	_ = ag.AddNode(node4)

	// Diamond shape: 1 -> 2, 1 -> 3, 2 -> 4, 3 -> 4
	_ = ag.AddEdge(node1, node2)
	_ = ag.AddEdge(node1, node3)
	_ = ag.AddEdge(node2, node4)
	_ = ag.AddEdge(node3, node4)

	result := <-ag.IsAcyclic()
	s.Require().True(result, "DAG with multiple paths should be acyclic")
}

func (s *IsAcyclicCorrectnessTestSuite) TestIsAcyclic_ComplexCycle() {
	ag := New()
	_ = ag.AddGroup("test")

	nodes := make([]GroupNode, 10)
	for i := 0; i < 10; i++ {
		nodes[i] = GroupNode{ID: uint64(i + 1), Group: "test"}
		_ = ag.AddNode(nodes[i])
	}

	// Create chain with cycle at the end
	for i := 0; i < 9; i++ {
		_ = ag.AddEdge(nodes[i], nodes[i+1])
	}
	// Close the cycle
	_ = ag.AddEdge(nodes[9], nodes[5])

	result := <-ag.IsAcyclic()
	s.Require().False(result, "complex cycle should be detected")
}

func (s *IsAcyclicCorrectnessTestSuite) TestIsAcyclic_DisconnectedComponents() {
	ag := New()
	_ = ag.AddGroup("test")

	// Component 1: 1 -> 2 -> 3
	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	node3 := GroupNode{ID: 3, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)
	_ = ag.AddNode(node3)
	_ = ag.AddEdge(node1, node2)
	_ = ag.AddEdge(node2, node3)

	// Component 2: 4 -> 5 -> 6
	node4 := GroupNode{ID: 4, Group: "test"}
	node5 := GroupNode{ID: 5, Group: "test"}
	node6 := GroupNode{ID: 6, Group: "test"}
	_ = ag.AddNode(node4)
	_ = ag.AddNode(node5)
	_ = ag.AddNode(node6)
	_ = ag.AddEdge(node4, node5)
	_ = ag.AddEdge(node5, node6)

	result := <-ag.IsAcyclic()
	s.Require().True(result, "disconnected acyclic components should be acyclic")
}

func (s *IsAcyclicCorrectnessTestSuite) TestIsAcyclic_DisconnectedWithCycle() {
	ag := New()
	_ = ag.AddGroup("test")

	// Component 1: 1 -> 2 -> 3 (acyclic)
	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	node3 := GroupNode{ID: 3, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)
	_ = ag.AddNode(node3)
	_ = ag.AddEdge(node1, node2)
	_ = ag.AddEdge(node2, node3)

	// Component 2: 4 -> 5 -> 6 -> 4 (cycle)
	node4 := GroupNode{ID: 4, Group: "test"}
	node5 := GroupNode{ID: 5, Group: "test"}
	node6 := GroupNode{ID: 6, Group: "test"}
	_ = ag.AddNode(node4)
	_ = ag.AddNode(node5)
	_ = ag.AddNode(node6)
	_ = ag.AddEdge(node4, node5)
	_ = ag.AddEdge(node5, node6)
	_ = ag.AddEdge(node6, node4)

	result := <-ag.IsAcyclic()
	s.Require().False(result, "cycle in one component should be detected")
}

// IsAcyclicPerformanceTestSuite tests performance of cycle detection
type IsAcyclicPerformanceTestSuite struct {
	suite.Suite
}

func (s *IsAcyclicPerformanceTestSuite) TestIsAcyclic_LargeAcyclicGraph() {
	ag := New()
	_ = ag.AddGroup("test")

	// Create a large acyclic graph with 1000 nodes in layers
	numNodes := 1000
	nodes := make([]GroupNode, numNodes)

	for i := 0; i < numNodes; i++ {
		nodes[i] = GroupNode{ID: uint64(i + 1), Group: "test"}
		_ = ag.AddNode(nodes[i])
	}

	// Connect each node to the next (linear chain)
	for i := 0; i < numNodes-1; i++ {
		_ = ag.AddEdge(nodes[i], nodes[i+1])
	}

	start := time.Now()
	result := <-ag.IsAcyclic()
	duration := time.Since(start)

	s.Require().True(result)
	// Should complete quickly (< 100ms for 1000 nodes)
	s.Require().Less(duration, 100*time.Millisecond,
		"IsAcyclic should complete quickly for 1000-node chain")
}

func (s *IsAcyclicPerformanceTestSuite) TestIsAcyclic_ComplexDAG() {
	ag := New()
	_ = ag.AddGroup("test")

	// Create a complex DAG with 100 nodes and many edges
	numNodes := 100
	nodes := make([]GroupNode, numNodes)

	for i := 0; i < numNodes; i++ {
		nodes[i] = GroupNode{ID: uint64(i + 1), Group: "test"}
		_ = ag.AddNode(nodes[i])
	}

	// Create edges from each node to several later nodes
	for i := 0; i < numNodes-1; i++ {
		for j := i + 1; j < numNodes && j < i+5; j++ {
			_ = ag.AddEdge(nodes[i], nodes[j])
		}
	}

	start := time.Now()
	result := <-ag.IsAcyclic()
	duration := time.Since(start)

	s.Require().True(result)
	s.Require().Less(duration, 50*time.Millisecond,
		"IsAcyclic should complete quickly for complex DAG")
}

// BackRefsTestSuite tests back-reference functionality
type BackRefsTestSuite struct {
	suite.Suite
}

func (s *BackRefsTestSuite) TestGetBackRefsOf() {
	ag := New()
	_ = ag.AddGroup("test")

	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	node3 := GroupNode{ID: 3, Group: "test"}
	node4 := GroupNode{ID: 4, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)
	_ = ag.AddNode(node3)
	_ = ag.AddNode(node4)

	// Multiple nodes point to node4
	_ = ag.AddEdge(node1, node4)
	_ = ag.AddEdge(node2, node4)
	_ = ag.AddEdge(node3, node4)

	backRefs, err := ag.GetBackRefsOf(node4)
	s.Require().NoError(err)
	s.Require().Equal(3, len(backRefs))
}

func (s *BackRefsTestSuite) TestGetBackRefsOf_NoBackRefs() {
	ag := New()
	_ = ag.AddGroup("test")

	node := GroupNode{ID: 1, Group: "test"}
	_ = ag.AddNode(node)

	backRefs, err := ag.GetBackRefsOf(node)
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrInvalidBackRef)
	s.Require().Nil(backRefs)
}

// ForEachNeighbourTestSuite tests neighbor iteration
type ForEachNeighbourTestSuite struct {
	suite.Suite
}

func (s *ForEachNeighbourTestSuite) TestForEachNeighbour() {
	ag := New()
	_ = ag.AddGroup("test")

	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	node3 := GroupNode{ID: 3, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)
	_ = ag.AddNode(node3)

	_ = ag.AddEdge(node1, node2)
	_ = ag.AddEdge(node1, node3)

	visited := make([]NodeID, 0)
	err := ag.ForEachNeighbour(node1, func(edge AdjacencyEdge, err error) {
		s.Require().NoError(err)
		visited = append(visited, edge.To)
	})

	s.Require().NoError(err)
	s.Require().Equal(2, len(visited))
}

func (s *ForEachNeighbourTestSuite) TestForEachNeighbour_PanicRecovery() {
	ag := New()
	_ = ag.AddGroup("test")

	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)
	_ = ag.AddEdge(node1, node2)

	var recoveredError error
	err := ag.ForEachNeighbour(node1, func(edge AdjacencyEdge, err error) {
		if err != nil {
			recoveredError = err
			return
		}
		panic("intentional panic")
	})

	s.Require().NoError(err)
	s.Require().NotNil(recoveredError)
	s.Require().ErrorIs(recoveredError, ErrRecoverFromPanic)
}

// GroupOperationsTestSuite tests group-related operations
type GroupOperationsTestSuite struct {
	suite.Suite
}

func (s *GroupOperationsTestSuite) TestListGroups() {
	ag := New()

	groups := ag.ListGroups()
	s.Require().Equal(0, len(groups))

	_ = ag.AddGroup("users")
	_ = ag.AddGroup("products")

	groups = ag.ListGroups()
	s.Require().Equal(2, len(groups))
}

func (s *GroupOperationsTestSuite) TestGetNodes() {
	ag := New()
	_ = ag.AddGroup("test")

	node1 := GroupNode{ID: 1, Group: "test"}
	node2 := GroupNode{ID: 2, Group: "test"}
	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)

	nodes, err := ag.GetNodes("test")
	s.Require().NoError(err)
	s.Require().Equal(2, len(nodes))
}

func (s *GroupOperationsTestSuite) TestGetNodes_NonExistentGroup() {
	ag := New()

	nodes, err := ag.GetNodes("nonexistent")
	s.Require().Error(err)
	s.Require().ErrorIs(err, ErrGroupNotFound)
	s.Require().Nil(nodes)
}

// ConcurrencyTestSuite tests concurrent operations
type ConcurrencyTestSuite struct {
	suite.Suite
}

func (s *ConcurrencyTestSuite) TestConcurrentIsAcyclic() {
	ag := New()
	_ = ag.AddGroup("test")

	// Build graph
	numNodes := 100
	nodes := make([]GroupNode, numNodes)
	for i := 0; i < numNodes; i++ {
		nodes[i] = GroupNode{ID: uint64(i + 1), Group: "test"}
		_ = ag.AddNode(nodes[i])
	}
	for i := 0; i < numNodes-1; i++ {
		_ = ag.AddEdge(nodes[i], nodes[i+1])
	}

	// Run multiple concurrent IsAcyclic checks
	var wg sync.WaitGroup
	results := make([]bool, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			results[idx] = <-ag.IsAcyclic()
		}(i)
	}

	wg.Wait()

	// All should return true
	for i, result := range results {
		s.Require().True(result, "result %d should be true", i)
	}
}

// Benchmark tests
func BenchmarkIsAcyclic_SmallGraph(b *testing.B) {
	ag := New()
	_ = ag.AddGroup("test")

	// 10 nodes linear chain
	nodes := make([]GroupNode, 10)
	for i := 0; i < 10; i++ {
		nodes[i] = GroupNode{ID: uint64(i + 1), Group: "test"}
		_ = ag.AddNode(nodes[i])
	}
	for i := 0; i < 9; i++ {
		_ = ag.AddEdge(nodes[i], nodes[i+1])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-ag.IsAcyclic()
	}
}

func BenchmarkIsAcyclic_MediumGraph(b *testing.B) {
	ag := New()
	_ = ag.AddGroup("test")

	// 100 nodes linear chain
	nodes := make([]GroupNode, 100)
	for i := 0; i < 100; i++ {
		nodes[i] = GroupNode{ID: uint64(i + 1), Group: "test"}
		_ = ag.AddNode(nodes[i])
	}
	for i := 0; i < 99; i++ {
		_ = ag.AddEdge(nodes[i], nodes[i+1])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-ag.IsAcyclic()
	}
}

func BenchmarkIsAcyclic_LargeGraph(b *testing.B) {
	ag := New()
	_ = ag.AddGroup("test")

	// 1000 nodes linear chain
	nodes := make([]GroupNode, 1000)
	for i := 0; i < 1000; i++ {
		nodes[i] = GroupNode{ID: uint64(i + 1), Group: "test"}
		_ = ag.AddNode(nodes[i])
	}
	for i := 0; i < 999; i++ {
		_ = ag.AddEdge(nodes[i], nodes[i+1])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-ag.IsAcyclic()
	}
}

func BenchmarkIsAcyclic_ComplexDAG(b *testing.B) {
	ag := New()
	_ = ag.AddGroup("test")

	// 100 nodes with multiple edges per node
	nodes := make([]GroupNode, 100)
	for i := 0; i < 100; i++ {
		nodes[i] = GroupNode{ID: uint64(i + 1), Group: "test"}
		_ = ag.AddNode(nodes[i])
	}
	for i := 0; i < 99; i++ {
		for j := i + 1; j < 100 && j < i+5; j++ {
			_ = ag.AddEdge(nodes[i], nodes[j])
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		<-ag.IsAcyclic()
	}
}

func BenchmarkAddEdge(b *testing.B) {
	ag := New()
	_ = ag.AddGroup("test")

	from := GroupNode{ID: 1, Group: "test"}
	to := GroupNode{ID: 2, Group: "test"}
	_ = ag.AddNode(from)
	_ = ag.AddNode(to)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ag.AddEdge(from, to)
	}
}

func BenchmarkRemoveEdge(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		ag := New()
		_ = ag.AddGroup("test")

		from := GroupNode{ID: 1, Group: "test"}
		to := GroupNode{ID: 2, Group: "test"}
		_ = ag.AddNode(from)
		_ = ag.AddNode(to)
		_ = ag.AddEdge(from, to)

		b.StartTimer()
		_ = ag.RemoveEdge(from, to)
		b.StopTimer()
	}
}

func BenchmarkHasEdge(b *testing.B) {
	ag := New()
	_ = ag.AddGroup("test")

	from := GroupNode{ID: 1, Group: "test"}
	to := GroupNode{ID: 2, Group: "test"}
	_ = ag.AddNode(from)
	_ = ag.AddNode(to)
	_ = ag.AddEdge(from, to)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ag.HasEdge(from, to)
	}
}

// Test suite runners
func TestBasicFunctionalityTestSuite(t *testing.T) {
	suite.Run(t, new(BasicFunctionalityTestSuite))
}

func TestMemoryConsistencyTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryConsistencyTestSuite))
}

func TestIsAcyclicCorrectnessTestSuite(t *testing.T) {
	suite.Run(t, new(IsAcyclicCorrectnessTestSuite))
}

func TestIsAcyclicPerformanceTestSuite(t *testing.T) {
	suite.Run(t, new(IsAcyclicPerformanceTestSuite))
}

func TestBackRefsTestSuite(t *testing.T) {
	suite.Run(t, new(BackRefsTestSuite))
}

func TestForEachNeighbourTestSuite(t *testing.T) {
	suite.Run(t, new(ForEachNeighbourTestSuite))
}

func TestGroupOperationsTestSuite(t *testing.T) {
	suite.Run(t, new(GroupOperationsTestSuite))
}

func TestConcurrencyTestSuite(t *testing.T) {
	suite.Run(t, new(ConcurrencyTestSuite))
}

// Example tests
func ExampleGraph_IsAcyclic() {
	ag := New()
	_ = ag.AddGroup("nodes")

	// Create a simple chain
	node1 := GroupNode{ID: 1, Group: "nodes"}
	node2 := GroupNode{ID: 2, Group: "nodes"}
	node3 := GroupNode{ID: 3, Group: "nodes"}

	_ = ag.AddNode(node1)
	_ = ag.AddNode(node2)
	_ = ag.AddNode(node3)

	_ = ag.AddEdge(node1, node2)
	_ = ag.AddEdge(node2, node3)

	isAcyclic := <-ag.IsAcyclic()
	fmt.Printf("Graph is acyclic: %t\n", isAcyclic)

	// Output:
	// Graph is acyclic: true
}
