package dag

import (
	"github.com/barnowlsnest/go-datalib/pkg/serial"
)

// Graph represents a complete graph structure with metadata.
//
// This structure combines the core graph functionality from AdjacencyGroups
// with additional metadata such as a unique identifier and human-readable name.
// It provides a high-level interface for working with graphs in applications.
//
// The Graph embeds AdjacencyGroups, inheriting all its methods and functionality
// while adding identity and naming capabilities.
type Graph struct {
	// AdjacencyGroups provides the core graph functionality.
	// All graph operations are delegated to this embedded structure.
	*AdjacencyGroups

	// id is the unique identifier for this graph instance.
	id ID

	// name is the human-readable name for this graph.
	name Name
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
	return g != nil && g.AdjacencyGroups != nil
}
