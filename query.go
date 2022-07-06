package girraph

// Get all parent nodes for all nodes with the specified id.
func FindAllParentsByID[T Node[T]](root T, id string) []T {
	var parents []T
	parentMap := make(map[string]bool)
	for _, node := range FindNodesByID[T](root, id) {
		for _, parent := range node.GetParents() {
			parentID := parent.GetID()
			_, exists := parentMap[parentID]
			if !exists {
				parents = append(parents, parent)
				parentMap[parentID] = true
			}
		}
	}
	return parents
}

func FindNodesByID[T Node[T]](root T, id string) []T {
	var found QueryResult[T]
	queryNodesByID[T](&found, root, id)
	return found.nodes
}

func Traverse[T Node[T]](input T, callback func(T)) {
	callback(input)
	for _, child := range input.GetChildren() {
		Traverse[T](child, callback)
	}
}

type QueryResult[T Node[T]] struct {
	nodes []T
}

func (q *QueryResult[T]) AddNode(atom T) {
	q.nodes = append(q.nodes, atom)
}

func (q *QueryResult[T]) GetNodes() []T {
	return q.nodes
}

func queryNodesByID[T Node[T]](query *QueryResult[T], node T, id string) {
	if node.GetID() == id {
		query.AddNode(node)
	}
	for _, child := range node.GetChildren() {
		queryNodesByID[T](query, child, id)
	}
}
