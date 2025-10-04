package dag

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/barnowlsnest/go-datalib/pkg/node"
	"github.com/barnowlsnest/go-datalib/pkg/queue"
	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

type (
	// NodeID represents a unique identifier for nodes in the graph.
	// It's an alias for uint64 to provide type safety and clarity.
	NodeID = uint64

	// EdgeID represents a unique identifier for edges in the graph.
	// It's an alias for uint64 to provide type safety and clarity.
	EdgeID = uint64

	// GroupName represents the name of a group within the graph structure.
	// Groups are used to organize nodes into logical collections.
	GroupName = string

	// Name represents a human-readable name for graph entities.
	// It's an alias for string to provide semantic clarity.
	Name = string

	// ID represents a universally unique identifier using UUID.
	// It's an alias for uuid.UUID to provide semantic clarity.
	ID = uuid.UUID

	// GroupNode represents a node that belongs to a specific group.
	//
	// This structure combines a node identifier with its group membership,
	// allowing for efficient organization and querying of nodes within
	// the graph structure.
	GroupNode struct {
		// ID is the unique identifier for this node.
		ID NodeID

		// Group is the name of the group this node belongs to.
		Group GroupName
	}

	// AdjacencyEdge represents a directed edge between two nodes in the graph.
	//
	// This structure contains all the information needed to describe a
	// relationship between nodes, including the source, destination, and
	// a unique edge identifier.
	AdjacencyEdge struct {
		// From is the source node ID of the edge.
		From NodeID

		// To is the destination node ID of the edge.
		To NodeID

		// Edge is the unique identifier for this edge.
		Edge EdgeID
	}

	// BackRefEdge represents a reverse reference edge for efficient traversal.
	//
	// Back-references allow for efficient reverse traversal of the graph
	// by maintaining inverse relationships. This is particularly useful
	// for finding all nodes that point to a specific node.
	BackRefEdge struct {
		// From is the source node ID of the back-reference.
		From NodeID

		// To is the destination node ID of the back-reference.
		To NodeID
	}

	// OnAdjacencyEdgeFn is a callback function type for processing adjacency edges.
	//
	// This function type is used in graph traversal operations where each
	// edge needs to be processed. The function receives an AdjacencyEdge
	// and an error parameter, allowing for both successful edge processing
	// and error handling during traversal.
	//
	// Parameters:
	//   - AdjacencyEdge: The edge being processed
	//   - error: Any error that occurred during edge processing, or nil
	OnAdjacencyEdgeFn func(AdjacencyEdge, error)

	// AdjacencyGroups implements a graph data structure with grouped nodes.
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
	// AdjacencyGroups is not thread-safe. Concurrent access requires external
	// synchronization mechanisms.
	AdjacencyGroups struct {
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

	// Graph represents a complete graph structure with metadata.
	//
	// This structure combines the core graph functionality from AdjacencyGroups
	// with additional metadata such as a unique identifier and human-readable name.
	// It provides a high-level interface for working with graphs in applications.
	//
	// The Graph embeds AdjacencyGroups, inheriting all its methods and functionality
	// while adding identity and naming capabilities.
	Graph struct {
		// AdjacencyGroups provides the core graph functionality.
		// All graph operations are delegated to this embedded structure.
		*AdjacencyGroups

		// id is the unique identifier for this graph instance.
		id ID

		// name is the human-readable name for this graph.
		name Name
	}
)

// NewAdjacencyGroups creates a new empty AdjacencyGroups instance.
//
// This constructor initializes all internal maps required for graph operations.
// The returned instance is ready for use and can immediately accept groups,
// nodes, and edges.
//
// Returns:
//   - A new AdjacencyGroups instance with initialized internal structures
//
// Example:
//
//	ag := NewAdjacencyGroups()
//	ag.AddGroup("users")
//	ag.AddNode(GroupNode{Id: 1, Group: "users"})
func NewAdjacencyGroups() *AdjacencyGroups {
	return &AdjacencyGroups{
		groups:    make(map[GroupName]map[NodeID]struct{}),
		backRefs:  make(map[NodeID]map[NodeID]struct{}),
		adjacency: make(map[NodeID]map[NodeID]EdgeID),
	}
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
func (ag *AdjacencyGroups) checkNodeExists(n GroupNode) error {
	groupNodes, groupExists := ag.groups[n.Group]
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
func (ag *AdjacencyGroups) forEachEdge(from NodeID, fn OnAdjacencyEdgeFn) {
	for to, edge := range ag.adjacency[from] {
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
func (ag *AdjacencyGroups) removeAdjacency(from, to NodeID) {
	delete(ag.adjacency[from], to)
	if len(ag.adjacency[from]) == 0 {
		delete(ag.adjacency, from)
	}
	delete(ag.backRefs[to], from)
	if len(ag.backRefs[to]) == 0 {
		delete(ag.backRefs, to)
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
//	ag := NewAdjacencyGroups()
//	err := ag.AddGroup("users")
//	if err != nil {
//		log.Printf("Failed to add group: %v", err)
//	}
//
//	// Adding the same group again will return an error
//	err = ag.AddGroup("users") // Returns ErrGroupAlreadyExists
func (ag *AdjacencyGroups) AddGroup(name GroupName) error {
	_, groupExists := ag.groups[name]
	if groupExists {
		return errors.Join(ErrGroupAlreadyExists, fmt.Errorf("group [%s]", name))
	}
	ag.groups[name] = make(map[NodeID]struct{})
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
//	ag := NewAdjacencyGroups()
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
func (ag *AdjacencyGroups) AddNode(n GroupNode) error {
	_, groupExists := ag.groups[n.Group]
	if !groupExists {
		return errors.Join(ErrGroupNotFound, fmt.Errorf("group [%s]", n.Group))
	}
	ag.groups[n.Group][n.ID] = struct{}{}
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
//	ag := NewAdjacencyGroups()
//	ag.AddGroup("users")
//	node := GroupNode{Id: 1, Group: "users"}
//	ag.AddNode(node)
//
//	err := ag.RemoveNode(node)
//	if err != nil {
//		log.Printf("Failed to remove node: %v", err)
//	}
func (ag *AdjacencyGroups) RemoveNode(gn GroupNode) error {
	if nodeErr := ag.checkNodeExists(gn); nodeErr != nil {
		return errors.Join(ErrInvalidEdge, nodeErr)
	}
	ag.forEachEdge(gn.ID, func(a AdjacencyEdge, err error) {
		ag.removeAdjacency(a.From, a.To)
	})
	delete(ag.groups[gn.Group], gn.ID)
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
//	ag := NewAdjacencyGroups()
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
func (ag *AdjacencyGroups) AddEdge(from, to GroupNode) error {
	if fromErr := ag.checkNodeExists(from); fromErr != nil {
		return errors.Join(ErrInvalidEdge, fromErr)
	}
	if toErr := ag.checkNodeExists(to); toErr != nil {
		return errors.Join(ErrInvalidEdge, toErr)
	}
	if _, hasNeighbours := ag.adjacency[from.ID]; !hasNeighbours {
		ag.adjacency[from.ID] = make(map[NodeID]EdgeID)
	}
	if _, hasRefs := ag.backRefs[to.ID]; !hasRefs {
		ag.backRefs[to.ID] = make(map[NodeID]struct{})
	}
	ag.adjacency[from.ID][to.ID] = serial.NSum(from.ID, to.ID)
	ag.backRefs[to.ID][from.ID] = struct{}{}
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
//	ag := NewAdjacencyGroups()
//	// ... setup nodes and edge ...
//
//	err := ag.RemoveEdge(from, to)
//	if err != nil {
//		log.Printf("Failed to remove edge: %v", err)
//	}
func (ag *AdjacencyGroups) RemoveEdge(from, to GroupNode) error {
	if fromErr := ag.checkNodeExists(from); fromErr != nil {
		return errors.Join(ErrInvalidEdge, fromErr)
	}
	if toErr := ag.checkNodeExists(to); toErr != nil {
		return errors.Join(ErrInvalidEdge, toErr)
	}
	ag.removeAdjacency(from.ID, to.ID)
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
//	ag := NewAdjacencyGroups()
//	ag.AddGroup("users")
//	node := GroupNode{Id: 1, Group: "users"}
//
//	if !ag.HasNode(node) {
//		ag.AddNode(node)
//	}
func (ag *AdjacencyGroups) HasNode(gn GroupNode) bool {
	if err := ag.checkNodeExists(gn); err != nil {
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
//	ag := NewAdjacencyGroups()
//	// ... setup nodes ...
//
//	if !ag.HasEdge(from, to) {
//		ag.AddEdge(from, to)
//	}
func (ag *AdjacencyGroups) HasEdge(from, to GroupNode) bool {
	if fromErr := ag.checkNodeExists(from); fromErr != nil {
		return false
	}
	if toErr := ag.checkNodeExists(to); toErr != nil {
		return false
	}
	toNodes, toNodeExists := ag.adjacency[from.ID]
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
//	ag := NewAdjacencyGroups()
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
func (ag *AdjacencyGroups) IsAcyclic() <-chan bool {
	ch := make(chan bool)

	go func() {
		defer close(ch)

		q := queue.New()
		in := make(map[NodeID]int)

		// Collect all nodes from the graph (both with outgoing and incoming edges)
		allNodes := make(map[NodeID]struct{})

		// Add nodes with outgoing edges
		for nodeID := range ag.adjacency {
			allNodes[nodeID] = struct{}{}
		}

		// Add nodes with incoming edges
		for nodeID := range ag.backRefs {
			allNodes[nodeID] = struct{}{}
		}

		// If there are no nodes, the graph is empty and is acyclic
		if len(allNodes) == 0 {
			ch <- true
			return
		}

		// Initialize in-degree for all nodes
		for nodeID := range allNodes {
			refs, exists := ag.backRefs[nodeID]
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
			neighbors, hasNeighbors := ag.adjacency[nodeID]
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
//	ag := NewAdjacencyGroups()
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
func (ag *AdjacencyGroups) ForEachNeighbour(gn GroupNode, fn OnAdjacencyEdgeFn) error {
	if nodeErr := ag.checkNodeExists(gn); nodeErr != nil {
		return errors.Join(ErrInvalidAdjacency, nodeErr)
	}
	ag.forEachEdge(gn.ID, fn)
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
//	ag := NewAdjacencyGroups()
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
func (ag *AdjacencyGroups) GetBackRefsOf(gn GroupNode) ([]GroupNode, error) {
	if nodeErr := ag.checkNodeExists(gn); nodeErr != nil {
		return nil, errors.Join(ErrInvalidBackRef, nodeErr)
	}
	backRefs, hasBackRefs := ag.backRefs[gn.ID]
	if !hasBackRefs {
		return nil, ErrInvalidBackRef
	}
	res := make([]GroupNode, len(backRefs))
	var i int
	for ref := range backRefs {
		for group, nodes := range ag.groups {
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
//	ag := NewAdjacencyGroups()
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
func (ag *AdjacencyGroups) GetNodes(group GroupName) ([]GroupNode, error) {
	groupNodes, groupExists := ag.groups[group]
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
//	ag := NewAdjacencyGroups()
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
func (ag *AdjacencyGroups) ListGroups() []GroupName {
	res := make([]GroupName, len(ag.groups))
	var i = 0
	for name := range ag.groups {
		res[i] = name
		i++
	}
	return res
}

// NewGraph creates a new Graph with the specified ID, name, and adjacency groups.
//
// This constructor creates a Graph instance that wraps an existing
// AdjacencyGroups structure with additional metadata. The graph inherits
// all functionality from the provided AdjacencyGroups.
//
// Parameters:
//   - id: The unique identifier for this graph instance
//   - name: The human-readable name for this graph
//   - groups: The AdjacencyGroups instance to wrap
//
// Returns:
//   - A new Graph instance with the specified configuration
//
// Example:
//
//	id := uuid.New()
//	ag := NewAdjacencyGroups()
//	ag.AddGroup("users")
//
//	graph := NewGraph(id, "user-relationships", ag)
//	fmt.Printf("Created graph: %s\n", graph.Name())
func NewGraph(id ID, name string, groups *AdjacencyGroups) *Graph {
	return &Graph{
		id:              id,
		name:            name,
		AdjacencyGroups: groups,
	}
}

// NewEmptyGraph creates a new Graph with empty adjacency groups.
//
// This convenience constructor creates a Graph with a new, empty
// AdjacencyGroups instance. It's equivalent to calling NewGraph
// with a newly created AdjacencyGroups.
//
// Parameters:
//   - id: The unique identifier for this graph instance
//   - name: The human-readable name for this graph
//
// Returns:
//   - A new Graph instance with empty adjacency groups
//
// Example:
//
//	id := uuid.New()
//	graph := NewEmptyGraph(id, "my-graph")
//
//	// Graph is ready for use
//	graph.AddGroup("nodes")
//	graph.AddNode(GroupNode{Id: 1, Group: "nodes"})
func NewEmptyGraph(id ID, name string) *Graph {
	return NewGraph(id, name, NewAdjacencyGroups())
}

// ID returns the unique identifier of this graph.
//
// The ID is immutable after graph creation and serves as the primary
// way to identify and reference this graph instance.
//
// Returns:
//   - The UUID assigned to this graph during creation
//
// Example:
//
//	graph := NewEmptyGraph(uuid.New(), "my-graph")
//	fmt.Printf("Graph ID: %s\n", graph.ID())
func (g *Graph) ID() ID {
	return g.id
}

// Name returns the human-readable name of this graph.
//
// The name is immutable after graph creation and provides a
// descriptive identifier for this graph instance.
//
// Returns:
//   - The name assigned to this graph during creation
//
// Example:
//
//	graph := NewEmptyGraph(uuid.New(), "user-relationships")
//	fmt.Printf("Graph name: %s\n", graph.Name())
func (g *Graph) Name() Name {
	return g.name
}

// NextID generates and returns the next sequential node ID for this graph.
//
// This method uses the global serial ID generator with the graph's name
// as the key, ensuring that each graph maintains its own independent
// sequence of node IDs. The generated IDs are unique within the context
// of this graph's name.
//
// Returns:
//   - The next sequential NodeID for this graph (starting from 1)
//
// Thread Safety:
// This method is thread-safe and can be called concurrently from
// multiple goroutines.
//
// Example:
//
//	graph := NewEmptyGraph(uuid.New(), "my-graph")
//
//	id1 := graph.NextID() // Returns 1
//	id2 := graph.NextID() // Returns 2
//	id3 := graph.NextID() // Returns 3
//
//	// Use the generated ID to create nodes
//	node := GroupNode{ID: id1, Group: "users"}
func (g *Graph) NextID() NodeID {
	return serial.Seq().Next(g.name)
}

// CurrentID returns the current node ID value for this graph without incrementing.
//
// This method provides read-only access to the current ID counter value
// for this graph. It's useful for checking the current state without
// generating a new ID.
//
// Returns:
//   - The current NodeID value for this graph (0 if never incremented)
//
// Thread Safety:
// This method is thread-safe and can be called concurrently from
// multiple goroutines.
//
// Example:
//
//	graph := NewEmptyGraph(uuid.New(), "my-graph")
//
//	current := graph.CurrentID() // Returns 0 (not yet incremented)
//	graph.NextID()               // Returns 1
//	current = graph.CurrentID()  // Returns 1
func (g *Graph) CurrentID() NodeID {
	return serial.Seq().Current(g.name)
}

// IsProvidable indicates whether this graph can be provided/used.
//
// This method always returns true, indicating that Graph instances
// are always ready for use. This method may be part of an interface
// contract for providable entities in the system.
//
// Returns:
//   - Always returns true
//
// Example:
//
//	graph := NewEmptyGraph(uuid.New(), "my-graph")
//	if graph.IsProvidable() {
//		// Graph is ready for use
//		graph.AddGroup("nodes")
//	}
func (g *Graph) IsProvidable() bool {
	return g != nil
}
