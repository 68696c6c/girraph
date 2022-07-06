package girraph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type CustomTree interface {
	SetName(string) CustomTree
	GetName() string
}

type customTree struct {
	Name string
}

func (t *customTree) SetName(name string) CustomTree {
	t.Name = name
	return t
}

func (t *customTree) GetName() string {
	return t.Name
}

func TestTree(t *testing.T) {
	tree := getTreeFixture()
	assert.Equal(t, "A", tree.GetID())
	assert.Equal(t, "node A", tree.GetMeta().GetName())
	assert.Len(t, tree.GetChildren(), 2)
}

func TestTree_JSON(t *testing.T) {

	// Make a graph.
	tree := getTreeFixture()

	// Convert it to JSON.
	expected, err := tree.JSON()
	println("expected", string(expected))
	require.Nil(t, err)

	// Convert back to a graph.
	treeFromJSON, err := TreeFromJSON[*customTree](expected)
	assert.Equal(t, "node A", treeFromJSON.GetMeta().GetName())
	require.Nil(t, err)

	// Convert back to JSON.
	result, err := treeFromJSON.JSON()
	println("result", string(result))
	require.Nil(t, err)

	// Should match the original JSON.
	assert.Equal(t, expected, result)
}

func getTreeFixture() Tree[CustomTree] {
	nodeA := MakeTree[CustomTree]().SetID("A").SetMeta(&customTree{})
	nodeA.GetMeta().SetName("node A")

	nodeB := MakeTree[CustomTree]().SetID("B").SetMeta(&customTree{})
	nodeB.GetMeta().SetName("node B")

	nodeC := MakeTree[CustomTree]().SetID("C").SetMeta(&customTree{})
	nodeC.GetMeta().SetName("node C")

	nodeD := MakeTree[CustomTree]().SetID("D").SetMeta(&customTree{})
	nodeD.GetMeta().SetName("node D")

	return nodeA.SetChildren([]Tree[CustomTree]{
		nodeB,
		nodeC.SetChildren([]Tree[CustomTree]{
			nodeD,
		}),
	})
}
