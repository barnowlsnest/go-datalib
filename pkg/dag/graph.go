// Package dag provides a directed acyclic graph (DAG) implementation with support
// for grouped nodes, edge management, and cycle detection.
package dag

import (
	"errors"
	"fmt"

	"github.com/barnowlsnest/go-datalib/pkg/list"
	"github.com/barnowlsnest/go-datalib/pkg/node"
	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

// Graph represents a directed graph with support for node grouping and acyclic verification.
// It provides efficient operations for adding/removing nodes and edges, querying relationships,
// and detecting cycles using Kahn's algorithm.
//
// Key features:
//   - Group-based node organization
//   - Bidirectional edge tracking (adjacency and back-references)
//   - Efficient cycle detection via topological sort
//   - Thread-unsafe: external synchronization required for concurrent access
type Graph struct {
	// name is the name of the graph.
	name string

	// id is the unique identifier of the graph.
	id ID

	// groups maps group names to sets of node IDs belonging to each group.
	// This allows for efficient group-based operations and queries.
	groups map[GroupName]map[NodeID]struct{}

	// backRefs maps each node to the set of nodes that have edges pointing to it.
	// This enables efficient reverse traversal and dependency analysis.
	backRefs map[NodeID]map[NodeID]struct{}

	// adjacency maps each source node to its outgoing edges.
	// The inner map associates destination nodes with edge IDs.
	adjacency map[NodeID]map[NodeID]EdgeID
}

// New creates and returns a new empty Graph instance with initialized internal maps.
func New() *Graph {
	return &Graph{
		groups:    make(map[GroupName]map[NodeID]struct{}),
		backRefs:  make(map[NodeID]map[NodeID]struct{}),
		adjacency: make(map[NodeID]map[NodeID]EdgeID),
	}
}

// Name returns the graph's name.
func (g *Graph) Name() string {
	return g.name
}

// ID returns the graph's unique identifier.
func (g *Graph) ID() ID {
	return g.id
}

// checkNodeExists verifies that a node exists in the specified group.
// Returns ErrGroupNotFound if the group doesn't exist, or ErrNodeNotFound if the node
// doesn't exist in the group.
func (g *Graph) checkNodeExists(n GroupNode) error {
	groupNodes, groupExists := g.groups[n.Group]
	if !groupExists {
		return errors.Join(ErrGroupNotFound, fmt.Errorf("group [%s]", n.Group))
	}
	_, nodeExists := groupNodes[n.ID]
	if !nodeExists {
		return errors.Join(ErrNodeNotFound, fmt.Errorf("group [%s] node [%d]", n.Group, n.ID))
	}
	return nil
}

// forEachEdge iterates over all outgoing edges from the specified node, invoking the
// provided callback function for each edge. Panics in the callback are recovered and
// passed to the callback as errors joined with ErrRecoverFromPanic.
func (g *Graph) forEachEdge(from NodeID, fn OnAdjacencyEdgeFn) {
	for to, edge := range g.adjacency[from] {
		func() {
			defer func() {
				if r := recover(); r != nil {
					var err error
					switch v := r.(type) {
					case error:
						err = v
					default:
						err = fmt.Errorf("recovered: %v", r)
					}
					fn(AdjacencyEdge{}, errors.Join(ErrRecoverFromPanic, err))
				}
			}()
			fn(AdjacencyEdge{
				From: from,
				To:   to,
				Edge: edge,
			}, nil)
		}()
	}
}

// removeAdjacency removes the edge from 'from' to 'to' and cleans up empty maps.
// This is a low-level helper that doesn't validate node existence.
func (g *Graph) removeAdjacency(from, to NodeID) {
	delete(g.adjacency[from], to)
	if len(g.adjacency[from]) == 0 {
		delete(g.adjacency, from)
	}
	delete(g.backRefs[to], from)
	if len(g.backRefs[to]) == 0 {
		delete(g.backRefs, to)
	}
}

// AddGroup creates a new group with the specified name.
// Returns ErrGroupAlreadyExists if a group with the same name already exists.
func (g *Graph) AddGroup(name GroupName) error {
	_, groupExists := g.groups[name]
	if groupExists {
		return errors.Join(ErrGroupAlreadyExists, fmt.Errorf("group [%s]", name))
	}
	g.groups[name] = make(map[NodeID]struct{})
	return nil
}

// AddNode adds a node to the specified group.
// Returns ErrGroupNotFound if the group doesn't exist.
// The node can be added multiple times without error (idempotent).
func (g *Graph) AddNode(n GroupNode) error {
	_, groupExists := g.groups[n.Group]
	if !groupExists {
		return errors.Join(ErrGroupNotFound, fmt.Errorf("group [%s]", n.Group))
	}
	g.groups[n.Group][n.ID] = struct{}{}
	return nil
}

// RemoveNode removes a node from its group and deletes all edges connected to it
// (both incoming and outgoing).
// Returns an error if the node doesn't exist.
func (g *Graph) RemoveNode(gn GroupNode) error {
	if nodeErr := g.checkNodeExists(gn); nodeErr != nil {
		return errors.Join(ErrInvalidEdge, nodeErr)
	}
	g.forEachEdge(gn.ID, func(a AdjacencyEdge, err error) {
		g.removeAdjacency(a.From, a.To)
	})
	delete(g.groups[gn.Group], gn.ID)
	return nil
}

// AddEdge creates a directed edge from 'from' to 'to'.
// The edge ID is computed as NSum(from.ID, to.ID).
// Returns ErrInvalidEdge if either node doesn't exist.
// Adding the same edge multiple times is idempotent.
func (g *Graph) AddEdge(from, to GroupNode) error {
	if fromErr := g.checkNodeExists(from); fromErr != nil {
		return errors.Join(ErrInvalidEdge, fromErr)
	}
	if toErr := g.checkNodeExists(to); toErr != nil {
		return errors.Join(ErrInvalidEdge, toErr)
	}
	if _, hasNeighbours := g.adjacency[from.ID]; !hasNeighbours {
		g.adjacency[from.ID] = make(map[NodeID]EdgeID)
	}
	if _, hasRefs := g.backRefs[to.ID]; !hasRefs {
		g.backRefs[to.ID] = make(map[NodeID]struct{})
	}
	g.adjacency[from.ID][to.ID] = serial.NSum(from.ID, to.ID)
	g.backRefs[to.ID][from.ID] = struct{}{}
	return nil
}

// RemoveEdge deletes the directed edge from 'from' to 'to'.
// Returns ErrInvalidEdge if either node doesn't exist.
// Removing a non-existent edge is a no-op (idempotent).
func (g *Graph) RemoveEdge(from, to GroupNode) error {
	if fromErr := g.checkNodeExists(from); fromErr != nil {
		return errors.Join(ErrInvalidEdge, fromErr)
	}
	if toErr := g.checkNodeExists(to); toErr != nil {
		return errors.Join(ErrInvalidEdge, toErr)
	}
	g.removeAdjacency(from.ID, to.ID)
	return nil
}

// HasNode returns true if the node exists in the specified group.
func (g *Graph) HasNode(gn GroupNode) bool {
	if err := g.checkNodeExists(gn); err != nil {
		return false
	}
	return true
}

// HasEdge returns true if a directed edge exists from 'from' to 'to'.
func (g *Graph) HasEdge(from, to GroupNode) bool {
	if fromErr := g.checkNodeExists(from); fromErr != nil {
		return false
	}
	if toErr := g.checkNodeExists(to); toErr != nil {
		return false
	}
	toNodes, toNodeExists := g.adjacency[from.ID]
	if !toNodeExists {
		return false
	}
	if _, edgeExists := toNodes[to.ID]; !edgeExists {
		return false
	}
	return true
}

// IsAcyclic performs cycle detection using Kahn's algorithm (topological sort).
// It returns a channel that will receive true if the graph is acyclic, false otherwise.
// The check runs asynchronously in a goroutine. An empty graph is considered acyclic.
//
// Algorithm: Kahn's topological sort
//   - Compute in-degree for all nodes
//   - Process nodes with zero in-degree
//   - If all nodes are processed, graph is acyclic
//   - If some nodes remain, a cycle exists
//
// Time complexity: O(V + E) where V is nodes and E is edges
// Space complexity: O(V)
func (g *Graph) IsAcyclic() <-chan bool {
	ch := make(chan bool)

	go func() {
		defer close(ch)

		q := list.NewQueue()
		in := make(map[NodeID]int)

		// Collect all nodes from the graph (both with outgoing and incoming edges)
		allNodes := make(map[NodeID]struct{})

		// Add nodes with outgoing edges
		for nodeID := range g.adjacency {
			allNodes[nodeID] = struct{}{}
		}

		// Add nodes with incoming edges
		for nodeID := range g.backRefs {
			allNodes[nodeID] = struct{}{}
		}

		// If there are no nodes, the graph is empty and is acyclic
		if len(allNodes) == 0 {
			ch <- true
			return
		}

		// Initialize in-degree for all nodes
		for nodeID := range allNodes {
			refs, exists := g.backRefs[nodeID]
			if exists {
				in[nodeID] = len(refs)
			} else {
				in[nodeID] = 0
			}
		}

		// Enqueue nodes with no incoming edges
		for nodeID, degree := range in {
			if degree == 0 {
				q.Enqueue(node.New(nodeID, nil, nil))
			}
		}

		var result []NodeID

		for q.Size() > 0 {
			n := q.Dequeue()

			if n == nil {
				break
			}

			nodeID := n.ID()
			result = append(result, nodeID)

			// Update in-degrees of neighbors
			neighbors, hasNeighbors := g.adjacency[nodeID]
			if hasNeighbors {
				for neighbor := range neighbors {
					in[neighbor]--
					if in[neighbor] == 0 {
						q.Enqueue(node.New(neighbor, nil, nil))
					}
				}
			}
		}

		// If we processed all nodes, the graph is acyclic
		if len(result) == len(allNodes) {
			ch <- true
		} else {
			ch <- false
		}
	}()

	return ch
}

// ForEachNeighbour iterates over all outgoing edges from the specified node,
// invoking the callback function for each neighbor.
// Returns ErrInvalidAdjacency if the node doesn't exist.
//
// The callback receives AdjacencyEdge with From, To, and Edge ID.
// Panics in the callback are recovered and passed as errors.
func (g *Graph) ForEachNeighbour(gn GroupNode, fn OnAdjacencyEdgeFn) error {
	if nodeErr := g.checkNodeExists(gn); nodeErr != nil {
		return errors.Join(ErrInvalidAdjacency, nodeErr)
	}
	g.forEachEdge(gn.ID, fn)
	return nil
}

// GetBackRefsOf returns all nodes that have edges pointing to the specified node.
// Returns ErrInvalidBackRef if the node doesn't exist or has no incoming edges.
//
// Note: The returned slice order is non-deterministic due to map iteration.
func (g *Graph) GetBackRefsOf(gn GroupNode) ([]GroupNode, error) {
	if nodeErr := g.checkNodeExists(gn); nodeErr != nil {
		return nil, errors.Join(ErrInvalidBackRef, nodeErr)
	}
	backRefs, hasBackRefs := g.backRefs[gn.ID]
	if !hasBackRefs {
		return nil, ErrInvalidBackRef
	}
	res := make([]GroupNode, len(backRefs))
	var i int
	for ref := range backRefs {
		for group, nodes := range g.groups {
			if _, exists := nodes[ref]; exists {
				res[i] = GroupNode{ref, group}
			}
		}
		i++
	}
	return res, nil
}

// GetNodes returns all nodes belonging to the specified group.
// Returns ErrGroupNotFound if the group doesn't exist.
//
// Note: The returned slice order is non-deterministic due to map iteration.
func (g *Graph) GetNodes(group GroupName) ([]GroupNode, error) {
	groupNodes, groupExists := g.groups[group]
	if !groupExists {
		return nil, errors.Join(ErrGroupNotFound, fmt.Errorf("group [%s]", group))
	}
	var i int
	res := make([]GroupNode, len(groupNodes))
	for n := range groupNodes {
		res[i] = GroupNode{n, group}
		i++
	}
	return res, nil
}

// ListGroups returns all group names in the graph.
//
// Note: The returned slice order is non-deterministic due to map iteration.
func (g *Graph) ListGroups() []GroupName {
	res := make([]GroupName, len(g.groups))
	var i = 0
	for name := range g.groups {
		res[i] = name
		i++
	}
	return res
}
