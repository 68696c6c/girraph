package girraph

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Graph[T any] interface {
	Node[Graph[T]]
	AddParent(Graph[T]) Graph[T]
	SetParents([]Graph[T]) Graph[T]
	SetMeta(T) Graph[T]
	GetMeta() T
}

type graph[T any] struct {
	ID       string
	Meta     T
	Children []Graph[T]
	parents  []Graph[T]
}

func MakeGraph[T any]() Graph[T] {
	return &graph[T]{
		ID:       uuid.New().String(),
		Children: []Graph[T]{},
		parents:  []Graph[T]{},
	}
}

func (g *graph[T]) SetID(id string) Graph[T] {
	g.ID = id
	return g
}

func (g *graph[T]) GetID() string {
	return g.ID
}

func (g *graph[T]) SetChildren(children []Graph[T]) Graph[T] {
	g.Children = children
	return g
}

func (g *graph[T]) GetChildren() []Graph[T] {
	return g.Children
}

func (g *graph[T]) AddChild(child Graph[T]) Graph[T] {
	g.Children = append(g.Children, child)
	return g
}

func (g *graph[T]) AddParent(parent Graph[T]) Graph[T] {
	g.parents = append(g.parents, parent)
	return g
}

func (g *graph[T]) SetParents(parents []Graph[T]) Graph[T] {
	g.parents = parents
	return g
}

func (g *graph[T]) GetParents() []Graph[T] {
	return g.parents
}

func (g *graph[T]) JSON() ([]byte, error) {
	return json.Marshal(g)
}

func (g *graph[T]) SetMeta(meta T) Graph[T] {
	g.Meta = meta
	return g
}

func (g *graph[T]) GetMeta() T {
	return g.Meta
}

func SetParents[T any](input Graph[T]) {
	for _, child := range input.GetChildren() {
		child.AddParent(input)
		SetParents(child)
	}
}

type nodeCount[T any] struct {
	count int
	node  Graph[T]
}

func GraphFromJSON[T any](input []byte) (Graph[T], error) {
	temp := &NodeJSON[T]{}
	err := json.Unmarshal(input, temp)
	if err != nil {
		return nil, err
	}
	result := GraphFromNode[T](temp)

	// Set the parents of all atoms.
	SetParents(result)

	// Count the number of times each atom id appears in the
	nodeCounter := make(map[string]nodeCount[T])
	Traverse[Graph[T]](result, func(node Graph[T]) {
		id := node.GetID()
		if tmp, exists := nodeCounter[id]; !exists {
			nodeCounter[id] = nodeCount[T]{
				count: 1,
				node:  node,
			}
		} else {
			tmp.count += 1
		}
	})

	// Replace each duplicate atom with a single atom that has all the parents of the duplicates.
	for nodeID, data := range nodeCounter {
		if data.count < 2 {
			continue
		}

		// We will replace all duplicate instances with a single instance.
		replacementChild := data.node

		// Get all the parents of each instance of this atom.
		parents := FindAllParentsByID[Graph[T]](result, nodeID)

		// Set the complete list of parents on our new only-child atom.
		replacementChild.SetParents(parents)

		// Replace each instance with our replacement.
		for _, parent := range parents {
			parentChildren := parent.GetChildren()
			for childIndex, parentChild := range parentChildren {
				if parentChild.GetID() == nodeID {
					parentChildren[childIndex] = replacementChild
				}
			}
		}
	}

	return result, nil
}

func GraphFromNode[T any](n *NodeJSON[T]) Graph[T] {
	result := &graph[T]{
		ID:       n.ID,
		Meta:     n.Meta,
		Children: []Graph[T]{},
		parents:  []Graph[T]{},
	}
	for _, child := range n.Children {
		result.AddChild(GraphFromNode[T](child))
	}
	return result
}
