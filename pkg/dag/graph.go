package dag

import (
	"errors"
	"fmt"

	"github.com/barnowlsnest/go-datalib/pkg/node"
	"github.com/barnowlsnest/go-datalib/pkg/queue"
	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

// Graph implements a graph data structure with grouped nodes.
//
// This structure provides an efficient representation of a directed graph
// where nodes are organized into named groups. It maintains three core
// data structures for optimal performance:
//   - groups: Maps group names to sets of nodes
//   - backRefs: Maps nodes to sets of nodes that reference them
//   - adjacency: Maps source nodes to their outgoing edges
//
// Key features:
//   - Efficient group-based node organization
//   - Fast edge lookup and traversal
//   - Cycle detection capabilities
//   - Back-reference tracking for reverse traversal
//   - Support for complex graph algorithms
//
// The structure is designed to support various graph algorithms including
// cycle detection, topological sorting, and efficient neighbor traversal.
//
// Thread Safety:
// Graph is not thread-safe. Concurrent access requires external
// synchronization mechanisms.
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

// New creates a new empty Graph instance.
//
// This constructor initializes all internal maps required for graph operations.
// The returned instance is ready for use and can immediately accept groups,
// nodes, and edges.
//
// Returns:
//   - A new Graph instance with initialized internal structures
//
// Example:
//
//	ag := New()
//	ag.AddGroup("users")
//	ag.AddNode(GroupNode{Id: 1, Group: "users"})
func New() *Graph {
	return &Graph{
		groups:    make(map[GroupName]map[NodeID]struct{}),
		backRefs:  make(map[NodeID]map[NodeID]struct{}),
		adjacency: make(map[NodeID]map[NodeID]EdgeID),
	}
}

// Name returns the name of the graph.
//
// This method returns the name of the graph. The name is used to identify
// the graph in logs and other contexts where it's useful to distinguish
// between multiple graphs.
//
// Returns:
//   - The name of the graph
func (g *Graph) Name() string {
	return g.name
}

// ID returns the unique ID of the graph.
//
// This method returns the unique ID of the graph. The ID is used to identify
// the graph in logs and other contexts where it's useful to distinguish
// between multiple graphs.
//
// Returns:
//   - The ID of the graph
func (g *Graph) ID() ID {
	return g.id
}

// checkNodeExists verifies that a node exists in its specified group.
//
// This internal helper method validates that both the group exists and
// the node is a member of that group. It's used throughout the codebase
// to ensure operations are performed on valid nodes.
//
// Parameters:
//   - n: The GroupNode to validate
//
// Returns:
//   - nil if the node exists in the specified group
//   - ErrGroupNotFound if the group doesn't exist
//   - ErrNodeNotFound if the node doesn't exist in the group
//
// Error Details:
// The returned errors include contextual information about which group
// and/or node caused the validation failure.
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

// forEachEdge iterates over all outgoing edges from a source node.
//
// This internal helper method safely executes a callback function for each
// edge originating from the specified node. It includes panic recovery to
// ensure that callback failures don't crash the entire operation.
//
// Parameters:
//   - from: The source node ID to iterate edges from
//   - fn: The callback function to execute for each edge
//
// Panic Recovery:
// If the callback function panics, the panic is recovered and converted
// to an error. The callback is then called with an empty AdjacencyEdge
// and an ErrRecoverFromPanic error.
//
// Behavior:
//   - If the source node has no outgoing edges, no callbacks are executed
//   - Each callback receives a properly constructed AdjacencyEdge
//   - Panics in callbacks are safely handled and reported as errors
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

// removeAdjacency removes an edge and its back-reference from the graph.
//
// This internal helper method performs the low-level cleanup required
// when removing edges. It updates both the adjacency list and back-reference
// maps, and cleans up empty map entries to prevent memory leaks.
//
// Parameters:
//   - from: The source node ID of the edge to remove
//   - to: The destination node ID of the edge to remove
//
// Cleanup Behavior:
//   - Removes the edge from the adjacency map
//   - Removes the back-reference from the backRefs map
//   - Deletes empty map entries to prevent memory leaks
//   - Safe to call even if the edge doesn't exist (no-op)
//
// Memory Management:
// This method ensures that empty map entries are cleaned up, preventing
// memory leaks in long-running applications with dynamic graph structures.
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

// AddGroup creates a new group in the graph.
//
// This method adds a new named group that can contain nodes. Groups provide
// a way to organize nodes into logical collections, enabling efficient
// group-based operations and queries.
//
// Parameters:
//   - name: The name of the group to create. Must be unique within the graph.
//
// Returns:
//   - nil if the group was created successfully
//   - ErrGroupAlreadyExists if a group with the same name already exists
//
// Example:
//
//	ag := New()
//	err := ag.AddGroup("users")
//	if err != nil {
//		log.Printf("Failed to add group: %v", err)
//	}
//
//	// Adding the same group again will return an error
//	err = ag.AddGroup("users") // Returns ErrGroupAlreadyExists
func (g *Graph) AddGroup(name GroupName) error {
	_, groupExists := g.groups[name]
	if groupExists {
		return errors.Join(ErrGroupAlreadyExists, fmt.Errorf("group [%s]", name))
	}
	g.groups[name] = make(map[NodeID]struct{})
	return nil
}

// AddNode adds a node to an existing group in the graph.
//
// This method adds a node with the specified ID to the given group.
// The group must already exist before nodes can be added to it.
// Nodes can participate in edges once they are added to the graph.
//
// Parameters:
//   - n: The GroupNode containing the node ID and group name
//
// Returns:
//   - nil if the node was added successfully
//   - ErrGroupNotFound if the specified group doesn't exist
//
// Behavior:
//   - If the node already exists in the group, this operation is idempotent
//   - The node becomes available for edge operations after addition
//   - No validation is performed on the node ID (any uint64 value is valid)
//
// Example:
//
//	ag := New()
//	ag.AddGroup("users")
//
//	node := GroupNode{Id: 1, Group: "users"}
//	err := ag.AddNode(node)
//	if err != nil {
//		log.Printf("Failed to add node: %v", err)
//	}
//
//	// Adding to non-existent group returns error
//	invalidNode := GroupNode{Id: 2, Group: "nonexistent"}
//	err = ag.AddNode(invalidNode) // Returns ErrGroupNotFound
func (g *Graph) AddNode(n GroupNode) error {
	_, groupExists := g.groups[n.Group]
	if !groupExists {
		return errors.Join(ErrGroupNotFound, fmt.Errorf("group [%s]", n.Group))
	}
	g.groups[n.Group][n.ID] = struct{}{}
	return nil
}

// RemoveNode removes a node and all its associated edges from the graph.
//
// This method removes the specified node from its group and automatically
// cleans up all edges (both incoming and outgoing) associated with the node.
// This ensures graph consistency and prevents dangling edge references.
//
// Parameters:
//   - node: The GroupNode to remove from the graph
//
// Returns:
//   - nil if the node was removed successfully
//   - ErrInvalidEdge joined with the underlying error if the node doesn't exist
//
// Cleanup Behavior:
//   - Removes all outgoing edges from the node
//   - Removes all incoming edges to the node
//   - Removes the node from its group
//   - Cleans up empty map entries to prevent memory leaks
//
// Example:
//
//	ag := New()
//	ag.AddGroup("users")
//	node := GroupNode{Id: 1, Group: "users"}
//	ag.AddNode(node)
//
//	err := ag.RemoveNode(node)
//	if err != nil {
//		log.Printf("Failed to remove node: %v", err)
//	}
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

// AddEdge creates a directed edge between two nodes in the graph.
//
// This method establishes a directed relationship from the source node to
// the destination node. Both nodes must exist in their respective groups
// before an edge can be created between them.
//
// Parameters:
//   - from: The source GroupNode for the edge
//   - to: The destination GroupNode for the edge
//
// Returns:
//   - nil if the edge was created successfully
//   - ErrInvalidEdge joined with the underlying error if either node doesn't exist
//
// Edge Properties:
//   - The edge ID is generated using NSum(from.Id, to.Id) for consistency
//   - Duplicate edges between the same nodes are allowed (overwrites existing)
//   - Back-references are automatically maintained for efficient reverse traversal
//
// Example:
//
//	ag := New()
//	ag.AddGroup("users")
//	from := GroupNode{Id: 1, Group: "users"}
//	to := GroupNode{Id: 2, Group: "users"}
//	ag.AddNode(from)
//	ag.AddNode(to)
//
//	err := ag.AddEdge(from, to)
//	if err != nil {
//		log.Printf("Failed to add edge: %v", err)
//	}
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

// RemoveEdge removes a directed edge between two nodes in the graph.
//
// This method removes the specified edge and its associated back-reference.
// Both nodes must exist for the operation to succeed, even if the edge
// doesn't exist (in which case it's a no-op).
//
// Parameters:
//   - from: The source GroupNode of the edge to remove
//   - to: The destination GroupNode of the edge to remove
//
// Returns:
//   - nil if the edge was removed successfully (or didn't exist)
//   - ErrInvalidEdge joined with the underlying error if either node doesn't exist
//
// Behavior:
//   - If the edge doesn't exist, the operation succeeds without error
//   - Back-references are automatically cleaned up
//   - Empty map entries are removed to prevent memory leaks
//
// Example:
//
//	ag := New()
//	// ... setup nodes and edge ...
//
//	err := ag.RemoveEdge(from, to)
//	if err != nil {
//		log.Printf("Failed to remove edge: %v", err)
//	}
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

// HasNode checks whether a node exists in the specified group.
//
// This method provides a simple boolean check for node existence without
// returning detailed error information. It's useful for conditional logic
// where you only need to know if a node exists.
//
// Parameters:
//   - node: The GroupNode to check for existence
//
// Returns:
//   - true if the node exists in the specified group
//   - false if the node or group doesn't exist
//
// Example:
//
//	ag := New()
//	ag.AddGroup("users")
//	node := GroupNode{Id: 1, Group: "users"}
//
//	if !ag.HasNode(node) {
//		ag.AddNode(node)
//	}
func (g *Graph) HasNode(gn GroupNode) bool {
	if err := g.checkNodeExists(gn); err != nil {
		return false
	}
	return true
}

// HasEdge checks whether a directed edge exists between two nodes.
//
// This method provides a simple boolean check for edge existence without
// returning detailed error information. Both nodes must exist for the
// method to return true.
//
// Parameters:
//   - from: The source GroupNode of the edge to check
//   - to: The destination GroupNode of the edge to check
//
// Returns:
//   - true if the edge exists between the specified nodes
//   - false if the edge doesn't exist or either node doesn't exist
//
// Example:
//
//	ag := New()
//	// ... setup nodes ...
//
//	if !ag.HasEdge(from, to) {
//		ag.AddEdge(from, to)
//	}
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

// IsAcyclic detects whether the graph is acyclic (contains no cycles).
//
// This method implements Kahn's algorithm for topological sorting to detect
// whether the directed graph is acyclic. The algorithm runs asynchronously in a
// separate goroutine and returns the result through a channel.
//
// Algorithm:
// 1. Collect all nodes that participate in edges (have incoming or outgoing edges)
// 2. Calculate in-degree for each node (number of incoming edges)
// 3. Start with nodes that have zero in-degree
// 4. Process nodes by removing them and decrementing in-degrees of neighbors
// 5. If all nodes are processed, the graph is acyclic; otherwise, it has cycles
//
// Returns:
//   - A receive-only channel that will contain the result:
//   - true if the graph is acyclic or empty
//   - false if the graph contains at least one cycle
//
// Concurrency:
// The cycle detection runs in a separate goroutine and the channel is closed
// when the computation completes. This allows for non-blocking cycle detection
// in concurrent applications.
//
// Performance:
// Time complexity: O(V + E) where V is the number of nodes and E is the number of edges
// Space complexity: O(V) for the in-degree map and processing queue
//
// Example:
//
//	ag := New()
//	// ... build graph ...
//
//	isAcyclic := <-ag.IsAcyclic()
//	if isAcyclic {
//		log.Println("Graph is acyclic")
//	} else {
//		log.Println("Graph contains cycles")
//	}
//
//	// Non-blocking usage
//	select {
//	case isAcyclic := <-ag.IsAcyclic():
//		fmt.Printf("Acyclic detection result: %t\n", isAcyclic)
//	case <-time.After(time.Second):
//		fmt.Println("Acyclic detection timed out")
//	}
func (g *Graph) IsAcyclic() <-chan bool {
	ch := make(chan bool)

	go func() {
		defer close(ch)

		q := queue.New()
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

// ForEachNeighbour iterates over all outgoing edges from a node.
//
// This method executes a callback function for each edge originating from
// the specified node. It provides a safe way to traverse a node's neighbors
// with automatic error handling and panic recovery.
//
// Parameters:
//   - node: The GroupNode to iterate neighbors from
//   - fn: The callback function to execute for each outgoing edge
//
// Returns:
//   - nil if the iteration completed successfully
//   - ErrInvalidAdjacency joined with the underlying error if the node doesn't exist
//
// Callback Behavior:
//   - The callback receives each AdjacencyEdge and any error that occurred
//   - Panics in the callback are recovered and reported as errors
//   - If the node has no outgoing edges, no callbacks are executed
//
// Example:
//
//	ag := New()
//	// ... setup graph ...
//
//	node := GroupNode{Id: 1, Group: "users"}
//	err := ag.ForEachNeighbour(node, func(edge AdjacencyEdge, err error) {
//		if err != nil {
//			log.Printf("Error processing edge: %v", err)
//			return
//		}
//		fmt.Printf("Edge from %d to %d (ID: %d)\n", edge.From, edge.To, edge.Edge)
//	})
func (g *Graph) ForEachNeighbour(gn GroupNode, fn OnAdjacencyEdgeFn) error {
	if nodeErr := g.checkNodeExists(gn); nodeErr != nil {
		return errors.Join(ErrInvalidAdjacency, nodeErr)
	}
	g.forEachEdge(gn.ID, fn)
	return nil
}

// GetBackRefsOf returns all nodes that have edges pointing to the specified node.
//
// This method finds all nodes that reference the given node, providing
// efficient reverse traversal capabilities. It's useful for dependency
// analysis and finding all predecessors of a node.
//
// Parameters:
//   - node: The GroupNode to find back-references for
//
// Returns:
//   - A slice of GroupNode instances that have edges pointing to the specified node
//   - ErrInvalidBackRef joined with the underlying error if the node doesn't exist
//   - ErrInvalidBackRef if the node exists but has no back-references
//
// Performance:
// This method has O(G*N) complexity in the worst case, where G is the number
// of groups and N is the average number of nodes per group, as it needs to
// find the group membership for each back-referencing node.
//
// Example:
//
//	ag := New()
//	// ... setup graph with edges pointing to target node ...
//
//	target := GroupNode{Id: 5, Group: "targets"}
//	backRefs, err := ag.GetBackRefsOf(target)
//	if err != nil {
//		log.Printf("Error getting back-references: %v", err)
//		return
//	}
//
//	fmt.Printf("Nodes pointing to %d: %v\n", target.Id, backRefs)
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
//
// This method retrieves all nodes that are members of the given group,
// providing an efficient way to query group membership and iterate
// over nodes within a specific group.
//
// Parameters:
//   - group: The GroupName to retrieve nodes from
//
// Returns:
//   - A slice of GroupNode instances belonging to the specified group
//   - ErrGroupNotFound if the group doesn't exist
//
// Ordering:
// The order of nodes in the returned slice is not guaranteed and may
// vary between calls due to the underlying map iteration.
//
// Example:
//
//	ag := New()
//	ag.AddGroup("users")
//	// ... add nodes to the group ...
//
//	nodes, err := ag.GetNodes("users")
//	if err != nil {
//		log.Printf("Error getting nodes: %v", err)
//		return
//	}
//
//	fmt.Printf("Users group contains %d nodes\n", len(nodes))
//	for _, node := range nodes {
//		fmt.Printf("Node ID: %d\n", node.Id)
//	}
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

// ListGroups returns the names of all groups in the graph.
//
// This method provides a way to enumerate all groups that have been
// created in the graph, regardless of whether they contain nodes.
//
// Returns:
//   - A slice of GroupName instances representing all groups in the graph
//
// Ordering:
// The order of group names in the returned slice is not guaranteed and
// may vary between calls due to the underlying map iteration.
//
// Example:
//
//	ag := New()
//	ag.AddGroup("users")
//	ag.AddGroup("products")
//	ag.AddGroup("orders")
//
//	groups := ag.ListGroups()
//	fmt.Printf("Graph contains %d groups: %v\n", len(groups), groups)
//
//	// Iterate over all groups
//	for _, groupName := range groups {
//		nodes, _ := ag.GetNodes(groupName)
//		fmt.Printf("Group %s has %d nodes\n", groupName, len(nodes))
//	}
func (g *Graph) ListGroups() []GroupName {
	res := make([]GroupName, len(g.groups))
	var i = 0
	for name := range g.groups {
		res[i] = name
		i++
	}
	return res
}
