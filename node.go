package girraph

type Node[T any] interface {
	SetID(string) T
	GetID() string
	SetChildren([]T) T
	GetChildren() []T
	AddChild(T) T
	GetParents() []T
	JSON() ([]byte, error)
}

type NodeJSON[T any] struct {
	ID       string
	Meta     T
	Children []*NodeJSON[T]
	parents  []*NodeJSON[T]
}
