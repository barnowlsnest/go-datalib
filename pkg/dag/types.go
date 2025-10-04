package dag

import (
	"github.com/google/uuid"
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
)
