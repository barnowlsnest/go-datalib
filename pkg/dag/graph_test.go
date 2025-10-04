package dag

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// GraphTestSuite tests Graph wrapper functionality
type GraphTestSuite struct {
	suite.Suite
}

func (s *GraphTestSuite) TestNewGraph() {
	id := uuid.New()
	ag := NewAdjacencyGroups()
	graph := NewGraph(id, "test-graph", ag)

	assert.NotNil(s.T(), graph)
	assert.Equal(s.T(), id, graph.ID())
	assert.Equal(s.T(), "test-graph", graph.Name())
	assert.True(s.T(), graph.IsProvidable())
}

func (s *GraphTestSuite) TestNewEmptyGraph() {
	id := uuid.New()
	graph := NewEmptyGraph(id, "test-graph")

	assert.NotNil(s.T(), graph)
	assert.Equal(s.T(), id, graph.ID())
	assert.Equal(s.T(), "test-graph", graph.Name())
}

func (s *GraphTestSuite) TestNextID() {
	// Use unique graph name to avoid serial ID collision with other tests
	graph := NewEmptyGraph(uuid.New(), uuid.New().String())

	id1 := graph.NextID()
	id2 := graph.NextID()
	id3 := graph.NextID()

	assert.Equal(s.T(), uint64(1), id1)
	assert.Equal(s.T(), uint64(2), id2)
	assert.Equal(s.T(), uint64(3), id3)
}

func (s *GraphTestSuite) TestCurrentID() {
	// Use unique graph name to avoid serial ID collision with other tests
	graph := NewEmptyGraph(uuid.New(), uuid.New().String())

	current := graph.CurrentID()
	assert.Equal(s.T(), uint64(0), current)

	graph.NextID()
	current = graph.CurrentID()
	assert.Equal(s.T(), uint64(1), current)
}

// Test suite runner
func TestGraphTestSuite(t *testing.T) {
	suite.Run(t, new(GraphTestSuite))
}
