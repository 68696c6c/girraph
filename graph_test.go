package girraph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CustomGraph interface {
	SetName(string) CustomGraph
	GetName() string
}

type customGraph struct {
	Name string
}

func (g *customGraph) SetName(name string) CustomGraph {
	g.Name = name
	return g
}

func (g *customGraph) GetName() string {
	return g.Name
}

func TestGraph(t *testing.T) {
	graph := getGraphFixture()
	assert.Equal(t, "A", graph.GetID())
	assert.Equal(t, "node A", graph.GetMeta().GetName())
	assert.Len(t, graph.GetChildren(), 2)
}

func TestGraph_JSON(t *testing.T) {

	// Make a graph.
	graph := getGraphFixture()

	// Convert it to JSON.
	expected, err := graph.JSON()
	require.Nil(t, err)

	// Convert back to a graph.
	graphFromJSON, err := GraphFromJSON[*customGraph](expected)
	assert.Equal(t, "node A", graphFromJSON.GetMeta().GetName())
	require.Nil(t, err)

	// Convert back to JSON.
	result, err := graphFromJSON.JSON()
	require.Nil(t, err)

	// Should match the original JSON.
	assert.Equal(t, expected, result)
}

func TestGraph_SetChildren_AddsParents(t *testing.T) {
	nodeA := MakeGraph[CustomGraph]().SetID("A").SetMeta(&customGraph{})
	nodeA.GetMeta().SetName("node A")

	nodeB := MakeGraph[CustomGraph]().SetID("B").SetMeta(&customGraph{})
	nodeB.GetMeta().SetName("node B")

	nodeC := MakeGraph[CustomGraph]().SetID("C").SetMeta(&customGraph{})
	nodeC.GetMeta().SetName("node C")

	nodeD := MakeGraph[CustomGraph]().SetID("D").SetMeta(&customGraph{})
	nodeD.GetMeta().SetName("node D")

	nodeA.SetChildren([]Graph[CustomGraph]{
		nodeB.SetChildren([]Graph[CustomGraph]{
			nodeD,
		}),
		nodeC.SetChildren([]Graph[CustomGraph]{
			nodeD,
		}),
	})

	assert.Len(t, nodeD.GetParents(), 2)
	assert.Len(t, nodeB.GetParents(), 1)
	assert.Len(t, nodeC.GetParents(), 1)
}

func TestGraph_AddChild_AddsParent(t *testing.T) {
	nodeA := MakeGraph[CustomGraph]().SetID("A").SetMeta(&customGraph{})
	nodeA.GetMeta().SetName("node A")

	nodeB := MakeGraph[CustomGraph]().SetID("B").SetMeta(&customGraph{})
	nodeB.GetMeta().SetName("node B")

	nodeC := MakeGraph[CustomGraph]().SetID("C").SetMeta(&customGraph{})
	nodeC.GetMeta().SetName("node C")

	nodeB.AddChild(nodeC)
	assert.Len(t, nodeC.GetParents(), 1)

	nodeA.AddChild(nodeC)
	assert.Len(t, nodeC.GetParents(), 2)
}

func getGraphFixture() Graph[CustomGraph] {
	nodeA := MakeGraph[CustomGraph]().SetID("A").SetMeta(&customGraph{})
	nodeA.GetMeta().SetName("node A")

	nodeB := MakeGraph[CustomGraph]().SetID("B").SetMeta(&customGraph{})
	nodeB.GetMeta().SetName("node B")

	nodeC := MakeGraph[CustomGraph]().SetID("C").SetMeta(&customGraph{})
	nodeC.GetMeta().SetName("node C")

	nodeD := MakeGraph[CustomGraph]().SetID("D").SetMeta(&customGraph{})
	nodeD.GetMeta().SetName("node D")

	return nodeA.SetChildren([]Graph[CustomGraph]{
		nodeB.SetChildren([]Graph[CustomGraph]{
			nodeD,
		}),
		nodeC.SetChildren([]Graph[CustomGraph]{
			nodeD,
		}),
	})
}
