package girraph

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Tree[T any] interface {
	Node[Tree[T]]
	SetParent(Tree[T]) Tree[T]
	GetParent() Tree[T]
	SetMeta(T) Tree[T]
	GetMeta() T
}

type TreeNode[T any] struct {
	ID       string
	Meta     T
	Children []Tree[T]
	parent   Tree[T]
}

func MakeTree[T any]() Tree[T] {
	return &TreeNode[T]{
		ID:       uuid.New().String(),
		Children: []Tree[T]{},
		parent:   nil,
	}
}

func (t *TreeNode[T]) SetID(id string) Tree[T] {
	t.ID = id
	return t
}

func (t *TreeNode[_]) GetID() string {
	return t.ID
}

func (t *TreeNode[T]) SetChildren(children []Tree[T]) Tree[T] {
	t.Children = children
	return t
}

func (t *TreeNode[T]) GetChildren() []Tree[T] {
	return t.Children
}

func (t *TreeNode[T]) AddChild(child Tree[T]) Tree[T] {
	t.Children = append(t.Children, child)
	return t
}

func (t *TreeNode[T]) AddParent(parent Tree[T]) Tree[T] {
	t.parent = parent
	return t
}

func (t *TreeNode[T]) SetParent(parent Tree[T]) Tree[T] {
	t.parent = parent
	return t
}

func (t *TreeNode[T]) GetParents() []Tree[T] {
	return []Tree[T]{t.parent}
}

func (t *TreeNode[T]) GetParent() Tree[T] {
	return t.parent
}

func (t *TreeNode[_]) JSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TreeNode[T]) SetMeta(meta T) Tree[T] {
	t.Meta = meta
	return t
}

func (t *TreeNode[T]) GetMeta() T {
	return t.Meta
}

func TreeFromJSON[T any](input []byte) (Tree[T], error) {
	temp := &NodeJSON[T]{}
	err := json.Unmarshal(input, temp)
	if err != nil {
		return nil, err
	}
	result := TreeFromNode[T](temp)
	return result, nil
}

func TreeFromNode[T any](n *NodeJSON[T]) Tree[T] {
	result := &TreeNode[T]{
		ID:       n.ID,
		Meta:     n.Meta,
		Children: []Tree[T]{},
		parent:   nil,
	}
	for _, child := range n.Children {
		result.AddChild(TreeFromNode[T](child))
	}
	return result
}
