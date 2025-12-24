package tree

import (
	"errors"

	"github.com/barnowlsnest/go-datalib/pkg/list"
	"github.com/barnowlsnest/go-datalib/pkg/node"
)

type (
	// ChildrenSlice represents a slice of child node values.
	ChildrenSlice = []string

	// HierarchyModel represents a tree structure as a map where keys are node values
	// and values are slices of their children's values.
	// The special RootTag key ("#root") must contain exactly one value representing the root node.
	HierarchyModel = map[string]ChildrenSlice
)

// RootTag is the special key in HierarchyModel that identifies the root node.
const RootTag = "#root"

// Hierarchy builds a tree from a HierarchyModel with cycle detection.
// It constructs a tree structure by traversing the model and creating Node instances.
//
// Parameters:
//   - m: HierarchyModel defining the tree structure
//   - maxBreadth: Maximum number of children per node
//   - nextID: Function to generate unique IDs for each node
//
// Returns an error if:
//   - RootTag is not defined in the model (ErrRootTagNotFound)
//   - Multiple roots are specified (ErrHierarchyModel)
//   - nextID is nil (ErrNil)
//   - maxBreadth < 1 (ErrHierarchyModel)
//   - Root reference doesn't exist in the model (ErrHierarchyModel)
//   - A cycle is detected in the hierarchy (ErrHierarchyModel with "cycle detected" message)
//   - MaxBreadth is exceeded for any node (ErrMaxBreadth)
//
// Cycle Detection:
// The function detects and prevents infinite loops caused by circular references
// (e.g., A→B→A or A→A). When a cycle is detected, it returns an error immediately.
//
// Example:
//
//	model := HierarchyModel{
//	    RootTag: {"Company"},
//	    "Company": {"Engineering", "Sales"},
//	    "Engineering": {"Frontend", "Backend"},
//	}
//	idGen := func() uint64 { return serial.Seq().Next("company") }
//	root, err := Hierarchy(model, 10, idGen)
func Hierarchy(m HierarchyModel, maxBreadth int, nextID func() uint64) (*Node[string], error) {
	rootDef, rootDefined := m[RootTag]
	switch {
	case !rootDefined:
		return nil, ErrRootTagNotFound
	case len(rootDef) != 1:
		return nil, errors.Join(ErrHierarchyModel, errors.New("only 1 root allowed"))
	case nextID == nil:
		return nil, ErrNil
	case maxBreadth < 1:
		return nil, errors.Join(ErrHierarchyModel, errors.New("max breadth should be at least 1"))
	}

	rootNodeVal := rootDef[0]
	rootChildren, rootExists := m[rootNodeVal]
	if !rootExists {
		return nil, errors.Join(ErrHierarchyModel, errors.New("root ref not found"))
	}

	rootID := nextID()
	rootNode, errRoot := NewNode[string](rootID, maxBreadth, ValueOpt[string](rootNodeVal))
	if errRoot != nil {
		return nil, errRoot
	}
	if ok := rootNode.asRoot(); !ok {
		return nil, errors.Join(ErrHierarchyModel, errors.New("unable set root state"))
	}

	stack := list.NewStack()
	lookup := make(map[uint64]*Node[string])
	visited := make(map[string]bool) // Track visited values to detect cycles
	visited[rootNodeVal] = true

	var (
		parent   *Node[string]
		children []string
	)
	parent = rootNode
	children = rootChildren
buildBranch:
	for _, childVal := range children {
		// Check for cycle: if we've seen this value before in our traversal path
		if visited[childVal] {
			return nil, errors.Join(ErrHierarchyModel, errors.New("cycle detected: value \""+childVal+"\" already exists in hierarchy"))
		}

		childID := nextID()
		childNode, errChild := NewNode[string](childID, maxBreadth, ValueOpt[string](childVal))
		if errChild != nil {
			return nil, errChild
		}
		if errAttach := parent.AttachChild(childNode); errAttach != nil {
			return nil, errAttach
		}

		visited[childVal] = true
		lookup[childID] = childNode
		stack.Push(node.ID(childID))
	}

	if stack.IsEmpty() {
		return rootNode, nil
	}

	n := stack.Pop()
	if childNode := lookup[n.ID()]; childNode != nil {
		parent = childNode
		children = m[childNode.Val()]
		goto buildBranch
	} else {
		return nil, ErrNil
	}
}

// ToModel converts a tree (starting from root node) back into a HierarchyModel.
// It performs a breadth-first traversal to build the model representation.
//
// Parameters:
//   - n: The root node of the tree to convert
//
// Returns an error if:
//   - n is nil (ErrNil)
//   - n is not a root node (ErrHierarchyModel)
//   - Internal inconsistencies are detected during traversal (ErrNil, ErrHierarchyModel)
//
// Example:
//
//	root, _ := Hierarchy(model, 10, idGen)
//	modelCopy, err := ToModel(root)  // Reconstruct the model
func ToModel(n *Node[string]) (HierarchyModel, error) {
	if n == nil {
		return nil, ErrNil
	}
	if ok := n.IsRoot(); !ok {
		return nil, errors.Join(ErrHierarchyModel, errors.New("not root"))
	}

	m, lookup := make(HierarchyModel), make(map[uint64]*Node[string])
	m[RootTag] = ChildrenSlice{n.Val()}
	rootID := n.ID()
	lookup[rootID] = n
	queue := list.NewQueue()
	queue.Enqueue(node.ID(rootID))
	for !queue.IsEmpty() {
		next := queue.Dequeue()
		if next == nil {
			return nil, ErrHierarchyModel
		}

		nextNode := lookup[next.ID()]
		if nextNode == nil {
			return nil, ErrNil
		}

		for id, child := range nextNode.ChildrenIter() {
			queue.Enqueue(node.ID(id))
			lookup[id] = child
			if m[nextNode.Val()] == nil {
				m[nextNode.Val()] = make(ChildrenSlice, 0, nextNode.MaxBreadth())
			}
			m[nextNode.Val()] = append(m[nextNode.Val()], child.Val())
		}
	}

	return m, nil
}
