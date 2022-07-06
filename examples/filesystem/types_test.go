package filesystem

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/68696c6c/girraph"
)

func TestDirectory(t *testing.T) {
	tree := getDirectoryFixture()
	assert.NotNil(t, tree.GetID())
	require.Equal(t, "A", tree.GetMeta().GetName())
}

func TestDirectory_JSON(t *testing.T) {

	// Make a graph.
	tree := getDirectoryFixture()

	// Convert it to JSON.
	expected, err := tree.JSON()
	println("expected", string(expected))
	require.Nil(t, err)

	// Convert back to a graph.
	treeFromJSON, err := girraph.TreeFromJSON[*directory](expected)
	assert.Equal(t, "A", treeFromJSON.GetMeta().GetName())
	require.Nil(t, err)

	// Convert back to JSON.
	result, err := treeFromJSON.JSON()
	println("result", string(result))
	require.Nil(t, err)

	// Should match the original JSON.
	assert.Equal(t, expected, result)
	// require.False(t, true)
}

func getDirectoryFixture() girraph.Tree[Directory] {
	dirA := MakeDirectory("A")
	dirA.GetMeta().SetFiles([]*file{
		MakeFile("one", "go"),
		MakeFile("two", "go"),
		MakeFile("three", "go"),
	})

	dirB := MakeDirectory("B")
	dirB.GetMeta().SetFiles([]*file{
		MakeFile("one", "js"),
		MakeFile("two", "js"),
		MakeFile("three", "js"),
	})

	dirC := MakeDirectory("C")

	dirD := MakeDirectory("D")
	dirD.GetMeta().SetFiles([]*file{
		MakeFile("one", "yml"),
		MakeFile("two", "yml"),
		MakeFile("three", "yml"),
	})

	return dirA.SetChildren([]girraph.Tree[Directory]{
		dirB,
		dirC.SetChildren([]girraph.Tree[Directory]{
			dirD,
		}),
	})
}
